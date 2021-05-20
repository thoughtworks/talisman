package gitrepo

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"talisman/utility"
)

//FilePath represents the absolute path of an added file
type FilePath string

//FileName represents the base name of an added file
type FileName string

//Addition represents the end state of a file
type Addition struct {
	Path    FilePath
	Name    FileName
	Commits []string
	Data    []byte
}

//GitRepo represents a Git repository located at the absolute path represented by root
type GitRepo struct {
	root string
}

//RepoLocatedAt returns a new GitRepo with it's root located at the location specified by the argument.
//If the argument is not an absolute path, it will be turned into one.
func RepoLocatedAt(path string) GitRepo {
	absoluteRoot, _ := filepath.Abs(path)
	return GitRepo{absoluteRoot}
}

// GetDiffForStagedFiles gets all the staged files and collects the diff section in each file
func (repo GitRepo) GetDiffForStagedFiles() []Addition {
	stagedContent := repo.executeRepoCommand("git", "diff", "--staged")
	content := strings.TrimSpace(string(stagedContent))
	lines := strings.Split(content, "\n")
	result := make([]Addition, 0)

	if len(lines) < 1 {
		return result
	}

	// Standard git diff header pattern
	// ref: https://git-scm.com/docs/diff-format#_generating_patches_with_p

	lineNumberOfFirstHeader := 0
	var additionFilename string
	for ; lineNumberOfFirstHeader < len(lines); lineNumberOfFirstHeader++ {
		match, stagedFilename := MatchGitDiffLine(lines[lineNumberOfFirstHeader])
		if match {
			additionFilename = stagedFilename
			break
		}
	}

	additionContentBuffer := &strings.Builder{}
	for i := lineNumberOfFirstHeader + 1; i < len(lines); i++ {
		match, stagedFilename := MatchGitDiffLine(lines[i])
		if match {
			// It is a new diff header
			// which means we have reached the next file's header

			// capture content written to buffer so far as addition content
			stagedChanges := repo.extractAdditions(additionContentBuffer.String())
			if stagedChanges != nil {
				addition := NewAddition(additionFilename, stagedChanges)
				result = append(
					result, addition,
				)
			}

			// get next file name and reset buffer for next iteration
			additionFilename = stagedFilename
			additionContentBuffer.Reset()
		} else {
			additionContentBuffer.WriteString(lines[i])
			additionContentBuffer.WriteRune('\n')
		}
	}

	// Save last file's diff content
	stagedChanges := repo.extractAdditions(additionContentBuffer.String())
	if stagedChanges != nil {
		addition := NewAddition(additionFilename, stagedChanges)
		result = append(result, addition)
	}

	log.WithFields(log.Fields{
		"additions": result,
	}).Debug("Generating staged additions.")

	return result
}

func MatchGitDiffLine(gitDiffString string) (bool, string) {
	if strings.Contains(gitDiffString, "diff --git") {
		fileNameLength := (len(gitDiffString) - len("diff --git a/ b/"))/2
		regexPattern := fmt.Sprintf("^diff --git a/(.{%v}) b/(.{%v})$", fileNameLength, fileNameLength);
		headerRegex := regexp.MustCompile(regexPattern)

		if headerRegex.MatchString(gitDiffString) {
			matches := headerRegex.FindStringSubmatch(gitDiffString)
			if matches[1] == matches[2] {
				return true, matches[1]
			}
		}
	}
	return false, ""
}

//StagedAdditions returns the files staged for commit in a GitRepo
func (repo GitRepo) StagedAdditions() []Addition {
	files := repo.stagedFiles()
	result := make([]Addition, len(files))
	for i, file := range files {
		data := repo.stagedVersionOfFile(file)
		result[i] = NewAddition(file, data)
	}

	log.WithFields(log.Fields{
		"additions": result,
	}).Info("Generating staged additions.")
	return result
}

//AllAdditions returns all the outgoing additions and modifications in a GitRepo. This does not include files that were deleted.
func (repo GitRepo) AllAdditions() []Addition {
	return repo.AdditionsWithinRange("origin/master", "master")
}

//AdditionsWithinRange returns the outgoing additions and modifications in a GitRepo that are in the given commit range. This does not include files that were deleted.
func (repo GitRepo) AdditionsWithinRange(oldCommit string, newCommit string) []Addition {
	files := repo.outgoingNonDeletedFiles(oldCommit, newCommit)
	result := make([]Addition, len(files))
	for i, file := range files {
		data, _ := repo.ReadRepoFile(file)
		result[i] = NewAddition(file, data)
	}
	log.WithFields(log.Fields{
		"oldCommit": oldCommit,
		"newCommit": newCommit,
		"additions": result,
	}).Info("Generating all additions in range.")
	return result
}

//NewAddition returns a new Addition for a file with supplied name and contents
func NewAddition(filePath string, content []byte) Addition {
	return Addition{
		Path: FilePath(filePath),
		Name: FileName(path.Base(filePath)),
		Data: content,
	}
}

//NewScannerAddition returns an new Addition for a file with supplied contents and all of the commits the file is in
func NewScannerAddition(filePath string, commits []string, content []byte) Addition {
	return Addition{
		Path:    FilePath(filePath),
		Name:    FileName(path.Base(filePath)),
		Commits: commits,
		Data:    content,
	}
}

