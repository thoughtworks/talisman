package main

import (
	"github.com/spf13/afero"
	flag "github.com/spf13/pflag"
	"runtime"
	"talisman/prompt"
	"talisman/talismanrc"
)

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"talisman/gitrepo"

	log "github.com/Sirupsen/logrus"
)

var (
	showVersion bool
	//Version : Version of talisman
	Version         = "Development Build"
	interactive     bool
)

const (
	//PrePush : Const for name of pre-push hook
	PrePush = "pre-push"
	//PreCommit : Const for name of of pre-commit hook
	PreCommit = "pre-commit"
)

func init() {
	log.SetOutput(os.Stderr)
}

type options struct {
	debug           bool
	githook         string
	pattern         string
	scan            bool
	ignoreHistory   bool
	checksum        string
	reportdirectory string
	scanWithHtml    bool
}

//Logger is the default log device, set to emit at the Error level by default
func main() {
	_options := options{}
	flag.BoolVarP(&_options.debug, "debug", "d", false, "enable debug mode (warning: very verbose)")
	flag.BoolVarP(&showVersion, "version", "v", false, "show current version of talisman")
	flag.StringVarP(&_options.pattern, "pattern", "p", "", "pattern (glob-like) of files to scan (ignores githooks)")
	flag.StringVarP(&_options.githook, "githook", "g", PrePush, "either pre-push or pre-commit")
	flag.BoolVarP(&_options.scan, "scan", "s", false, "scanner scans the git commit history for potential secrets")
	flag.BoolVar(&_options.ignoreHistory, "ignoreHistory", false, "scanner scans all files on current head, will not scan through git commit history")
	flag.StringVarP(&_options.checksum, "checksum", "c", "", "checksum calculator calculates checksum and suggests .talismanrc format")
	flag.StringVarP(&_options.reportdirectory, "reportdirectory", "r", "", "directory where the scan reports will be stored")
	flag.BoolVarP(&_options.scanWithHtml, "scanWithHtml", "w", false, "generate html report (**Make sure you have installed talisman_html_report to use this, as mentioned in Readme**)")
	flag.BoolVarP(&interactive, "interactive", "i", false, "interactively update talismanrc (only makes sense with -g/--githook)")

	flag.Parse()

	if showVersion {
		fmt.Printf("talisman %s\n", Version)
		os.Exit(0)
	}

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if _options.githook != "" {
		if !(_options.githook == PreCommit || _options.githook == PrePush) {
			fmt.Println(fmt.Errorf("githook should be %s or %s, but got %s", PreCommit, PrePush, _options.githook))
			os.Exit(1)
		}
	}

	prompter := prompt.NewPrompt()
	promptContext := prompt.NewPromptContext(interactive, prompter)

	os.Exit(run(os.Stdin, _options, promptContext))
}

func run(stdin io.Reader, _options options, promptContext prompt.PromptContext) (returnCode int) {
	if err := validateGitExecutable(afero.NewOsFs(), runtime.GOOS); err != nil {
		log.Printf("error validating git executable:" +
			" %v", err)
		return 1
	}

	if _options.debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.ErrorLevel)
	}

	if _options.githook == "" {
		_options.githook = PrePush
	}

	var additions []gitrepo.Addition
	if _options.checksum != "" {
		log.Infof("Running %s patterns against checksum calculator", _options.checksum)
		return NewRunner(make([]gitrepo.Addition, 0)).RunChecksumCalculator(strings.Fields(_options.checksum))
	} else if _options.scan {
		log.Infof("Running scanner")
		return NewRunner(make([]gitrepo.Addition, 0)).Scan(_options.reportdirectory, talismanrc.Get() , _options.ignoreHistory)
	} else if _options.scanWithHtml {
		log.Infof("Running scanner with html report")
		return NewRunner(make([]gitrepo.Addition, 0)).Scan("talisman_html_report", talismanrc.Get(), _options.ignoreHistory)
	} else if _options.pattern != "" {
		log.Infof("Running %s pattern", _options.pattern)
		directoryHook := NewDirectoryHook()
		additions = directoryHook.GetFilesFromDirectory(_options.pattern)
	} else if _options.githook == PreCommit {
		log.Infof("Running %s hook", _options.githook)
		preCommitHook := NewPreCommitHook()
		additions = preCommitHook.GetRepoAdditions()
	} else {
		log.Infof("Running %s hook", _options.githook)
		prePushHook := NewPrePushHook(readRefAndSha(stdin))
		additions = prePushHook.GetRepoAdditions()
	}

	return NewRunner(additions).RunWithoutErrors(promptContext)
}

func readRefAndSha(file io.Reader) (string, string, string, string) {
	text, _ := bufio.NewReader(file).ReadString('\n')
	refsAndShas := strings.Split(strings.Trim(string(text), "\n"), " ")
	if len(refsAndShas) < 4 {
		return EmptySha, EmptySha, "", ""
	}
	return refsAndShas[0], refsAndShas[1], refsAndShas[2], refsAndShas[3]
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
