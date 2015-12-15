package git_repo

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/badrij/go-codereview/patch"
)

//FilePath represents the absolute path of an added file
type FilePath string

//FileName represents the base name of an added file
type FileName string

//Addition represents the end state of a file
type Addition struct {
	Path FilePath
	Name FileName
	Data []byte
}

type gitRepo struct {
	root string
}

func RepoLocatedAt(path string) gitRepo {
	absoluteRoot, _ := filepath.Abs(path)
	return gitRepo{absoluteRoot}
}

func (repo gitRepo) AllChanges() []*patch.File {
	return repo.generateDiff("origin/master", "master").File
}

func (repo gitRepo) AllAdditions() []Addition {
	return repo.Additions("origin/master", "master")
}

func (repo gitRepo) Additions(oldCommit string, newCommit string) []Addition {
	files := repo.outgoingNonDeletedFiles(oldCommit, newCommit)
	result := make([]Addition, len(files))
	for i, file := range files {
		data, _ := repo.ReadRepoFile(file.Dst)
		result[i] = NewAddition(file.Dst, data)
	}
	log.WithFields(log.Fields{
		"oldCommit": oldCommit,
		"newCommit": newCommit,
		"additions": result,
	}).Info("Generating all additions in range.")
	return result
}

func NewAddition(filePath string, content []byte) Addition {
	return Addition{
		Path: FilePath(filePath),
		Name: FileName(path.Base(filePath)),
		Data: content,
	}
}

//ReadRepoFile returns the contents of the supplied relative filename by locating it in the git repo
func (repo gitRepo) ReadRepoFile(fileName string) ([]byte, error) {
	return ioutil.ReadFile(path.Join(repo.root, fileName))
}

//ReadRepoFileOrNothing returns the contents of the supplied relative filename by locating it in the git repo.
//If the given file cannot be located in theb repo, then an empty array of bytes is returned for the content.
func (repo gitRepo) ReadRepoFileOrNothing(fileName string) ([]byte, error) {
	filepath := path.Join(repo.root, fileName)
	if _, err := os.Stat(filepath); err == nil {
		return repo.ReadRepoFile(fileName)
	} else {
		return make([]byte, 0), nil
	}
}

//Matches states whether the addition matches the given pattern.
//If the pattern ends in a path separator, then all files inside a directory with that name are matched. However, files with that name itself will not be matched.
//If a pattern contains the path separator in any other location, the match works according to the pattern logic of the default golang glob mechanism
//If there is no path separator anywhere in the pattern, the pattern is matched against the base name of the file. Thus, the pattern will match files with that name anywhere in the repository.
func (a Addition) Matches(pattern string) bool {
	var result bool
	var err error
	if pattern[len(pattern)-1] == os.PathSeparator {
		result, err = strings.HasPrefix(string(a.Path), pattern), nil
	} else if strings.ContainsRune(pattern, os.PathSeparator) {
		result, err = path.Match(pattern, string(a.Path))
	} else {
		result, err = path.Match(pattern, string(a.Name))
	}
	log.WithFields(log.Fields{
		"pattern":  pattern,
		"filePath": a.Path,
		"match":    result,
		"error":    err,
	}).Debug("Checking addition for match.")
	check(err)
	return result
}

func (fn FileName) Matches(r *regexp.Regexp) bool {
	return r.MatchString(string(fn))
}

func (fn FileName) String() string {
	return string(fn)
}

func (fp FilePath) String() string {
	return string(fp)
}

func (a Addition) String() string {
	return fmt.Sprintf("%s\n\n%s\n", a.Path, string(a.Data))
}

func (repo gitRepo) generateDiff(oldCommit string, newCommit string) *patch.Set {
	patchSet, err := patch.Parse(repo.fetchRawOutgoingDiff(oldCommit, newCommit))
	check(err)
	return patchSet
}

func (repo gitRepo) outgoingNonDeletedFiles(oldCommit, newCommit string) []*patch.File {
	patchSet, err := patch.Parse(repo.fetchRawOutgoingDiff(oldCommit, newCommit))
	check(err)
	result := make([]*patch.File, 0)
	for _, file := range patchSet.File {
		if file.Verb != patch.Delete {
			result = append(result, file)
		}
	}
	return result
}

func (repo gitRepo) fetchRawOutgoingDiff(oldCommit string, newCommit string) []byte {
	gitRange := oldCommit + ".." + newCommit
	o, err := repo.command("git", "diff", gitRange, "--binary").Output()
	log.WithFields(log.Fields{
		"output": string(o),
		"error":  err,
	}).Debug("Raw outgoing diff.")
	check(err)
	return o
}

func (repo gitRepo) command(commandName string, args ...string) *exec.Cmd {
	log.WithFields(log.Fields{
		"command": commandName,
		"args":    args,
	}).Debug("Building repo command.")
	result := exec.Command(commandName, args...)
	result.Dir = repo.root
	return result
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
