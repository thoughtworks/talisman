package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"talisman/git_repo"

	log "github.com/Sirupsen/logrus"
)

var (
	fdebug      bool
	githook     string
	showVersion bool
	pattern     string
	//Version : Version of talisman
	Version      = "Development Build"
	blob_details string
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
	debug        bool
	githook      string
	pattern      string
	blob_details string
}

//Logger is the default log device, set to emit at the Error level by default
func main() {
	flag.BoolVar(&fdebug, "d", false, "short form of debug")
	flag.BoolVar(&fdebug, "debug", false, "enable debug mode (warning: very verbose)")
	flag.BoolVar(&showVersion, "v", false, "short form of version")
	flag.BoolVar(&showVersion, "version", false, "show current version of talisman")
	flag.StringVar(&pattern, "p", "", "short form of pattern")
	flag.StringVar(&pattern, "pattern", "", "pattern (glob-like) of files to scan (ignores githooks)")
	flag.StringVar(&githook, "githook", PrePush, "either pre-push or pre-commit")
	flag.StringVar(&blob_details, "blob", "", "blob details for scanner")

	flag.Parse()

	if showVersion {
		fmt.Printf("talisman %s\n", Version)
		os.Exit(0)
	}

	_options := options{
		debug:        fdebug,
		githook:      githook,
		pattern:      pattern,
		blob_details: blob_details,
	}

	os.Exit(run(os.Stdin, _options))
}

func run(stdin io.Reader, _options options) (returnCode int) {
	if _options.debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.ErrorLevel)
	}

	if _options.blob_details != "" {
		var additions []git_repo.Addition
		return NewRunner(additions).Scan(_options.blob_details)
	}

	if _options.githook == "" {
		_options.githook = PrePush
	}

	var additions []git_repo.Addition
	if _options.pattern != "" {
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
