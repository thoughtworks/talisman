package main

import (
	"os"
	"talisman/git_testing"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChecksumCalculatorShouldExitSuccess(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file.txt")
		git.CreateFileWithContents("private.pem", "secret")
		git.CreateFileWithContents("another/private.pem", "secret")
		git.CreateFileWithContents("sample.txt", "password")
		os.Chdir(git.Root())

		checksumCmd := NewChecksumCmd([]string{"*.txt"})
		assert.Equal(t, 0, checksumCmd.Run(), "Expected run() to return 0 as given patterns are found and .talsimanrc is suggested")
		options.Checksum = ""
	})
}

func TestChecksumCalculatorShouldExitFailure(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file.txt")
		git.CreateFileWithContents("private.pem", "secret")
		git.CreateFileWithContents("another/private.pem", "secret")
		git.CreateFileWithContents("sample.txt", "password")
		os.Chdir(git.Root())

		checksumCmd := NewChecksumCmd([]string{"*.java"})
		assert.Equal(t, 1, checksumCmd.Run(), "Expected run() to return 1 as given patterns are found and .talsimanrc is suggested")
		options.Checksum = ""
	})
}

func TestChecksumCalculatorShouldExitFailureWhenHasherIsEmpty(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		checksumCmd := ChecksumCmd{[]string{"*.java"}, nil, git.Root()}
		assert.Equal(t, 1, checksumCmd.Run(), "Expected run() to return 1 because hasher failed to start")
	})
}
