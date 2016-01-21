package git_testing

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/drhodes/golorem"
)

func Init(gitRoot string) {
	os.MkdirAll(gitRoot, 0777)
	ExecCommand(gitRoot, "git", "init", ".")
}

func GitClone(existingRepoRoot string, cloneName string) string {
	ExecCommand("", "git", "clone", existingRepoRoot, cloneName)
	return cloneName
}

func SetupBaselineFiles(gitRoot string, filenames ...string) {
	for _, filename := range filenames {
		CreateFileWithContents(gitRoot, filename, lorem.Sentence(8, 10), lorem.Sentence(8, 10))
	}
	AddAndcommit(gitRoot, "*", "initial commit")
}

func EarliestCommit(gitRoot string) string {
	return ExecCommand(gitRoot, "git", "rev-list", "--max-parents=0", "HEAD")
}

func LatestCommit(gitRoot string) string {
	return ExecCommand(gitRoot, "git", "rev-parse", "HEAD")
}

func CreateFileWithContents(gitRoot string, name string, contents ...string) string {
	fileName := path.Join(gitRoot, name)
	os.MkdirAll(path.Dir(fileName), 0777)
	f, err := os.Create(fileName)
	die(err)
	defer f.Close()
	for _, line := range contents {
		f.WriteString(line)
	}
	return name
}

func OverwriteFileContent(gitRoot string, name string, contents ...string) {
	fileName := path.Join(gitRoot, name)
	f, err := os.Create(fileName)
	die(err)
	defer f.Close()
	for _, line := range contents {
		f.WriteString(line)
	}
	f.Sync()
}

func AppendFileContent(gitRoot string, name string, contents ...string) {
	fileName := path.Join(gitRoot, name)
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0660)
	die(err)
	defer f.Close()
	for _, line := range contents {
		f.WriteString(line)
	}
	f.Sync()
}

func RemoveFile(gitRoot string, filename string) {
	os.Remove(path.Join(gitRoot, filename))
}

func FileContents(gitRoot string, name string) []byte {
	fileName := path.Join(gitRoot, name)
	result, err := ioutil.ReadFile(fileName)
	die(err)
	return result
}

func AddAndcommit(gitRoot string, fileName string, message string) {
	ExecCommand(gitRoot, "git", "add", fileName)
	ExecCommand(gitRoot, "git", "commit", fileName, "-m", message)
}

func Add(gitRoot string, fileName string) {
	ExecCommand(gitRoot, "git", "add", fileName)
}

func ExecCommand(gitRoot string, commandName string, args ...string) string {
	result := exec.Command(commandName, args...)
	if len(gitRoot) > 0 {
		result.Dir = gitRoot
	}
	o, err := result.Output()
	die(err)
	return strings.Trim(string(o), "\n")
}

func die(err error) {
	if err != nil {
		panic(err)
	}
}
