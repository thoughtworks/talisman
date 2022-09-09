package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime/pprof"
	"strings"
	"talisman/utility"
	"time"

	"talisman/prompt"
	"talisman/talismanrc"

	"runtime"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	flag "github.com/spf13/pflag"
)

var (
	showVersion bool
	//Version : Version of talisman
	Version       = "Development Build"
	interactive   bool
	talismanInput io.Reader
)

const (
	//PrePush : Const for name of pre-push hook
	PrePush = "pre-push"
	//PreCommit : Const for name of of pre-commit hook
	PreCommit = "pre-commit"
	//EXIT_SUCCESS : Const to indicate successful talisman invocation
	EXIT_SUCCESS = 0
	//EXIT_FAILURE : Const to indicate failed successful invocation
	EXIT_FAILURE = 1
)

var options struct {
	Debug           bool
	LogLevel        string
	GitHook         string
	Pattern         string
	Scan            bool
	IgnoreHistory   bool
	Checksum        string
	ReportDirectory string
	ScanWithHtml    bool
	ShouldProfile   bool
}

//var options Options

func init() {
	log.SetOutput(os.Stderr)
	talismanInput = os.Stdin
	flag.BoolVarP(&options.Debug,
		"debug", "d", false,
		"enable debug mode (warning: very verbose)")
	flag.StringVarP(&options.LogLevel,
		"loglevel", "l", "error",
		"set log level for talisman (allowed values: error|info|warn|debug, default: error)")
	flag.BoolVarP(&showVersion,
		"version", "v", false,
		"show current version of talisman")
	flag.StringVarP(&options.Pattern,
		"pattern", "p", "",
		"pattern (glob-like) of files to scan (ignores githooks)")
	flag.StringVarP(&options.GitHook,
		"githook", "g", PrePush,
		"either pre-push or pre-commit")
	flag.BoolVarP(&options.Scan,
		"scan", "s", false,
		"scanner scans the git commit history for potential secrets")
	flag.BoolVarP(&options.IgnoreHistory,
		"ignoreHistory", "^", false,
		"scanner scans all files on current head, will not scan through git commit history")
	flag.StringVarP(&options.Checksum,
		"checksum", "c", "",
		"checksum calculator calculates checksum and suggests .talismanrc entry")
	flag.StringVarP(&options.ReportDirectory,
		"reportDirectory", "r", "talisman_report",
		"directory where the scan report will be stored")
	flag.BoolVarP(&options.ScanWithHtml,
		"scanWithHtml", "w", false,
		"generate html report (**Make sure you have installed talisman_html_report to use this, as mentioned in talisman Readme**)")
	flag.BoolVarP(&interactive,
		"interactive", "i", false,
		"interactively update talismanrc (only makes sense with -g/--githook)")
	flag.BoolVarP(&options.ShouldProfile,
		"profile", "f", false,
		"profile cpu and memory usage of talisman")
}

func main() {
	flag.Parse()

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		os.Exit(EXIT_SUCCESS)
	}

	if showVersion {
		fmt.Printf("talisman %s\n", Version)
		os.Exit(EXIT_SUCCESS)
	}

	if options.GitHook != "" {
		if !(options.GitHook == PreCommit || options.GitHook == PrePush) {
			fmt.Println(fmt.Errorf("githook should be %s or %s, but got %s", PreCommit, PrePush, options.GitHook))
			os.Exit(EXIT_FAILURE)
		}
	}

	if options.ShouldProfile {
		stopProfFunc := setupProfiling()
		defer stopProfFunc()
	}

	promptContext := prompt.NewPromptContext(interactive, prompt.NewPrompt())
	os.Exit(run(promptContext))
}

func run(promptContext prompt.PromptContext) (returnCode int) {
	start := time.Now()
	defer func() { fmt.Printf("Talisman done in %v\n", time.Since(start)) }()

	if err := validateGitExecutable(afero.NewOsFs(), runtime.GOOS); err != nil {
		log.Errorf("error validating git executable: %v", err)
		return 1
	}

	setLogLevel()

	if options.GitHook == "" {
		options.GitHook = PrePush
	}

	optionsBytes, _ := json.Marshal(options)
	fields := make(map[string]interface{})
	_ = json.Unmarshal(optionsBytes, &fields)
	log.WithFields(fields).Debug("Talisman execution environment")
	defer  utility.DestroyHashers()
	if options.Checksum != "" {
		log.Infof("Running %s patterns against checksum calculator", options.Checksum)
		return NewChecksumCmd(strings.Fields(options.Checksum)).Run()
	} else if options.Scan {
		log.Infof("Running scanner")
		return NewScannerCmd(options.IgnoreHistory, options.ReportDirectory).Run(talismanrc.ForScan(options.IgnoreHistory))
	} else if options.ScanWithHtml {
		log.Infof("Running scanner with html report")
		return NewScannerCmd(options.IgnoreHistory, "talisman_html_report").Run(talismanrc.ForScan(options.IgnoreHistory))
	} else if options.Pattern != "" {
		log.Infof("Running scan for %s", options.Pattern)
		return NewPatternCmd(options.Pattern).Run(talismanrc.For(talismanrc.HookMode), promptContext)
	} else if options.GitHook == PreCommit {
		log.Infof("Running %s hook", options.GitHook)
		return NewPreCommitHook().Run(talismanrc.For(talismanrc.HookMode), promptContext)
	} else {
		log.Infof("Running %s hook", options.GitHook)
		return NewPrePushHook(talismanInput).Run(talismanrc.For(talismanrc.HookMode), promptContext)
	}
}

func validateGitExecutable(fs afero.Fs, operatingSystem string) error {
	if operatingSystem == "windows" {
		extensions := strings.ToLower(os.Getenv("PATHEXT"))
		windowsExecutables := strings.Split(extensions, ";")
		for _, executable := range windowsExecutables {
			gitExecutable := fmt.Sprintf("git%s", executable)
			exists, err := afero.Exists(fs, gitExecutable)
			if err != nil {
				return err
			}
			if exists {
				return fmt.Errorf("not allowed to have git executable located in repository: %s", gitExecutable)
			}
		}
	}
	return nil
}

func setLogLevel() {
	switch options.LogLevel {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	default:
		log.SetLevel(log.ErrorLevel)
	}
	if options.Debug {
		log.SetLevel(log.DebugLevel)
	}
}

func setupProfiling() func() {
	log.Info("Profiling initiated")

	cpuProfFile, err := os.Create("talisman.cpuprof")
	if err != nil {
		log.Fatalf("Unable to create cpu profiling output file talisman.cpuprof: %v", err)
	}

	memProfFile, err := os.Create("talisman.memprof")
	if err != nil {
		log.Fatalf("Unable to create memory profiling output file talisman.memprof: %v", err)
	}

	_ = pprof.StartCPUProfile(cpuProfFile)
	progEnded := false

	go func() {
		memProfTimer := time.NewTimer(500 * time.Millisecond)

		for !progEnded {
			<-memProfTimer.C
			_ = pprof.WriteHeapProfile(memProfFile)
		}
		_ = memProfFile.Close()
	}()

	return func() {
		progEnded = true
		pprof.StopCPUProfile()
		log.Info("Profiling completed")
		_ = cpuProfFile.Close()
	}
}
