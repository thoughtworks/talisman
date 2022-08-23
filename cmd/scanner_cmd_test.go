package main

import (
	"os"
	"talisman/git_testing"
	"talisman/talismanrc"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScannerCmdRunsSuccessfully(t *testing.T) {
	withNewTmpGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("some-dir/should-be-included.txt", "safeContents")
		git.AddAndcommit("*", "Start of Scan")
		os.Chdir(git.GetRoot())

		scannerCmd := NewScannerCmd(true, git.GetRoot())
		scannerCmd.Run(&talismanrc.TalismanRC{})
		assert.Equal(t, 0, scannerCmd.exitStatus(), "Expected ScannerCmd.exitStatus() to return 0 since no secret is found")
	})
}

func TestScannerCmdDetectsSecretAndFails(t *testing.T) {
	withNewTmpGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("some-dir/file-with-secret.txt", awsAccessKeyIDExample)
		git.AddAndcommit("*", "Initial Commit")
		git.RemoveFile("some-dir/file-with-secret.txt")
		git.AddAndcommit("*", "Removed secret")
		git.CreateFileWithContents("some-dir/safe-file.txt", "safeContents")
		git.AddAndcommit("*", "Start of Scan")
		os.Chdir(git.GetRoot())

		scannerCmd := NewScannerCmd(false, git.GetRoot())
		scannerCmd.Run(&talismanrc.TalismanRC{})
		assert.Equal(t, 1, scannerCmd.exitStatus(), "Expected ScannerCmd.exitStatus() to return 1 since secret present in history")
	})
}

func TestScannerCmdAddingSecretKeyShouldExitZeroIfFileIsWithinConfiguredScope(t *testing.T) {
	withNewTmpGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("go.sum", awsAccessKeyIDExample)
		git.CreateFileWithContents("go.mod", awsAccessKeyIDExample)
		git.CreateFileWithContents(".talismanrc", talismanRCDataWithScopeAsGo)
		git.AddAndcommit("*", "go sum file")
		os.Chdir(git.GetRoot())

		scannerCmd := NewScannerCmd(false, git.GetRoot())
		scannerCmd.Run(&talismanrc.TalismanRC{ScopeConfig: []talismanrc.ScopeConfig{{ScopeName: "go"}}})
		assert.Equal(t, 0, scannerCmd.exitStatus(), "Expected ScannerCmd.exitStatus() to return 0 since no secret is found")
	})
}
