package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime/pprof"
	"strings"
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
	Version     = "Development Build"
	interactive bool
)

const (
	//PrePush : Const for name of pre-push hook
	PrePush = "pre-push"
	//PreCommit : Const for name of of pre-commit hook
	PreCommit = "pre-commit"
)

var options struct {
	debug           bool
	loglevel        string
	githook         string
	pattern         string
	scan            bool
	ignoreHistory   bool
	checksum        string
	reportDirectory string
	scanWithHtml    bool
	input           io.Reader
	shouldProfile   bool
}

//var options Options

func init() {
	log.SetOutput(os.Stderr)
	flag.BoolVarP(&options.debug, "debug", "d", false, "enable debug mode (warning: very verbose)")
	flag.StringVarP(&options.loglevel, "loglevel", "l", "error", "enable debug mode (warning: very verbose)")
	flag.BoolVarP(&showVersion, "version", "v", false, "show current version of talisman")
	flag.StringVarP(&options.pattern, "pattern", "p", "", "pattern (glob-like) of files to scan (ignores githooks)")
	flag.StringVarP(&options.githook, "githook", "g", PrePush, "either pre-push or pre-commit")
	flag.BoolVarP(&options.scan, "scan", "s", false, "scanner scans the git commit history for potential secrets")
	flag.BoolVar(&options.ignoreHistory, "ignoreHistory", false, "scanner scans all files on current head, will not scan through git commit history")
	flag.StringVarP(&options.checksum, "checksum", "c", "", "checksum calculator calculates checksum and suggests .talismanrc format")
	flag.StringVarP(&options.reportDirectory, "reportDirectory", "r", "talisman_report", "directory where the scan reports will be stored")
	flag.BoolVarP(&options.scanWithHtml, "scanWithHtml", "w", false, "generate html report (**Make sure you have installed talisman_html_report to use this, as mentioned in Readme**)")
	flag.BoolVarP(&interactive, "interactive", "i", false, "interactively update talismanrc (only makes sense with -g/--githook)")
	flag.BoolVarP(&options.shouldProfile, "profile", "", false, "interactively update talismanrc (only makes sense with -g/--githook)")

}

func main() {
	flag.Parse()

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if showVersion {
		fmt.Printf("talisman %s\n", Version)
		os.Exit(0)
	}

	if options.githook != "" {
		if !(options.githook == PreCommit || options.githook == PrePush) {
			fmt.Println(fmt.Errorf("githook should be %s or %s, but got %s", PreCommit, PrePush, options.githook))
			os.Exit(1)
		}
	}

	prompter := prompt.NewPrompt()
	promptContext := prompt.NewPromptContext(interactive, prompter)
	options.input = os.Stdin
	os.Exit(run(promptContext))
}

func run(promptContext prompt.PromptContext) (returnCode int) {
	if options.shouldProfile {
		log.Info("Profiling initiated")
		defer func() { log.Info("Profiling completed") }()
		f, err := os.Create("talisman.pprof")
		if err != nil {
			log.Fatal(err)
		}
		_ = pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
		progEnded := false
		defer func() { progEnded = true }()
		go func() {
			t := time.NewTimer(500 * time.Millisecond)
			f1, _ := os.Create("talisman.mprof")
			for !progEnded {
				select {
				case <-t.C:
					if f1 != nil {
						_ = pprof.WriteHeapProfile(f1)
					} else {
						log.Error("Could not write memory profiling info")
					}

				}
			}
			if f1 != nil {
				_ = f1.Close()
			}
		}()
	}
	start := time.Now()
	defer func() { fmt.Printf("Talisman done in %v", time.Now().Sub(start)) }()
	if err := validateGitExecutable(afero.NewOsFs(), runtime.GOOS); err != nil {
		log.Errorf("error validating git executable:"+" %v", err)
		return 1
	}

	switch options.loglevel {
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
	if options.debug {
		log.SetLevel(log.DebugLevel)
	}

	if options.githook == "" {
		options.githook = PrePush
	}
	bytes, _ := json.Marshal(options)
	fields := make(map[string]interface{})
	_ = json.Unmarshal(bytes, fields)
	log.WithFields(fields).Debug("Execution environment")

	if options.checksum != "" {
		log.Infof("Running %s patterns against checksum calculator", options.checksum)
		return NewChecksumCmd(strings.Fields(options.checksum)).Run()
	} else if options.scan {
		log.Infof("Running scanner")
		return NewScannerCmd(options.ignoreHistory, options.reportDirectory).Run(talismanrc.For(talismanrc.ScanMode))
	} else if options.scanWithHtml {
		log.Infof("Running scanner with html report")
		return NewScannerCmd(options.ignoreHistory, "talisman_html_report").Run(talismanrc.For(talismanrc.ScanMode))
	} else if options.pattern != "" {
		log.Infof("Running scan for %s", options.pattern)
		return NewPatternCmd(options.pattern).Run(talismanrc.For(talismanrc.HookMode), promptContext)
	} else if options.githook == PreCommit {
		log.Infof("Running %s hook", options.githook)
		return NewPreCommitHook().Run(talismanrc.For(talismanrc.HookMode), promptContext)
	} else {
		log.Infof("Running %s hook", options.githook)
		return NewPrePushHook(options.input).Run(talismanrc.For(talismanrc.HookMode), promptContext)
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
