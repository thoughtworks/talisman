package git_repo

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	git "github.com/thoughtworks/talisman/git_testing"
)

func TestEmptyRepoReturnsNoFileChanges(t *testing.T) {
	cleanTestData()
	_, repo := setupOriginAndClones("data/testLocation1", "data/cloneLocation")
	assert.Len(t, repo.AllChanges(), 0, "Empty git repo should not have any changes")
}

func TestAdditionsReturnsEditsAndAdds(t *testing.T) {
	cleanTestData()
	_, repo := setupOriginAndClones("data/testLocation1", "data/cloneLocation")
	git.AppendFileContent(repo.root, "a.txt", "New content.\n", "Spanning multiple lines, even.")
	git.CreateFileWithContents(repo.root, "new.txt", "created contents")
	git.AddAndcommit(repo.root, "*", "added to lorem-ipsum content with my own stuff!")

	additions := repo.Additions("HEAD~1", "HEAD")
	assert.Len(t, additions, 2)
	assert.True(t, strings.HasSuffix(string(additions[0].Data), "New content.\nSpanning multiple lines, even."))
}

func TestNewlyAddedFilesAreCountedAsChanges(t *testing.T) {
	cleanTestData()
	_, repo := setupOriginAndClones("data/testLocation1", "data/cloneLocation")
	git.CreateFileWithContents(repo.root, "h", "Hello")
	git.CreateFileWithContents(repo.root, "foo/bar/w", ", World!")
	git.AddAndcommit(repo.root, "*", "added hello world")
	assert.Len(t, repo.AllChanges(), 2)
}

func TestOutgoingContentOfNewlyAddedFilesIsAvailableInChanges(t *testing.T) {
	cleanTestData()
	_, repo := setupOriginAndClones("data/testLocation1", "data/cloneLocation")
	git.CreateFileWithContents(repo.root, "foo/bar/w", "new contents")
	git.AddAndcommit(repo.root, "*", "added new files")

	assert.Len(t, repo.AllChanges(), 1)
	assert.True(t, strings.HasSuffix(string(repo.AllAdditions()[0].Data), "new contents"))
}

func TestOutgoingContentOfModifiedFilesIsAvailableInChanges(t *testing.T) {
	cleanTestData()
	_, repo := setupOriginAndClones("data/testLocation1", "data/cloneLocation")
	git.AppendFileContent(repo.root, "a.txt", "New content.\n", "Spanning multiple lines, even.")
	git.AddAndcommit(repo.root, "a.txt", "added to lorem-ipsum content with my own stuff!")
	assert.Len(t, repo.AllChanges(), 1)
	assert.True(t, strings.HasSuffix(string(repo.AllAdditions()[0].Data), "New content.\nSpanning multiple lines, even."))
}

func TestMultipleOutgoingChangesToTheSameFileAreAvailableInAdditions(t *testing.T) {
	cleanTestData()
	_, repo := setupOriginAndClones("data/testLocation1", "data/cloneLocation")
	git.AppendFileContent(repo.root, "a.txt", "New content.\n")
	git.AddAndcommit(repo.root, "a.txt", "added some new content")

	git.AppendFileContent(repo.root, "a.txt", "More new content.\n")
	git.AddAndcommit(repo.root, "a.txt", "added some more new content")

	assert.Len(t, repo.AllChanges(), 1)
	assert.True(t, strings.HasSuffix(string(repo.AllAdditions()[0].Data), "New content.\nMore new content.\n"))
}

func TestContentOfDeletedFilesIsNotAvailableInChanges(t *testing.T) {
	cleanTestData()
	_, repo := setupOriginAndClones("data/testLocation1", "data/cloneLocation")
	git.RemoveFile(repo.root, "a.txt")
	git.AddAndcommit(repo.root, "a.txt", "Deleted this file. After all, it only had lorem-ipsum content.")
	assert.Len(t, repo.AllChanges(), 1)
	assert.Equal(t, 0, len(repo.AllAdditions()), "There should be no additions because there only an outgoing deletion")
}

func setupOriginAndClones(originLocation, cloneLocation string) (gitRepo, gitRepo) {
	origin := RepoLocatedAt(originLocation)
	git.Init(origin.root)
	git.SetupBaselineFiles(origin.root, "a.txt", "alice/bob/b.txt")
	git.GitClone(origin.root, cloneLocation)
	return origin, RepoLocatedAt(cloneLocation)
}
