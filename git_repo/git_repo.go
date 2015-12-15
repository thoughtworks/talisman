package git_repo

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
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
	return repo.changedFiles("origin/master", "master")
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
	if pattern[len(pattern)-1] == os.PathSeparator {
		result = strings.HasPrefix(string(a.Path), pattern)
	} else if strings.ContainsRune(pattern, os.PathSeparator) {
		result, _ = path.Match(pattern, string(a.Path))
	} else {
		result, _ = path.Match(pattern, string(a.Name))
	}
	log.WithFields(log.Fields{
		"pattern":  pattern,
		"filePath": a.Path,
		"match":    result,
	}).Debug("Checking addition for match.")
	return result
}

func (repo gitRepo) outgoingNonDeletedFiles(oldCommit, newCommit string) []*patch.File {
	result := make([]*patch.File, 0)
	for _, file := range repo.changedFiles(oldCommit, newCommit) {
		if file.Verb != patch.Delete {
			result = append(result, file)
		}
	}
	return result
}

func (repo gitRepo) changedFiles(oldCommit, newCommit string) []*patch.File {
	diff := repo.fetchRawOutgoingDiff(oldCommit, newCommit)
	logEntry := log.WithFields(log.Fields{
		"oldCommit": oldCommit,
		"newCommit": newCommit,
	})
	patchSet, err := patch.Parse(diff)
	if err != nil {
		logEntry.WithError(err).Fatal("Unable to parse the the diff")
	} else {
		logEntry.Debug("Diff parsed successfully")
	}
	return patchSet.File
}

func (repo gitRepo) fetchRawOutgoingDiff(oldCommit string, newCommit string) []byte {
	gitRange := oldCommit + ".." + newCommit
	return repo.executeRepoCommand("git", "diff", gitRange, "--binary")
}

func (repo gitRepo) executeRepoCommand(commandName string, args ...string) []byte {
	log.WithFields(log.Fields{
		"command": commandName,
		"args":    args,
	}).Debug("Building repo command")
	result := exec.Command(commandName, args...)
	result.Dir = repo.root
	o, err := result.Output()
	logEntry := log.WithFields(log.Fields{
		"output": string(o),
	})
	if err == nil {
		logEntry.Debug("Command excuted successfully")
	} else {
		logEntry.WithError(err).Fatal("Command execution failed")
	}
	return o
}
