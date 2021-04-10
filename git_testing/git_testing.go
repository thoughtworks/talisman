package git_testing

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	lorem "github.com/drhodes/golorem"
)

var Logger *logrus.Entry

type GitTesting struct {
	gitRoot string
}

func Init(gitRoot string) *GitTesting {
	os.MkdirAll(gitRoot, 0777)
	testingRepo := &GitTesting{gitRoot}
	output := testingRepo.ExecCommand("git", "init", ".")
	logrus.Debugf("Git init result %v", string(output))
	if os.Getenv("CI") != "" {
		testingRepo.ExecCommand("git", "config", "--global", "user.email", "talisman-test-user@example.com")
		testingRepo.ExecCommand("git", "config", "--global", "user.name", "Talisman Test User")
	}
	return testingRepo
}

func (git *GitTesting) GitClone(cloneName string) *GitTesting {
	result := git.ExecCommand("git", "clone", git.gitRoot, cloneName)
	Logger.Debugf("Clone result : %s\n", result)
	Logger.Debugf("GitRoot :%s \t CloneRoot: %s\n", git.gitRoot, cloneName)
	retval := &GitTesting{cloneName}
	return retval
}

func (git *GitTesting) SetupBaselineFiles(filenames ...string) {
	Logger.Debugf("Creating %v in %s\n", filenames, git.gitRoot)
	git.doInGitRoot(func() {
		for _, filename := range filenames {
			git.CreateFileWithContents(filename, lorem.Sentence(8, 10), lorem.Sentence(8, 10))
		}
	})
	git.AddAndcommit("*", "initial commit")
}

func (git *GitTesting) EarliestCommit() string {
	return git.ExecCommand("git", "rev-list", "--max-parents=0", "HEAD")
}

func (git *GitTesting) LatestCommit() string {
	return git.ExecCommand("git", "rev-parse", "HEAD")
}

func (git *GitTesting) CreateFileWithContents(filePath string, contents ...string) string {
	git.doInGitRoot(func() {
		os.MkdirAll(filepath.Dir(filePath), 0700)
		f, err := os.Create(filePath)
		git.die(fmt.Sprintf("when creating file %s", filePath), err)
		defer f.Close()
		for _, line := range contents {
			f.WriteString(line)
		}
	})
	return filePath
}

func (git *GitTesting) OverwriteFileContent(filePath string, contents ...string) {
	git.doInGitRoot(func() {
		os.MkdirAll(filepath.Dir(filePath), 0770)
		f, err := os.Create(filePath)
		git.die(fmt.Sprintf("when overwriting file %s", filePath), err)
		defer f.Close()
		for _, line := range contents {
			f.WriteString(line)
		}
		f.Sync()
	})
}

func (git *GitTesting) AppendFileContent(filePath string, contents ...string) {
	git.doInGitRoot(func() {
		os.MkdirAll(filepath.Dir(filePath), 0770)
		f, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0660)
		git.die(fmt.Sprintf("when appending file %s", filePath), err)
		defer f.Close()
		for _, line := range contents {
			f.WriteString(line)
		}
		f.Sync()
	})
}

func (git *GitTesting) RemoveFile(filename string) {
	git.doInGitRoot(func() {
		os.Remove(filename)
	})
}

func (git *GitTesting) FileContents(filePath string) []byte {
	var result []byte
	var err error
	git.doInGitRoot(func() {
		result, err = ioutil.ReadFile(filePath)
		git.die(fmt.Sprintf("when reading file %s", filePath), err)
	})
	return result
}

func (git *GitTesting) AddAndcommit(fileName string, message string) {
	git.Add(fileName)
	git.Commit(fileName, message)
}

func (git *GitTesting) Add(fileName string) {
	git.ExecCommand("git", "add", fileName)
}

func (git *GitTesting) Commit(fileName string, message string) {
	git.ExecCommand("git", "commit", "-m", message)
}

//GetBlobDetails returns git blob details for a path
func (git *GitTesting) GetBlobDetails(fileName string) string {
	var output []byte
	objectHashAndFilename := ""
	git.doInGitRoot(func() {
		fmt.Println("hello")
		result := exec.Command("git", "rev-list", "--objects", "--all")
		output, _ = result.Output()
		objects := strings.Split(string(output), "\n")
		for _, object := range objects {
			objectDetails := strings.Split(object, " ")
			if len(objectDetails) == 2 && objectDetails[1] == fileName {
				objectHashAndFilename = object
				return
			}
		}
	})
	return objectHashAndFilename
}

//ExecCommand executes a command with given arguments in the git repo directory
func (git *GitTesting) ExecCommand(commandName string, args ...string) string {
	var output []byte
	git.doInGitRoot(func() {
		result := exec.Command(commandName, args...)
		var err error
		output, err = result.Output()
		summaryMessage := fmt.Sprintf("Command: %s %v\nWorkingDirectory: %s\nOutput %s\nError: %v", commandName, args, git.gitRoot, string(output), err)
		git.die(summaryMessage, err)
		Logger.Debug(summaryMessage)
	})
	if len(output) > 0 {
		return strings.Trim(string(output), "\n")
	}
	return ""
}

func (git *GitTesting) die(msg string, err error) {
	if err != nil {
		Logger.Debugf(msg)
		panic(msg)
	}
}

func (git *GitTesting) doInGitRoot(operation func()) {
	wd, _ := os.Getwd()
	os.Chdir(git.gitRoot)
	defer func() { os.Chdir(wd) }()
	operation()
}

//GetRoot returns the root directory of the git-testing repo
func (git *GitTesting) GetRoot() string {
	return git.gitRoot
}

//RemoveHooks removes all file-system hooks from git-test repo
func (git *GitTesting) RemoveHooks() {
	git.ExecCommand("rm", "-rf", ".git/hooks/")
}
