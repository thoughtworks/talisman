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
	//Version : Version of talisman
	Version = "Development Build"
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
	debug   bool
	githook string
}

//Logger is the default log device, set to emit at the Error level by default
func main() {
	flag.BoolVar(&fdebug, "debug", false, "enable debug mode (warning: very verbose)")
	flag.BoolVar(&fdebug, "d", false, "short form of debug (warning: very verbose)")
	flag.BoolVar(&showVersion, "v", false, "show current version of talisman")
	flag.BoolVar(&showVersion, "version", false, "show current version of talisman")
	flag.StringVar(&githook, "githook", PrePush, "either pre-push or pre-commit")
	flag.Parse()

	if showVersion {
		fmt.Printf("talisman %s\n", Version)
		os.Exit(0)
	}

	_options := options{
		debug:   fdebug,
		githook: githook,
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

	log.Infof("Running %s hook", _options.githook)

	var additions []git_repo.Addition
	if _options.githook == PreCommit {
		preCommitHook := NewPreCommitHook()
		additions = preCommitHook.GetRepoAdditions()
	} else {
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
