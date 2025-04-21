package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"talisman/prompt"
	"testing"

	"talisman/git_testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const awsAccessKeyIDExample string = "accessKey=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

const talismanRCDataWithIgnoreDetectorWithFilename = `
fileignoreconfig:
- filename: private.pem
  checksum: 05db785bf1e1712f69b81eeb9956bd797b956e7179ebe3cb7bb2cd9be37a24c
  ignore_detectors: [filename]
`

const talismanRCDataWithScopeAsGo = `
scopeconfig:
 - scope: go
`

const talismanRCDataWithIgnoreDetectorWithFilecontent = `
fileignoreconfig:
- filename: private.pem
  checksum: 05db785bf1e1712f69b81eeb9956bd797b956e7179ebe3cb7bb2cd9be37a24c
  ignore_detectors: [filecontent]
`

const talismanRCDataWithFileNameAndCorrectChecksum = `
fileignoreconfig:
- filename: private.pem
  checksum: 1db800b79e6e9695adc451f77be974dc47bcd84d42873560d7767bfca30db8b1
  ignore_detectors: []
`
const invalidTalismanRC = `
fileignoreconfig:
- filename:
private.pem
  checksum: checksum_value
  ignore_detectors: []
  `

const talismanRCForHelloTxtFile = `
fileignoreconfig:
- filename: hello.txt
  checksum: edf37f1d33525d1710c3f5dc3437778c483def1d2e98b5a3495fc627eb1be588
  ignore_detectors: []
`

func init() {
	git_testing.Logger = logrus.WithField("Environment", "Debug")
	git_testing.Logger.Debug("Acceptance test started")
}

func TestNotHavingAnyOutgoingChangesShouldNotFail(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file")
		assert.Equal(t, 0, runTalismanInPrePushMode(git), "Expected run() to return 0 if no input is available on stdin. This happens when there are no outgoing changes")
	})
}

func TestAddingSimpleFileShouldExitZero(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file")
		assert.Equal(t, 0, runTalismanInPrePushMode(git), "Expected run() to return 0 and pass as no suspicious files are in the repo")
	})
}

func TestAddingSecretKeyShouldExitOne(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("private.pem", "secret")
		git.AddAndcommit("*", "add private key")

		assert.Equal(t, 1, runTalismanInPrePushMode(git), "Expected run() to return 1 and fail as pem file was present in the repo")
	})
}

func TestAddingSecretKeyAsFileContentShouldExitOne(t *testing.T) {

	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("contains_keys.properties", awsAccessKeyIDExample)
		git.AddAndcommit("*", "add private key as content")

		assert.Equal(t, 1, runTalismanInPrePushMode(git), "Expected run() to return 1 and fail as file contains some secrets")
	})
}

func TestAddingSecretKeyShouldExitZeroIfPEMFileIsIgnored(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("private.pem", "secret")
		git.CreateFileWithContents(".talismanrc", talismanRCDataWithFileNameAndCorrectChecksum)
		git.AddAndcommit("private.pem", "add private key")

		assert.Equal(t, 0, runTalismanInPrePushMode(git), "Expected run() to return 0 and pass as pem file was ignored")
	})
}

func TestScanningSimpleFileShouldExitZero(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		options.Scan = false

		git.SetupBaselineFiles("simple-file")
		assert.Equal(t, 0, runTalismanInPrePushMode(git), "Expected run() to return 0 and pass as pem file was ignored")
	})
}

func TestChecksumCalculatorShouldExitOne(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		options.Checksum = "*txt1"
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("private.pem", "secret")
		git.CreateFileWithContents("another/private.pem", "secret")
		git.CreateFileWithContents("sample.txt", "password")
		assert.Equal(t, 1, runTalisman(git), "Expected run() to return 0 as given patterns are found and .talsimanrc is suggested")
		options.Checksum = ""
	})
}

func TestShouldExitOneWhenSecretIsCommitted(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		options.Debug = false
		options.GitHook = PreCommit
		options.Scan = false
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("sample.txt", "password=somepassword \n")
		git.Add("*")
		assert.Equal(t, 1, runTalisman(git), "Expected run() to return 1 as given patterns are found")
	})
}

func TestShouldExitZeroWhenNonSecretIsCommittedButFileContainsSecretPreviously(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		options.Debug = false
		options.GitHook = PreCommit
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("sample.txt", "password=somepassword \n")
		git.AddAndcommit("*", "Initial Commit With Secret")

		git.AppendFileContent("sample.txt", "some text \n")
		git.Add("*")

		assert.Equal(t, 0, runTalisman(git), "Expected run() to return 1 as given patterns are not found")
	})
}

