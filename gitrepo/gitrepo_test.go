package gitrepo

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"talisman/git_testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var testLocation = filepath.Join("data", "testLocation")
var cloneLocation = filepath.Join("data", "cloneLocation")

func init() {
	git_testing.Logger = logrus.WithField("Environment", "Debug")
	git_testing.Logger.Debug("Accetpance test started")
	logger = git_testing.Logger
}

func TestEmptyRepoReturnsNoFileChanges(t *testing.T) {
	cleanTestData()
	_, repo := setupOriginAndClones(testLocation, cloneLocation)
	assert.Len(t, repo.AllAdditions(), 0, "Empty git repo should not have any changes")
}

func TestGetDiffForStagedFiles(t *testing.T) {
	cleanTestData()
	git, repo := setupOriginAndClones(testLocation, cloneLocation)
	git.AppendFileContent("a.txt", "New content.\n", "Spanning multiple lines, even.")
	git.CreateFileWithContents("new.txt", "created contents")
	git.Add("a.txt")
	git.Add("new.txt")
	additions := repo.GetDiffForStagedFiles()

	if assert.Len(t, additions, 2) {
		modifiedAddition := additions[0]
		createdAddition := additions[1]

		aTxtFileContents, err := os.ReadFile(path.Join(cloneLocation, "a.txt"))
		assert.NoError(t, err)
		newTxtFileContents, err := os.ReadFile(path.Join(cloneLocation, "new.txt"))
		assert.NoError(t, err)

		expectedModifiedAddition := Addition{
			Path: FilePath("a.txt"),
			Name: FileName("a.txt"),
			Data: []byte(fmt.Sprintf("%s\n", string(aTxtFileContents))),
		}

		expectedCreatedAddition := Addition{
			Path: FilePath("new.txt"),
			Name: FileName("new.txt"),
			Data: []byte(fmt.Sprintf("%s\n", string(newTxtFileContents))),
		}

		// For human-readable comparison
		assert.Equal(t, string(expectedModifiedAddition.Data), string(modifiedAddition.Data))
		assert.Equal(t, string(expectedCreatedAddition.Data), string(createdAddition.Data))

		assert.Equal(t, expectedModifiedAddition, modifiedAddition)
		assert.Equal(t, expectedCreatedAddition, createdAddition)
	}

}

func TestGetDiffForStagedFilesWithSpacesInPath(t *testing.T) {
	cleanTestData()
	git, repo := setupOriginAndClones(testLocation, cloneLocation)
	git.AppendFileContent("folder b/c.txt", "New content.\n", "Spanning multiple lines, even.")
	git.Add("folder b/c.txt")
	additions := repo.GetDiffForStagedFiles()

	if assert.Len(t, additions, 1) {
		modifiedAddition := additions[0]

		aTxtFileContents, err := os.ReadFile(path.Join(cloneLocation, "folder b/c.txt"))
		assert.NoError(t, err)

		expectedModifiedAddition := Addition{
			Path: FilePath("folder b/c.txt"),
			Name: FileName("c.txt"),
			Data: []byte(fmt.Sprintf("%s\n", string(aTxtFileContents))),
		}

		// For human-readable comparison
		assert.Equal(t, string(expectedModifiedAddition.Data), string(modifiedAddition.Data))

		assert.Equal(t, expectedModifiedAddition, modifiedAddition)
	}
}

func TestAdditionsReturnsEditsAndAdds(t *testing.T) {
	cleanTestData()
	git, repo := setupOriginAndClones(testLocation, cloneLocation)
	git.AppendFileContent("a.txt", "New content.\n", "Spanning multiple lines, even.")
	git.CreateFileWithContents("new.txt", "created contents")
	git.AddAndcommit("*", "added to lorem-ipsum content with my own stuff!")

	additions := repo.AdditionsWithinRange("HEAD~1", "HEAD")
	assert.Len(t, additions, 2)
	assert.True(t, strings.HasSuffix(string(additions[0].Data), "New content.\nSpanning multiple lines, even."))
}

func TestNewlyAddedFilesAreCountedAsChanges(t *testing.T) {
	cleanTestData()
	git, repo := setupOriginAndClones(testLocation, cloneLocation)
	git.CreateFileWithContents("h", "Hello")
	git.CreateFileWithContents("foo/bar/w", ", World!")
	git.AddAndcommit("*", "added hello world")
	assert.Len(t, repo.AllAdditions(), 2)
}

func TestOutgoingContentOfNewlyAddedFilesIsAvailableInChanges(t *testing.T) {
	cleanTestData()
	git, repo := setupOriginAndClones(testLocation, cloneLocation)
	git.CreateFileWithContents("foo/bar/w", "new contents")
	git.AddAndcommit("*", "added new files")

	assert.Len(t, repo.AllAdditions(), 1)
	assert.True(t, strings.HasSuffix(string(repo.AllAdditions()[0].Data), "new contents"))
}

func TestOutgoingContentOfModifiedFilesIsAvailableInChanges(t *testing.T) {
	cleanTestData()
	git, repo := setupOriginAndClones(testLocation, cloneLocation)
	git.AppendFileContent("a.txt", "New content.\n", "Spanning multiple lines, even.")
	git.AddAndcommit("a.txt", "added to lorem-ipsum content with my own stuff!")
	assert.Len(t, repo.AllAdditions(), 1)
	assert.True(t, strings.HasSuffix(string(repo.AllAdditions()[0].Data), "New content.\nSpanning multiple lines, even."))
}

