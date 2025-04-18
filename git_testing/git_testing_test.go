package git_testing

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var logger *logrus.Entry

func init() {
	Logger = logrus.WithField("Environment", "Debug")
	Logger.Debug("GitTesting test started")
	logrus.SetOutput(os.Stderr)
	logger = Logger
}

func TestInitializingANewRepoSetsUpFolderAndGitStructures(t *testing.T) {
	repo := Init()
	defer repo.Clean()
	assert.True(t, exists(repo.gitRoot), "GitTesting initialization should create the directory structure required")
	assert.True(t, isGitRepo(repo.gitRoot), "Repo root does not contain the .git folder")
}

func TestSettingUpBaselineFilesSetsUpACommitInRepo(t *testing.T) {
	repo := Init()
	defer repo.Clean()
	repo.SetupBaselineFiles("a.txt", filepath.Join("alice", "bob", "b.txt"))
	verifyPresenceOfGitRepoWithCommits(t, 1, repo.gitRoot)
}

func TestEditingFilesInARepoWorks(t *testing.T) {
	repo := Init()
	defer repo.Clean()
	repo.SetupBaselineFiles("a.txt", filepath.Join("alice", "bob", "b.txt"))
	repo.AppendFileContent("a.txt", "\nmonkey see.\n", "monkey do.")
	content := repo.FileContents("a.txt")
	assert.True(t, strings.HasSuffix(string(content), "monkey see.\nmonkey do."))
	repo.AddAndcommit("a.txt", "modified content")
	verifyPresenceOfGitRepoWithCommits(t, 2, repo.gitRoot)
}

func TestRemovingFilesInARepoWorks(t *testing.T) {
	repo := Init()
	defer repo.Clean()
	repo.SetupBaselineFiles("a.txt", filepath.Join("alice", "bob", "b.txt"))
	repo.RemoveFile("a.txt")
	assert.False(t, exists(filepath.Join("data", "testLocation1", "a.txt")), "Unexpected. Deleted file a.txt still exists inside the repo")
	repo.AddAndcommit("a.txt", "removed it")
	verifyPresenceOfGitRepoWithCommits(t, 2, repo.gitRoot)
}

func TestCloningARepoToAnotherWorks(t *testing.T) {
	repo := Init()
	defer repo.Clean()
	repo.SetupBaselineFiles("a.txt", filepath.Join("alice", "bob", "b.txt"))
	fs := afero.NewMemMapFs()
	anotherRepoLocation, _ := afero.TempDir(fs, afero.GetTempDir(fs, "talisman-clone-test"), "")
	cloned := repo.GitClone(anotherRepoLocation)
	defer cloned.Clean()
	verifyPresenceOfGitRepoWithCommits(t, 1, repo.gitRoot)
	logger.Debug("Finished with first verification")
	logger.Debugf("Another location is %s\n", anotherRepoLocation)
	verifyPresenceOfGitRepoWithCommits(t, 1, anotherRepoLocation)
}

func TestEarliestCommits(t *testing.T) {
	repo := Init()
	defer repo.Clean()
	repo.SetupBaselineFiles("a.txt")
	initialCommit := repo.EarliestCommit()
	repo.AppendFileContent("a.txt", "\nmonkey see.\n", "monkey do.")
	repo.AddAndcommit("a.txt", "modified content")
	assert.Equal(t, initialCommit, repo.EarliestCommit(), "First commit is not expected to change on repo modifications")
}

func TestLatestCommits(t *testing.T) {
	repo := Init()
	defer repo.Clean()
	repo.SetupBaselineFiles("a.txt")
	repo.AppendFileContent("a.txt", "\nmonkey see.\n", "monkey do.")
	repo.AddAndcommit("a.txt", "modified content")
	repo.AppendFileContent("a.txt", "\nline n-1.\n", "line n.")
	repo.AddAndcommit("a.txt", "more modified content")
	assert.NotEqual(t, repo.EarliestCommit(), repo.LatestCommit()) //bad test.
}

func verifyPresenceOfGitRepoWithCommits(t *testing.T, expectedCommitCount int, repoLocation string) {
	wd, _ := os.Getwd()
	os.Chdir(repoLocation)
	defer func() { os.Chdir(wd) }()

	cmd := exec.Command("git", "log", "--pretty=short")
	o, err := cmd.CombinedOutput()
	dieOnError(err)
	matches := regExp("(?m)^commit\\s[a-z0-9]+\\s+.*$").FindAllString(string(o), -1)
	assert.Len(t, matches, expectedCommitCount, "Repo root does not contain exactly %d commits.", expectedCommitCount)
}

func regExp(pattern string) *regexp.Regexp {
	return regexp.MustCompile(pattern)
}

func isGitRepo(loc string) bool {
	return exists(path.Join(loc, ".git"))
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if (err != nil) && (os.IsNotExist(err)) {
		return false
	} else if err != nil {
		dieOnError(err)
		return true
	} else {
		return true
	}
}

func dieOnError(err error) {
	if err != nil {
		panic(err)
	}
}