// Need to work on this test case as talismanrc does  not yet support comments
func TestAddingSecretKeyShouldExitZeroIfPEMFilesAreIgnoredAndCommented(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("private.pem", "secret")
		git.CreateFileWithContents(".talismanrc", talismanRCDataWithIgnoreDetectorWithFilename)
		git.AddAndcommit("*", "add private key")

		assert.Equal(t, 0, runTalismanInPrePushMode(git), "Expected run() to return 0 and pass as pem file was ignored")
	})
}

func TestAddingSecretKeyShouldExitOneIfTheyContainBadContentButOnlyFilenameDetectorWasIgnored(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("private.pem", awsAccessKeyIDExample)
		git.CreateFileWithContents(".talismanrc", talismanRCDataWithIgnoreDetectorWithFilename)
		git.AddAndcommit("private.pem", "add private key")

		assert.Equal(t, 1, runTalismanInPrePushMode(git), "Expected run() to return 1 and fail as only filename was ignored")
	})
}

func TestAddingSecretKeyShouldExitZeroIfFileIsWithinConfiguredScope(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("glide.lock", awsAccessKeyIDExample)
		git.CreateFileWithContents("glide.yaml", awsAccessKeyIDExample)
		git.CreateFileWithContents(".talismanrc", talismanRCDataWithScopeAsGo)
		git.AddAndcommit("*", "add private key")

		assert.Equal(t, 0, runTalismanInPrePushMode(git), "Expected run() to return 1 and fail as only filename was ignored")
	})
}

func TestAddingSecretKeyShouldExitOneIfFileIsNotWithinConfiguredScope(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("danger.pem", awsAccessKeyIDExample)
		git.CreateFileWithContents("glide.yaml", awsAccessKeyIDExample)
		git.CreateFileWithContents(".talismanrc", talismanRCDataWithScopeAsGo)
		git.AddAndcommit("*", "add private key")

		assert.Equal(t, 1, runTalismanInPrePushMode(git), "Expected run() to return 1 and fail as only filename was ignored")
	})
}

func TestAddingSecretKeyShouldExitOneIfFileNameIsSensitiveButOnlyFilecontentDetectorWasIgnored(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("private.pem", awsAccessKeyIDExample)
		git.CreateFileWithContents(".talismanrc", talismanRCDataWithIgnoreDetectorWithFilecontent)
		git.AddAndcommit("private.pem", "add private key")

		assert.Equal(t, 1, runTalismanInPrePushMode(git), "Expected run() to return 1 and fail as only filename was ignored")
	})
}

func TestStagingSecretKeyShouldExitOneWhenPreCommitFlagIsSet(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("private.pem", "secret")
		git.Add("*")

		options.Debug = false
		options.GitHook = PreCommit

		assert.Equal(t, 1, runTalisman(git), "Expected run() to return 1 and fail as pem file was present in the repo")
	})
}

func TestPatternFindsSecretKey(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		options.Debug = false
		options.Pattern = "./*.*"

		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("private.pem", "secret")

		assert.Equal(t, 1, runTalisman(git), "Expected run() to return 1 and fail as pem file was present in the repo")
	})
}

func TestPatternFindsNestedSecretKey(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		options.Debug = false
		options.Pattern = "./**/*.*"

		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("some-dir/private.pem", "secret")

		assert.Equal(t, 1, runTalisman(git), "Expected run() to return 1 and fail as nested pem file was present in the repo")
	})
}

func TestPatternFindsSecretInNestedFile(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		options.Debug = false
		options.Pattern = "./**/*.*"

		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("some-dir/some-file.txt", awsAccessKeyIDExample)

		assert.Equal(t, 1, runTalisman(git), "Expected run() to return 1 and fail as nested pem file was present in the repo")
	})
}

func TestFilesWithSameNameWithinRepositoryAreHandledAsSeparateFiles(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file", "some-dir/hello.txt")
		git.CreateFileWithContents("hello.txt", awsAccessKeyIDExample)
		git.AddAndcommit("*", "Start of Scan before talismanrc")
		assert.Equal(t, 1, runTalismanInPrePushMode(git), "Expected run() to return 1 since secret is detected in hello")
		git.CreateFileWithContents(".talismanrc", talismanRCForHelloTxtFile)

		git.AppendFileContent("some-dir/hello.txt", "More safe content")
		assert.Equal(t, 0, runTalismanInPrePushMode(git), "Expected run() to return 0 since hello checksum is added to fileignoreconfig in talismanrc")
		git.AddAndcommit("some-dir/hello.txt", "Start of second scan")
		assert.Equal(t, 0, runTalismanInPrePushMode(git), "Expected run() to return 0 since hello in subfolder was only changed")
		git.AppendFileContent("hello.txt", "More safe content in unsafe file")
		git.AddAndcommit("hello.txt", "Start of third scan - new hash is required")
		assert.Equal(t, 1, runTalismanInPrePushMode(git), "Expected run() to return 1 since hello checksum is changed")
	})
}

