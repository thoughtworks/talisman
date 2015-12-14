package git_testing

import (
	"github.com/drhodes/golorem"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

func Init(gitRoot string) {
	os.MkdirAll(gitRoot, 0777)
	check(Command(gitRoot, "git", "init", ".").Run())
}

func GitClone(existingRepoRoot string, cloneName string) string {
	check(exec.Command("git", "clone", existingRepoRoot, cloneName).Run())
	return cloneName
}

func SetupBaselineFiles(gitRoot string, filenames ...string) {
	for _, filename := range filenames {
		CreateFileWithContents(gitRoot, filename, lorem.Sentence(8, 10), lorem.Sentence(8, 10))
	}
	AddAndcommit(gitRoot, "*", "initial commit")
}

func EarliestCommit(gitRoot string) string {
	commit, err := Command(gitRoot, "git", "rev-list", "--max-parents=0", "HEAD").Output()
	check(err)
	return strings.Trim(string(commit), "\n")
}

func LatestCommit(gitRoot string) string {
	latestShaOutput, err := Command(gitRoot, "git", "rev-parse", "HEAD").Output()
	check(err)
	return strings.Trim(string(latestShaOutput), "\n")
}

func CreateFileWithContents(gitRoot string, name string, contents ...string) string {
	fileName := path.Join(gitRoot, name)
	os.MkdirAll(path.Dir(fileName), 0777)
	f, err := os.Create(fileName)
	check(err)
	defer f.Close()
	for _, line := range contents {
		f.WriteString(line)
	}
	return name
}

func AppendFileContent(gitRoot string, name string, contents ...string) {
	fileName := path.Join(gitRoot, name)
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0660)
	check(err)
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
	check(err)
	return result
}

func AddAndcommit(gitRoot string, fileName string, message string) {
	check(Command(gitRoot, "git", "add", fileName).Run())
	check(Command(gitRoot, "git", "commit", fileName, "-m", message).Run())
}

func Command(gitRoot string, commandName string, args ...string) *exec.Cmd {
	result := exec.Command(commandName, args...)
	result.Dir = gitRoot
	return result
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
