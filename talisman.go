package main

import flag "github.com/spf13/pflag"

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
	fdebug      bool
	githook     string
	showVersion bool
	pattern     string
	//Version : Version of talisman
	Version         = "Development Build"
	scan            bool
	checksum        string
	reportdirectory string
	scanWithHtml    bool
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
	checksum        string
	reportdirectory string
	scanWithHtml    bool
}

//Logger is the default log device, set to emit at the Error level by default
func main() {
	flag.BoolVarP(&fdebug, "debug", "d", false, "enable debug mode (warning: very verbose)")
	flag.BoolVarP(&showVersion, "version", "v", false, "show current version of talisman")
	flag.StringVarP(&pattern, "pattern", "p", "", "pattern (glob-like) of files to scan (ignores githooks)")
	flag.StringVarP(&githook, "githook", "g", PrePush, "either pre-push or pre-commit")
	flag.BoolVarP(&scan, "scan", "s", false, "scanner scans the git commit history for potential secrets")
	flag.StringVarP(&checksum, "checksum", "c", "", "checksum calculator calculates checksum and suggests .talsimarc format")
	flag.StringVarP(&reportdirectory, "reportdirectory", "r", "", "directory where the scan reports will be stored")
	flag.BoolVarP(&scanWithHtml, "scanWithHtml", "w", false, "generate html report (**Make sure you have installed talisman_html_report to use this, as mentioned in Readme**)")

	flag.Parse()

	if showVersion {
		fmt.Printf("talisman %s\n", Version)
		os.Exit(0)
	}

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	_options := options{
		debug:           fdebug,
		githook:         githook,
		pattern:         pattern,
		scan:            scan,
		checksum:        checksum,
		reportdirectory: reportdirectory,
		scanWithHtml:    scanWithHtml,
	}

	os.Exit(run(os.Stdin, _options))
}

func run(stdin io.Reader, _options options) (returnCode int) {
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
		return NewRunner(make([]gitrepo.Addition, 0)).Scan(_options.reportdirectory)
	} else if _options.scanWithHtml {
		log.Infof("Running scanner with html report")
		return NewRunner(make([]gitrepo.Addition, 0)).Scan("talisman_html_report")
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

	return NewRunner(additions).RunWithoutErrors()
}

func readRefAndSha(file io.Reader) (string, string, string, string) {
	text, _ := bufio.NewReader(file).ReadString('\n')
	refsAndShas := strings.Split(strings.Trim(string(text), "\n"), " ")
	if len(refsAndShas) < 4 {
		return EmptySha, EmptySha, "", ""
	}
	return refsAndShas[0], refsAndShas[1], refsAndShas[2], refsAndShas[3]
}