func TestScan(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		options.Debug = false
		options.Pattern = "./**/*.*"
		options.Scan = true

		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("some-dir/should-not-be-included.txt", awsAccessKeyIDExample)
		git.AddAndcommit("*", "Initial Commit")
		git.CreateFileWithContents("some-dir/should-be-included.txt", "safeContents")
		git.AddAndcommit("*", "Start of Scan")
		git.RemoveFile("some-dir/should-not-be-included.txt")
		git.AddAndcommit("*", "Removed secret")

		t.Run("Detects removed secrets", func(t *testing.T) {
			options.IgnoreHistory = false
			assert.Equal(t, 1, runTalisman(git), "Expected run() to return 1 because of removed secret in history")
		})

		t.Run("Does not detect removed secrets when ignoring history", func(t *testing.T) {
			options.IgnoreHistory = true
			assert.Equal(t, 0, runTalisman(git), "Expected run() to return 0 since secret was removed from head")
		})
	})
}

func TestScanDetectsIgnoredSecret(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		options.Scan = true
		options.IgnoreHistory = false

		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("private.pem", "secret")
		git.CreateFileWithContents(".talismanrc", talismanRCDataWithFileNameAndCorrectChecksum)
		git.AddAndcommit("private.pem", "add private key")
		assert.Equal(t, 1, runTalisman(git), "Expected run() to return 1 because ignores aren't check when scanning history")
	})
}

func TestIgnoreHistoryDetectsExistingIssuesOnHead(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		options.Debug = false
		options.Pattern = "./**/*.*"
		options.Scan = true
		options.IgnoreHistory = true

		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents("some-dir/file-with-issue.txt", awsAccessKeyIDExample)
		git.AddAndcommit("*", "Commit with Secret")
		git.CreateFileWithContents("some-dir/should-be-included.txt", "safeContents")
		git.AddAndcommit("*", "Another Commit")

		assert.Equal(t, 1, runTalisman(git), "Expected run() to return 1 since secret exists on head")
	})
}

func TestTalismanFailsIfTalismanrcIsInvalidYamlInPrePushMode(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents(".talismanrc", invalidTalismanRC)
		git.AddAndcommit("*", "Incorrect Talismanrc commit")

		assert.Equal(t, 1, runTalismanInPrePushMode(git), "Expected run() to return 1 and fails as talismanrc is invalid")
	})
}

func TestTalismanFailsIfTalismanrcIsInvalidYamlInPreCommitMode(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		options.Debug = true
		options.GitHook = PreCommit
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents(".talismanrc", invalidTalismanRC)
		git.AddAndcommit("*", "Incorrect Talismanrc commit")

		assert.Equal(t, 1, runTalisman(git), "Expected run() to return 0 and pass as pem file was ignored")
	})
}

func TestTalismanFailsIfTalismanrcIsInvalidYamlInScanMode(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		options.Debug = true
		options.Scan = true
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents(".talismanrc", invalidTalismanRC)
		git.AddAndcommit("*", "Incorrect Talismanrc commit")

		assert.Equal(t, 1, runTalisman(git), "Expected run() to return 0 and pass as pem file was ignored")
	})
}

func TestTalismanFailsIfTalismanrcIsInvalidYamlInScanWithHTMLMode(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		options.Debug = true
		options.ScanWithHtml = true
		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents(".talismanrc", invalidTalismanRC)
		git.AddAndcommit("*", "Incorrect Talismanrc commit")

		assert.Equal(t, 1, runTalisman(git), "Expected run() to return 0 and pass as pem file was ignored")
	})
}

func TestTalismanFailsIfTalismanrcIsInvalidYamlInPatternMode(t *testing.T) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		options.Debug = false
		options.Pattern = "./*.*"

		git.SetupBaselineFiles("simple-file")
		git.CreateFileWithContents(".talismanrc", invalidTalismanRC)

		assert.Equal(t, 1, runTalisman(git), "Expected run() to return 1 and fail as pem file was present in the repo")
	})
}

func runTalismanInPrePushMode(git *git_testing.GitTesting) int {
	options.Debug = true
	options.GitHook = PrePush
	return runTalisman(git)
}

func runTalisman(git *git_testing.GitTesting) int {
	wd, _ := os.Getwd()
	os.Chdir(git.Root())
	defer func() { os.Chdir(wd) }()
	prompter := prompt.NewPrompt()
	promptContext := prompt.NewPromptContext(false, prompter)
	if options.GitHook == PrePush {
		talismanInput = mockStdIn(git.EarliestCommit(), git.LatestCommit())
	}
	return run(promptContext)
}

func mockStdIn(oldSha string, newSha string) io.Reader {
	return strings.NewReader(fmt.Sprintf("master %s master %s\n", newSha, oldSha))
}