func TestMultipleOutgoingChangesToTheSameFileAreAvailableInAdditions(t *testing.T) {
	cleanTestData()
	git, repo := setupOriginAndClones(testLocation, cloneLocation)
	git.AppendFileContent("a.txt", "New content.\n")
	git.AddAndcommit("a.txt", "added some new content")

	git.AppendFileContent("a.txt", "More new content.\n")
	git.AddAndcommit("a.txt", "added some more new content")

	assert.Len(t, repo.AllAdditions(), 1)
	assert.True(t, strings.HasSuffix(string(repo.AllAdditions()[0].Data), "New content.\nMore new content.\n"))
}

func TestContentOfDeletedFilesIsNotAvailableInChanges(t *testing.T) {
	cleanTestData()
	git, repo := setupOriginAndClones(testLocation, cloneLocation)
	git.RemoveFile("a.txt")
	git.AddAndcommit("a.txt", "Deleted this file. After all, it only had lorem-ipsum content.")
	assert.Equal(t, 0, len(repo.AllAdditions()), "There should be no additions because there only an outgoing deletion")
}

func TestDiffContainingBinaryFileChangesDoesNotBlowUp(t *testing.T) {
	cleanTestData()
	git, repo := setupOriginAndClones(testLocation, cloneLocation)
	exec.Command("cp", "./pixel.jpg", repo.root).Run()
	git.AddAndcommit("pixel.jpg", "Testing binary diff.")
	assert.Len(t, repo.AllAdditions(), 1)
	assert.Equal(t, "pixel.jpg", string(repo.AllAdditions()[0].Name))
}

func TestStagedAdditionsIncludeStagedFiles(t *testing.T) {
	cleanTestData()
	git, repo := setupOriginAndClones(testLocation, cloneLocation)
	git.OverwriteFileContent("a.txt", "New content.\n")
	git.Add("a.txt")

	git.AppendFileContent("a.txt", "More new content\n")
	git.AppendFileContent("alice/bob/b.txt", "New content to b\n")

	stagedAdditions := repo.StagedAdditions()
	assert.Len(t, stagedAdditions, 1)
	assert.Equal(t, "a.txt", string(stagedAdditions[0].Name))
	assert.Equal(t, "New content.\n", string(stagedAdditions[0].Data))
}

func TestStagedAdditionsIncludeStagedNewFiles(t *testing.T) {
	cleanTestData()
	git, repo := setupOriginAndClones(testLocation, cloneLocation)
	git.CreateFileWithContents("new.txt", "New content.\n")
	git.Add("new.txt")

	stagedAdditions := repo.StagedAdditions()
	assert.Len(t, stagedAdditions, 1)
	assert.Equal(t, "new.txt", string(stagedAdditions[0].Name))
	assert.Equal(t, "New content.\n", string(stagedAdditions[0].Data))
}

func TestStagedAdditionsShouldNotIncludeDeletedFiles(t *testing.T) {
	cleanTestData()
	git, repo := setupOriginAndClones(testLocation, cloneLocation)
	git.RemoveFile("a.txt")
	git.Add(".")

	stagedAdditions := repo.StagedAdditions()
	assert.Len(t, stagedAdditions, 0)
}

func TestMatchShouldMatchExactFileIfNoPatternIsProvided(t *testing.T) {
	file1 := Addition{Path: "bigfile", Name: "bigfile"}
	file2 := Addition{Path: "subfolder/bigfile", Name: "bigfile"}
	file3 := Addition{Path: "somefile", Name: "somefile"}
	pattern := "bigfile"

	assert.True(t, file1.Matches(pattern))
	assert.False(t, file2.Matches(pattern))
	assert.False(t, file3.Matches(pattern))
}

func TestMatchingWithPatterns(t *testing.T) {
	files := []Addition{
		NewAddition("GitRepoPath1/File1.txt", nil),
		NewAddition("GitRepoPath1/File2.txt", nil),
		NewAddition("GitRepoPath1/somefile", nil),
		NewAddition("somefile.jpg", nil),
		NewAddition("somefile.txt", nil),
		NewAddition("File1.txt", nil),
		NewAddition("File3.txt", nil),
	}

	expectedToMatch := map[string][]bool{
		"GitRepoPath1/*.txt": {true, true, false, false, false, false, false},
		"*.txt":              {true, true, false, false, true, true, true},
		"File?.txt":          {true, true, false, false, false, true, true},
		"File[1].txt":        {true, false, false, false, false, true, false},
		"File[1-2].txt":      {true, true, false, false, false, true, false},
		"File\\1.txt":        {true, false, false, false, false, true, false},
	}

	for pattern, expectedResults := range expectedToMatch {
		t.Run(fmt.Sprintf("Testing matches for pattern %s", pattern), func(t *testing.T) {
			for i := range files {
				assert.Equal(t, expectedResults[i], files[i].Matches(pattern))
			}
		})
	}
}

func TestMatchingAdditionBasename(t *testing.T) {
	addition := NewAddition("subdirectory/nested-file", nil)
	assert.False(t, addition.Matches("nested-file"))
	assert.True(t, addition.NameMatches("nested-file"))
}

func setupOriginAndClones(originLocation, cloneLocation string) (*git_testing.GitTesting, GitRepo) {
	origin := RepoLocatedAt(originLocation)
	git := git_testing.InitAt(origin.root)
	git.SetupBaselineFiles("a.txt", filepath.Join("alice", "bob", "b.txt"))
	git.SetupBaselineFiles("c.txt", filepath.Join("folder b", "c.txt"))
	cwd, _ := os.Getwd()
	gitClone := git.GitClone(filepath.Join(cwd, cloneLocation))
	return gitClone, RepoLocatedAt(cloneLocation)
}
