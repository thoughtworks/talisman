package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	git "github.com/thoughtworks/talisman/git_testing"
)

func TestAddingSecretKeyShouldExitOne(t *testing.T) {
	// Given
	gitPath, err := ioutil.TempDir(os.TempDir(), "talisman-acceptance-test")
	check(err)
	defer os.RemoveAll(gitPath)

	git.Init(gitPath)
	git.SetupBaselineFiles(gitPath, "simple-file")
	git.CreateFileWithContents(gitPath, "private.pem", "secret")
	git.AddAndcommit(gitPath, "*", "add private key")

	os.Chdir(gitPath)
	stdIn := formatStdIn(git.EarliestCommit(gitPath), git.LatestCommit(gitPath))

	// Then
	assert.Equal(t, 1, run(strings.NewReader(stdIn)), "Expected run() to return 1 and fail as pem file was present in the repo")
}

func TestAddingSecretKeyShouldExitZeroIfPEMFilesAreIgnored(t *testing.T) {
	// Given
	gitPath, _ := ioutil.TempDir(os.TempDir(), "talisman-acceptance-test")
	defer os.RemoveAll(gitPath)

	git.Init(gitPath)
	git.SetupBaselineFiles(gitPath, "simple-file")
	git.CreateFileWithContents(gitPath, "private.pem", "secret")
	git.CreateFileWithContents(gitPath, ".talismanignore", "*.pem")
	git.AddAndcommit(gitPath, "*", "add private key")

	os.Chdir(gitPath)
	stdIn := formatStdIn(git.EarliestCommit(gitPath), git.LatestCommit(gitPath))

	// Then
	assert.Equal(t, 0, run(strings.NewReader(stdIn)), "Expected run() to return 0 and pass as pem file was ignored")
}

func TestAddingSimpleFileShouldExitZero(t *testing.T) {
	// Given
	gitPath, err := ioutil.TempDir(os.TempDir(), "talisman-acceptance-test")
	check(err)
	defer os.RemoveAll(gitPath)

	git.Init(gitPath)
	git.SetupBaselineFiles(gitPath, "simple-file")

	os.Chdir(gitPath)
	stdIn := formatStdIn(git.EarliestCommit(gitPath), git.LatestCommit(gitPath))

	assert.Equal(t, 0, run(strings.NewReader(stdIn)), "Expected run() to return 0 and pass as no suspicious files are in the repo")
}

func TestNotHavingAnyOutgoingChangesShouldNotFail(t *testing.T) {
	assert.Equal(t, 0, run(strings.NewReader("")), "Expected run() to return 0 if no input is available on stdin. This happens when there are no outgoing changes")
}

func formatStdIn(oldSha string, newSha string) string {
	return fmt.Sprintf("master %s master %s\n", newSha, oldSha)
}
