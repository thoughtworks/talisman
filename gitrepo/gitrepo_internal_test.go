package gitrepo

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"talisman/git_testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var testLocation1 = filepath.Join("data", "testLocation1")

var logger *logrus.Entry

func init() {
	git_testing.Logger = logrus.WithField("Environment", "Debug")
	git_testing.Logger.Debug("Acceptance test started")
	logrus.SetOutput(os.Stderr)
	logger = git_testing.Logger
}

func TestNewRepoGetsCreatedWithAbsolutePath(t *testing.T) {
	cleanTestData()
	repo := RepoLocatedAt(testLocation1)
	assert.True(t, path.IsAbs(repo.root))
}

func TestInitializingANewRepoSetsUpFolderAndGitStructures(t *testing.T) {
	cleanTestData()
	repo := RepoLocatedAt(filepath.Join("data", "dir", "sub_dir", "testLocation2"))
	git_testing.Init(repo.root)
	assert.True(t, exists(repo.root), "Git Repo initialization should create the directory structure required")
	assert.True(t, isGitRepo(repo.root), "Repo root does not contain the .git folder")
}

func TestSettingUpBaselineFilesSetsUpACommitInRepo(t *testing.T) {
	cleanTestData()
	repo := RepoLocatedAt(testLocation1)
	git := git_testing.Init(repo.root)
	git.SetupBaselineFiles("a.txt", filepath.Join("alice", "bob", "b.txt"))
	verifyPresenceOfGitRepoWithCommits(testLocation1, 1, t)
}

func TestEditingFilesInARepoWorks(t *testing.T) {
	cleanTestData()
	repo := RepoLocatedAt(testLocation1)
	git := git_testing.Init(repo.root)
	git.SetupBaselineFiles("a.txt", filepath.Join("alice", "bob", "b.txt"))
	git.AppendFileContent("a.txt", "\nmonkey see.\n", "monkey do.")
	content := git.FileContents("a.txt")
	assert.True(t, strings.HasSuffix(string(content), "monkey see.\nmonkey do."))
	git.AddAndcommit("a.txt", "modified content")
	verifyPresenceOfGitRepoWithCommits(testLocation1, 2, t)
}

func TestRemovingFilesInARepoWorks(t *testing.T) {
	cleanTestData()
	repo := RepoLocatedAt(testLocation1)
	git := git_testing.Init(repo.root)
	git.SetupBaselineFiles("a.txt", filepath.Join("alice", "bob", "b.txt"))
	git.RemoveFile("a.txt")
	assert.False(t, exists(filepath.Join("data", "testLocation1", "a.txt")), "Unexpected. Deleted file a.txt still exists inside the repo")
	git.AddAndcommit("a.txt", "removed it")
	verifyPresenceOfGitRepoWithCommits(testLocation1, 2, t)
}

func TestCloningARepoToAnotherWorks(t *testing.T) {
	cleanTestData()
	repo := RepoLocatedAt(testLocation1)
	git := git_testing.Init(repo.root)
	git.SetupBaselineFiles("a.txt", filepath.Join("alice", "bob", "b.txt"))
	cwd, _ := os.Getwd()
	anotherRepoLocation := filepath.Join(cwd, "data", "somewhereElse", "testLocationClone")
	git.GitClone(anotherRepoLocation)
	verifyPresenceOfGitRepoWithCommits(testLocation1, 1, t)
	logger.Debug("Finished with first verification")
	logger.Debugf("Another location is %s\n", anotherRepoLocation)
	verifyPresenceOfGitRepoWithCommits(anotherRepoLocation, 1, t)
}

func TestEarliestCommits(t *testing.T) {
	cleanTestData()
	repo := RepoLocatedAt(testLocation1)
	git := git_testing.Init(repo.root)
	git.SetupBaselineFiles("a.txt")
	initialCommit := git.EarliestCommit()
	git.AppendFileContent("a.txt", "\nmonkey see.\n", "monkey do.")
	git.AddAndcommit("a.txt", "modified content")
	assert.Equal(t, initialCommit, git.EarliestCommit(), "First commit is not expected to change on repo modifications")
}

func TestLatestCommits(t *testing.T) {
	cleanTestData()
	repo := RepoLocatedAt(testLocation1)
	git := git_testing.Init(repo.root)
	git.SetupBaselineFiles("a.txt")
	git.AppendFileContent("a.txt", "\nmonkey see.\n", "monkey do.")
	git.AddAndcommit("a.txt", "modified content")
	git.AppendFileContent("a.txt", "\nline n-1.\n", "line n.")
	git.AddAndcommit("a.txt", "more modified content")
	assert.NotEqual(t, git.EarliestCommit(), git.LatestCommit()) //bad test.
}

func verifyPresenceOfGitRepoWithCommits(location string, expectedCommitCount int, t *testing.T) {
	wd, _ := os.Getwd()
	os.Chdir(location)
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

func cleanTestData() {
	dataDir := testDataDir()
	if !exists(dataDir) {
		os.MkdirAll(dataDir, 0777)
	}
	d, err := os.Open(dataDir)
	dieOnError(err)
	defer d.Close()
	names, err := d.Readdirnames(-1)
	dieOnError(err)
	for _, name := range names {
		dieOnError(os.RemoveAll(filepath.Join(dataDir, name)))
	}
}

func testDataDir() string {
	wd, _ := os.Getwd()
	dataDir, _ := filepath.Abs(path.Join(wd, "data"))
	return dataDir
}

func TestMain(m *testing.M) {
	testExitStatus := m.Run()
	os.RemoveAll(testDataDir())
	os.Exit(testExitStatus)
}

func dieOnError(err error) {
	if err != nil {
		panic(err)
	}
}
