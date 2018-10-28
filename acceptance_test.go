package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"talisman/git_testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const awsAccessKeyIDExample string = "accessKey=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

func init() {
	git_testing.Logger = logrus.WithField("Environment", "Debug")
	git_testing.Logger.Debug("Accetpance test started")
}

func TestNotHavingAnyOutgoingChangesShouldNotFail(t *testing.T) {
	withNewTmpGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file")
		assert.Equal(t, 0, runTalisman(git), "Expected run() to return 0 if no input is available on stdin. This happens when there are no outgoing changes")
	})
}

func TestAddingSimpleFileShouldExitZero(t *testing.T) {
	withNewTmpGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file")
		assert.Equal(t, 0, runTalisman(git), "Expected run() to return 0 and pass as no suspicious files are in the repo")
	})
}

func TestAddingSecretKeyShouldExitOne(t *testing.T) {
	withNewTmpGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("private.pem", "secret")
		git.AddAndcommit("*", "add private key")

		assert.Equal(t, 1, runTalisman(git), "Expected run() to return 1 and fail as pem file was present in the repo")
	})
}

func TestAddingSecretKeyAsFileContentShouldExitOne(t *testing.T) {

	withNewTmpGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("contains_keys.properties", awsAccessKeyIDExample)
		git.AddAndcommit("*", "add private key as content")

		assert.Equal(t, 1, runTalisman(git), "Expected run() to return 1 and fail as file contains some secrets")
	})
}

func TestAddingSecretKeyShouldExitZeroIfPEMFilesAreIgnored(t *testing.T) {
	withNewTmpGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("private.pem", "secret")
		git.CreateFileWithContents(".talismanignore", "*.pem")
		git.AddAndcommit("*", "add private key")

		assert.Equal(t, 0, runTalisman(git), "Expected run() to return 0 and pass as pem file was ignored")
	})
}

func TestAddingSecretKeyShouldExitZeroIfPEMFilesAreIgnoredAndCommented(t *testing.T) {
	withNewTmpGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("private.pem", "secret")
		git.CreateFileWithContents(".talismanignore", "*.pem # I know what I'm doing")
		git.AddAndcommit("*", "add private key")

		assert.Equal(t, 0, runTalisman(git), "Expected run() to return 0 and pass as pem file was ignored")
	})
}

func TestAddingSecretKeyShouldExitOneIfTheyContainBadContentButOnlyFilenameDetectorWasIgnored(t *testing.T) {
	withNewTmpGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("private.pem", awsAccessKeyIDExample)
		git.CreateFileWithContents(".talismanignore", "*.pem # ignore:filename")
		git.AddAndcommit("*", "add private key")

		assert.Equal(t, 1, runTalisman(git), "Expected run() to return 0 and pass as pem file was ignored")
	})
}

func TestStagingSecretKeyShouldExitOneWhenPreCommitFlagIsSet(t *testing.T) {
	withNewTmpGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("private.pem", "secret")
		git.Add("*")

		_options := options{
			debug:   false,
			githook: PreCommit,
		}

		assert.Equal(t, 1, runTalismanWithOptions(git, _options), "Expected run() to return 1 and fail as pem file was present in the repo")
	})
}

func runTalisman(git *git_testing.GitTesting) int {
	_options := options{
		debug:   false,
		githook: "pre-push",
	}
	return runTalismanWithOptions(git, _options)
}

func runTalismanWithOptions(git *git_testing.GitTesting, _options options) int {
	wd, _ := os.Getwd()
	os.Chdir(git.GetRoot())
	defer func() { os.Chdir(wd) }()
	return run(mockStdIn(git.EarliestCommit(), git.LatestCommit()), _options)
}

type Operation func(dirName string)

func withNewTmpDirNamed(dirName string, operation Operation) {
	path, err := ioutil.TempDir(os.TempDir(), dirName)
	if err != nil {
		panic(err)
	}
	operation(path)
}

type GitOperation func(*git_testing.GitTesting)

func withNewTmpGitRepo(doGitOperation GitOperation) {
	withNewTmpDirNamed("talisman-acceptance-test", func(gitPath string) {
		gt := git_testing.Init(gitPath)
		doGitOperation(gt)
		os.RemoveAll(gitPath)
	})
}

func mockStdIn(oldSha string, newSha string) io.Reader {
	return strings.NewReader(fmt.Sprintf("master %s master %s\n", newSha, oldSha))
}
