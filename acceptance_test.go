package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	git "github.com/thoughtworks/talisman/git_testing"
)

const awsAccessKeyIDExample string = "accessKey=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

func TestNotHavingAnyOutgoingChangesShouldNotFail(t *testing.T) {
	withNewTmpGitRepo(func(gitPath string) {
		git.SetupBaselineFiles(gitPath, "simple-file")
		assert.Equal(t, 0, runTalisman(gitPath), "Expected run() to return 0 if no input is available on stdin. This happens when there are no outgoing changes")
	})
}

func TestAddingSimpleFileShouldExitZero(t *testing.T) {
	withNewTmpGitRepo(func(gitPath string) {
		git.SetupBaselineFiles(gitPath, "simple-file")
		exitStatus := runTalisman(gitPath)
		assert.Equal(t, 0, exitStatus, "Expected run() to return 0 and pass as no suspicious files are in the repo")
	})
}

func TestAddingSecretKeyShouldExitOne(t *testing.T) {
	withNewTmpGitRepo(func(gitPath string) {
		git.SetupBaselineFiles(gitPath, "simple-file")
		git.CreateFileWithContents(gitPath, "private.pem", "secret")
		git.AddAndcommit(gitPath, "*", "add private key")

		exitStatus := runTalisman(gitPath)
		assert.Equal(t, 1, exitStatus, "Expected run() to return 1 and fail as pem file was present in the repo")
	})
}

func TestAddingSecretKeyAsFileContentShouldExitOne(t *testing.T) {

	withNewTmpGitRepo(func(gitPath string) {
		git.SetupBaselineFiles(gitPath, "simple-file")
		git.CreateFileWithContents(gitPath, "contains_keys.properties", awsAccessKeyIDExample)
		git.AddAndcommit(gitPath, "*", "add private key as content")

		exitStatus := runTalisman(gitPath)
		assert.Equal(t, 1, exitStatus, "Expected run() to return 1 and fail as file contains some secrets")
	})
}

func TestAddingSecretKeyShouldExitZeroIfPEMFilesAreIgnored(t *testing.T) {
	withNewTmpGitRepo(func(gitPath string) {
		git.SetupBaselineFiles(gitPath, "simple-file")
		git.CreateFileWithContents(gitPath, "private.pem", "secret")
		git.CreateFileWithContents(gitPath, ".talismanignore", "*.pem")
		git.AddAndcommit(gitPath, "*", "add private key")

		exitStatus := runTalisman(gitPath)
		assert.Equal(t, 0, exitStatus, "Expected run() to return 0 and pass as pem file was ignored")
	})
}

func TestAddingSecretKeyShouldExitZeroIfPEMFilesAreIgnoredAndCommented(t *testing.T) {
	withNewTmpGitRepo(func(gitPath string) {
		git.SetupBaselineFiles(gitPath, "simple-file")
		git.CreateFileWithContents(gitPath, "private.pem", "secret")
		git.CreateFileWithContents(gitPath, ".talismanignore", "*.pem # I know what I'm doing")
		git.AddAndcommit(gitPath, "*", "add private key")

		exitStatus := runTalisman(gitPath)
		assert.Equal(t, 0, exitStatus, "Expected run() to return 0 and pass as pem file was ignored")
	})
}

func TestAddingSecretKeyShouldExitOneIfTheyContainBadContentButOnlyFilenameDetectorWasIgnored(t *testing.T) {
	withNewTmpGitRepo(func(gitPath string) {
		git.SetupBaselineFiles(gitPath, "simple-file")
		git.CreateFileWithContents(gitPath, "private.pem", awsAccessKeyIDExample)
		git.CreateFileWithContents(gitPath, ".talismanignore", "*.pem # ignore:filename")
		git.AddAndcommit(gitPath, "*", "add private key")

		exitStatus := runTalisman(gitPath)
		assert.Equal(t, 1, exitStatus, "Expected run() to return 0 and pass as pem file was ignored")
	})
}


func TestStagingSecretKeyShouldExitOneWhenPreCommitFlagIsSet(t *testing.T) {
	withNewTmpGitRepo(func(gitPath string) {
		git.SetupBaselineFiles(gitPath, "simple-file")
		git.CreateFileWithContents(gitPath, "private.pem", "secret")
		git.Add(gitPath, "*")

		options := Options{
			debug:   false,
			githook: "pre-commit",
		}

		exitStatus := runTalismanWithOptions(gitPath, options)
		assert.Equal(t, 1, exitStatus, "Expected run() to return 1 and fail as pem file was present in the repo")
	})
}

func runTalisman(gitPath string) int {
	options := Options{
		debug:   false,
		githook: "pre-push",
	}
	return runTalismanWithOptions(gitPath, options)
}

func runTalismanWithOptions(gitPath string, options Options) int {
	os.Chdir(gitPath)
	return run(mockStdIn(git.EarliestCommit(gitPath), git.LatestCommit(gitPath)), options)
}

func withNewTmpGitRepo(gitOp func(gitPath string)) {
	WithNewTmpDirNamed("talisman-acceptance-test", func(gitPath string) {
		git.Init(gitPath)
		gitOp(gitPath)
	})
}

type DirOp func(dirName string)

func WithNewTmpDirNamed(dirName string, dop DirOp) {
	path, err := ioutil.TempDir(os.TempDir(), dirName)
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(path)
	dop(path)
}

func mockStdIn(oldSha string, newSha string) io.Reader {
	return strings.NewReader(fmt.Sprintf("master %s master %s\n", newSha, oldSha))
}
