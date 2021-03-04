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
var gitConfigFile string

type GitTesting struct {
	gitRoot string
}

func Init(gitRoot string) *GitTesting {
	os.MkdirAll(gitRoot, 0777)
	testingRepo := &GitTesting{gitRoot}
	testingRepo.ExecCommand("git", "init", ".")
	gitConfigFileObject, _ := ioutil.TempFile(os.TempDir(), "gitConfigForTalismanTests")
	gitConfigFile = gitConfigFileObject.Name()
	testingRepo.CreateFileWithContents(gitConfigFile, `[user]
	email = talisman-test-user@example.com
	name = Talisman Test User`)
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

func (git *GitTesting) GetBlobDetails(fileName string) string {
	var output []byte
	object_hash_and_filename := ""
	git.doInGitRoot(func() {
		fmt.Println("hello")
		result := exec.Command("git", "rev-list", "--objects", "--all")
		output, _ = result.Output()
		objects := strings.Split(string(output), "\n")
		for _, object := range objects {
			object_details := strings.Split(object, " ")
			if len(object_details) == 2 && object_details[1] == fileName {
				object_hash_and_filename = object
				return
			}
		}
	})
	return object_hash_and_filename
}


//ExecCommand executes a command with given arguments in the git repo directory
func (git *GitTesting) ExecCommand(commandName string, args ...string) string {
	var output []byte
	git.doInGitRoot(func() {
		result := exec.Command(commandName, args...)
		//Passes locally, but fails on CI
		//result.Env = []string{"GIT_CONFIG=" + gitConfigFile}
		var err error
		output, err = result.Output()
		fmt.Println("------------------>>>>>>>>>")
		fmt.Printf("when executing command %s %v in %s\nError: %v\n", commandName, args, git.gitRoot, err)
		git.die(fmt.Sprintf("when executing command %s %v in %s\nError: %v", commandName, args, git.gitRoot, err), err)
		Logger.Debugf("Output of command %s %v in %s is: %s\n", commandName, args, git.gitRoot, string(output))
	})
	if len(output) > 0 {
		return strings.Trim(string(output), "\n")
	}
	return ""
}

func (git *GitTesting) die(msg string, err error) {
	if err != nil {
		msgtxt := fmt.Sprintf("Error %s: %s\n", msg, err.Error())
		Logger.Debugf(msgtxt)
		fmt.Println(msgtxt);
		panic(msgtxt)
	}
}

func (git *GitTesting) doInGitRoot(operation func()) {
	wd, _ := os.Getwd()
	os.Chdir(git.gitRoot)
	defer func() { os.Chdir(wd) }()
	operation()
}

func (git *GitTesting) GetRoot() string {
	return git.gitRoot
}

func (git *GitTesting) RemoveHooks() {
	git.ExecCommand("rm", "-rf", ".git/hooks/")
}