//ReadRepoFile returns the contents of the supplied relative filename by locating it in the git repo
func (repo GitRepo) ReadRepoFile(fileName string) ([]byte, error) {
	path := filepath.Join(repo.root, fileName)
	log.Debugf("reading file %s", path)
	return utility.SafeReadFile(path)
}

//ReadRepoFileOrNothing returns the contents of the supplied relative filename by locating it in the git repo.
//If the given file cannot be located in the repo, then an empty array of bytes is returned for the content.
func (repo GitRepo) ReadRepoFileOrNothing(fileName string) ([]byte, error) {
	filepath := path.Join(repo.root, fileName)
	if _, err := os.Stat(filepath); err == nil {
		return repo.ReadRepoFile(fileName)
	}
	return make([]byte, 0), nil
}

//CheckIfFileExists checks if the file exists on the file system. Does not look into the file contents
//Returns TRUE if file exists
//Returns FALSE if the file is not found
func (repo GitRepo) CheckIfFileExists(fileName string) bool {
	filepath := path.Join(repo.root, fileName)
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	}
	return true
}

//Matches states whether the addition matches the given pattern.
//If the pattern ends in a path separator, then all files inside a directory with that name are matched. However, files with that name itself will not be matched.
//If a pattern contains the path separator in any other location, the match works according to the pattern logic of the default golang glob mechanism
//If there is no path separator anywhere in the pattern, the pattern is matched against the base name of the file. Thus, the pattern will match files with that name anywhere in the repository.
func (a Addition) Matches(pattern string) bool {
	var result bool
	if pattern[len(pattern)-1] == '/' {
		result = strings.HasPrefix(string(a.Path), pattern)
	} else if strings.ContainsRune(pattern, '/') {
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

//TrackedFilesAsAdditions returns all of the tracked files in a GitRepo as Additions
func (repo GitRepo) TrackedFilesAsAdditions() []Addition {
	trackedFilePaths := repo.trackedFilePaths()
	var additions []Addition
	for _, path := range trackedFilePaths {
		additions = append(additions, NewAddition(path, make([]byte, 0)))
	}
	return additions
}

func (repo GitRepo) trackedFilePaths() []string {
	branchName := repo.currentBranch()
	if len(branchName) == 0 {
		return make([]string, 0)
	}
	byteArray := repo.executeRepoCommand("git", "ls-tree", branchName, "--name-only", "-r")
	trackedFilePaths := strings.Split(string(byteArray), "\n")
	return trackedFilePaths
}

func (repo GitRepo) stagedFiles() []string {
	stagedFiles := strings.Split(repo.fetchStagedChanges(), "\n")
	var result []string
	for _, c := range stagedFiles {
		if len(c) != 0 {
			changeTypeAndFile := strings.Split(c, "\t")
			if len(changeTypeAndFile) > 0 {
				file := changeTypeAndFile[1]
				result = append(result, file)
			}
		}
	}
	return result
}

func (repo GitRepo) currentBranch() string {
	if !repo.hasBranch() {
		return ""
	}
	byteArray := repo.executeRepoCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchName := strings.TrimSpace(string(byteArray))
	return branchName
}

func (repo GitRepo) hasBranch() bool {
	byteArray := repo.executeRepoCommand("git", "branch")
	if len(string(byteArray)) != 0 {
		return true
	}
	return false
}

func (repo GitRepo) stagedVersionOfFile(file string) []byte {
	return repo.executeRepoCommand("git", "show", ":"+file)
}

func (repo GitRepo) outgoingNonDeletedFiles(oldCommit, newCommit string) []string {
	allChanges := strings.Split(repo.fetchRawOutgoingDiff(oldCommit, newCommit), "\n")
	var result []string
	for _, c := range allChanges {
		if len(c) != 0 {
			result = append(result, c)
		}
	}
	return result
}

func (repo *GitRepo) fetchStagedChanges() string {
	return string(repo.executeRepoCommand("git", "diff", "--cached", "--name-status", "--diff-filter=ACM"))
}

// extractAdditions will accept git diff --staged {file} output and filters the command output
// to get only the modified sections of the file
func (repo *GitRepo) extractAdditions(diffContent string) []byte {
	var result []byte
	changes := strings.Split(diffContent, "\n")
	for _, c := range changes {
		if !strings.HasPrefix(c, "+++") && !strings.HasPrefix(c, "---") && strings.HasPrefix(c, "+") {

			result = append(result, strings.TrimPrefix(c, "+")...)
			result = append(result, "\n"...)
		}
	}
	return result
}

func (repo GitRepo) fetchRawOutgoingDiff(oldCommit string, newCommit string) string {
	gitRange := oldCommit + ".." + newCommit
	return string(repo.executeRepoCommand("git", "diff", gitRange, "--name-only", "--diff-filter=ACM"))
}

func (repo GitRepo) executeRepoCommand(commandName string, args ...string) []byte {
	log.WithFields(log.Fields{
		"command": commandName,
		"args":    args,
	}).Debug("Building repo command")
	result := exec.Command(commandName, args...)
	result.Dir = repo.root
	co, err := result.CombinedOutput()
	logEntry := log.WithFields(log.Fields{
		"dir":     repo.root,
		"command": fmt.Sprintf("%s %s", commandName, strings.Join(args, " ")),
		"output":  string(co),
		"error":   err,
	})
	if err == nil {
		logEntry.Debug("Git command excuted successfully")
	} else {
		logEntry.Fatal("Git command execution failed")
	}
	return co
}
