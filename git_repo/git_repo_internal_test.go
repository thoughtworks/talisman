package git_repo

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	git "github.com/badrij/talisman/git_testing"
	"github.com/stretchr/testify/assert"
)

func TestNewRepoGetsCreatedWithAbsolutePath(t *testing.T) {
	cleanTestData()
	repo := RepoLocatedAt("data/testLocation1")
	assert.True(t, path.IsAbs(repo.root))
}

func TestInitializingANewRepoSetsUpFolderAndGitStructures(t *testing.T) {
	cleanTestData()
	repo := RepoLocatedAt("data/dir/sub_dir/testLocation2")
	git.Init(repo.root)
	assert.True(t, exists(repo.root), "Git Repo initialization should create the directory structure required")
	assert.True(t, isGitRepo(repo.root), "Repo root does not contain the .git folder")
}

func TestSettingUpBaselineFilesSetsUpACommitInRepo(t *testing.T) {
	cleanTestData()
	repo := RepoLocatedAt("data/testLocation1")
	git.Init(repo.root)
	git.SetupBaselineFiles(repo.root, "a.txt", "alice/bob/b.txt")
	verifyPresenseOfGitRepoWithCommits("data/testLocation1", 1, t)
}

func TestEditingFilesInARepoWorks(t *testing.T) {
	cleanTestData()
	repo := RepoLocatedAt("data/testLocation1")
	git.Init(repo.root)
	git.SetupBaselineFiles(repo.root, "a.txt", "alice/bob/b.txt")
	git.AppendFileContent(repo.root, "a.txt", "\nmonkey see.\n", "monkey do.")
	content := git.FileContents(repo.root, "a.txt")
	assert.True(t, strings.HasSuffix(string(content), "monkey see.\nmonkey do."))
	git.AddAndcommit(repo.root, "a.txt", "modified content")
	verifyPresenseOfGitRepoWithCommits("data/testLocation1", 2, t)
}

func TestRemovingFilesInARepoWorks(t *testing.T) {
	cleanTestData()
	repo := RepoLocatedAt("data/testLocation1")
	git.Init(repo.root)
	git.SetupBaselineFiles(repo.root, "a.txt", "alice/bob/b.txt")
	git.RemoveFile(repo.root, "a.txt")
	assert.False(t, exists("data/testLocation1/a.txt"), "Unexpected. Deleted file a.txt still exists inside the repo")
	git.AddAndcommit(repo.root, "a.txt", "removed it")
	verifyPresenseOfGitRepoWithCommits("data/testLocation1", 2, t)
}

func TestCloningARepoToAnotherWorks(t *testing.T) {
	cleanTestData()
	repo := RepoLocatedAt("data/testLocation1")
	git.Init(repo.root)
	git.SetupBaselineFiles(repo.root, "a.txt", "alice/bob/b.txt")
	git.GitClone(repo.root, "data/somewhereElse/testLocationClone")
	verifyPresenseOfGitRepoWithCommits("data/testLocation1", 1, t)
	verifyPresenseOfGitRepoWithCommits("data/somewhereElse/testLocationClone", 1, t)
}

func TestEarliestCommits(t *testing.T) {
	cleanTestData()
	repo := RepoLocatedAt("data/testLocation1")
	git.Init(repo.root)
	git.SetupBaselineFiles(repo.root, "a.txt")
	initialCommit := git.EarliestCommit(repo.root)
	git.AppendFileContent(repo.root, "a.txt", "\nmonkey see.\n", "monkey do.")
	git.AddAndcommit(repo.root, "a.txt", "modified content")
	assert.Equal(t, initialCommit, git.EarliestCommit(repo.root), "First commit is not expected to change on repo modifications")
}

func TestLatestCommits(t *testing.T) {
	cleanTestData()
	repo := RepoLocatedAt("data/testLocation1")
	git.Init(repo.root)
	git.SetupBaselineFiles(repo.root, "a.txt")
	git.AppendFileContent(repo.root, "a.txt", "\nmonkey see.\n", "monkey do.")
	git.AddAndcommit(repo.root, "a.txt", "modified content")
	git.AppendFileContent(repo.root, "a.txt", "\nline n-1.\n", "line n.")
	git.AddAndcommit(repo.root, "a.txt", "more modified content")
	assert.NotEqual(t, git.EarliestCommit(repo.root), git.LatestCommit(repo.root)) //bad test.
}

func verifyPresenseOfGitRepoWithCommits(location string, expectedCommitCount int, t *testing.T) {
	cmd := exec.Command("git", "log")
	cmd.Dir = location
	o, err := cmd.CombinedOutput()
	check(err)
	matches := regExp("(?m)^commit\\s[a-z0-9]+$").FindAllString(string(o), -1)
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
		check(err)
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
	check(err)
	defer d.Close()
	names, err := d.Readdirnames(-1)
	check(err)
	for _, name := range names {
		check(os.RemoveAll(filepath.Join(dataDir, name)))
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
