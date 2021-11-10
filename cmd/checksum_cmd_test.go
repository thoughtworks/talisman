package main

import (
	"fmt"
	"os"
	"talisman/git_testing"
	mock "talisman/internal/mock/utility"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestChecksumCalculatorShouldExitSuccess(t *testing.T) {
	withNewTmpGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file.txt")
		git.CreateFileWithContents("private.pem", "secret")
		git.CreateFileWithContents("another/private.pem", "secret")
		git.CreateFileWithContents("sample.txt", "password")
		os.Chdir(git.GetRoot())

		checksumCmd := NewChecksumCmd([]string{"*.txt"})
		assert.Equal(t, 0, checksumCmd.Run(), "Expected run() to return 0 as given patterns are found and .talsimanrc is suggested")
		options.Checksum = ""
	})
}

func TestChecksumCalculatorShouldExitFailure(t *testing.T) {
	withNewTmpGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file.txt")
		git.CreateFileWithContents("private.pem", "secret")
		git.CreateFileWithContents("another/private.pem", "secret")
		git.CreateFileWithContents("sample.txt", "password")
		os.Chdir(git.GetRoot())

		checksumCmd := NewChecksumCmd([]string{"*.java"})
		assert.Equal(t, 1, checksumCmd.Run(), "Expected run() to return 1 as given patterns are found and .talsimanrc is suggested")
		options.Checksum = ""
	})
}

func TestChecksumCalculatorShouldExitFailureWhenHasherStarFails(t *testing.T) {
	withNewTmpGitRepo(func(git *git_testing.GitTesting) {
		ctrl := gomock.NewController(t)
		hasher := mock.NewMockSHA256Hasher(ctrl)
		checksumCmd := ChecksumCmd{[]string{"*.java"}, hasher, git.GetRoot()}
		hasher.EXPECT().Start().Return(fmt.Errorf("fail this test because hasher failed to start"))
		assert.Equal(t, 1, checksumCmd.Run(), "Expected run() to return 1 because hasher failed to start")
	})
}
