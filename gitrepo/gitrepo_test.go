package gitrepo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"talisman/git_testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	git_testing.Logger = logrus.WithField("Environment", "Debug")
	git_testing.Logger.Debug("GitRepo test started")
}

func (repo GitRepo) additionsInLastCommit() []Addition {
	return repo.AdditionsWithinRange("HEAD~1", "HEAD")
}

func TestNewRepoGetsCreatedWithAbsolutePath(t *testing.T) {
	var testLocation1 = filepath.Join("data", "testLocation1")
	repo := RepoLocatedAt(testLocation1)
	assert.True(t, filepath.IsAbs(repo.root))
}

func TestNoAdditionsBetweenSameRef(t *testing.T) {
	doInRepoWithCommit(func(git *git_testing.GitTesting) {
		assert.Len(t, RepoLocatedAt(git.Root()).AdditionsWithinRange("HEAD", "HEAD"), 0,
			"There should be no additions between a ref and itself.")
	})
}

func TestGetDiffForStagedFiles(t *testing.T) {
	doInRepoWithCommit(func(git *git_testing.GitTesting) {
		git.AppendFileContent("a.txt", "New content.\n", "Spanning multiple lines, even.")
		git.CreateFileWithContents("new.txt", "created contents")
		git.Add("a.txt")
		git.Add("new.txt")
		repo := RepoLocatedAt(git.Root())
		additions := repo.GetDiffForStagedFiles()

		if assert.Len(t, additions, 2) {
			modifiedAddition := additions[0]
			createdAddition := additions[1]

			aTxtFileContents, err := os.ReadFile(filepath.Join(repo.root, "a.txt"))
			assert.NoError(t, err)
			newTxtFileContents, err := os.ReadFile(filepath.Join(repo.root, "new.txt"))
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
	})
}

func TestGetDiffForStagedFilesWithSpacesInPath(t *testing.T) {
	doInRepoWithCommit(func(git *git_testing.GitTesting) {
		git.AppendFileContent("folder b/c.txt", "New content.\n", "Spanning multiple lines, even.")
		git.Add("folder b/c.txt")
		repo := RepoLocatedAt(git.Root())
		additions := repo.GetDiffForStagedFiles()

		if assert.Len(t, additions, 1) {
			modifiedAddition := additions[0]

			aTxtFileContents, err := os.ReadFile(filepath.Join(repo.root, "folder b/c.txt"))
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
	})
}

func TestAdditionsReturnsEditsAndAdds(t *testing.T) {
	doInRepoWithCommit(func(git *git_testing.GitTesting) {
		git.AppendFileContent("a.txt", "New content.\n", "Spanning multiple lines, even.")
		git.CreateFileWithContents("new.txt", "created contents")
		git.AddAndcommit("*", "added to lorem-ipsum content with my own stuff!")

		additions := RepoLocatedAt(git.Root()).additionsInLastCommit()
		assert.Len(t, additions, 2)
		assert.True(t, strings.HasSuffix(string(additions[0].Data), "New content.\nSpanning multiple lines, even."))
	})
}

func TestNewlyAddedFilesAreCountedAsChanges(t *testing.T) {
	doInRepoWithCommit(func(git *git_testing.GitTesting) {
		git.CreateFileWithContents("h", "Hello")
		git.CreateFileWithContents("foo/bar/w", ", World!")
		git.AddAndcommit("*", "added hello world")
		assert.Len(t, RepoLocatedAt(git.Root()).additionsInLastCommit(), 2)
	})
}

func TestOutgoingContentOfNewlyAddedFilesIsAvailableInChanges(t *testing.T) {
	doInRepoWithCommit(func(git *git_testing.GitTesting) {
		git.CreateFileWithContents("foo/bar/w", "new contents")
		git.AddAndcommit("*", "added new files")
		repo := RepoLocatedAt(git.Root())
		assert.Len(t, repo.additionsInLastCommit(), 1)
		assert.True(t, strings.HasSuffix(string(repo.AdditionsWithinRange("HEAD~1", "HEAD")[0].Data), "new contents"))
	})
}

func TestOutgoingContentOfModifiedFilesIsAvailableInChanges(t *testing.T) {
	doInRepoWithCommit(func(git *git_testing.GitTesting) {
		git.AppendFileContent("a.txt", "New content.\n", "Spanning multiple lines, even.")
		git.AddAndcommit("a.txt", "added to lorem-ipsum content with my own stuff!")
		repo := RepoLocatedAt(git.Root())
		assert.Len(t, repo.additionsInLastCommit(), 1)
		assert.True(t, strings.HasSuffix(string(repo.additionsInLastCommit()[0].Data), "New content.\nSpanning multiple lines, even."))
	})
}

func TestMultipleOutgoingChangesToTheSameFileAreAvailableInAdditions(t *testing.T) {
	doInRepoWithCommit(func(git *git_testing.GitTesting) {
		git.AppendFileContent("a.txt", "New content.\n")
		git.AddAndcommit("a.txt", "added some new content")

		git.AppendFileContent("a.txt", "More new content.\n")
		git.AddAndcommit("a.txt", "added some more new content")

		repo := RepoLocatedAt(git.Root())
		assert.Len(t, repo.additionsInLastCommit(), 1)
		assert.True(t, strings.HasSuffix(string(repo.AdditionsWithinRange("HEAD~2", "HEAD")[0].Data), "New content.\nMore new content.\n"))
	})
}

func TestContentOfDeletedFilesIsNotAvailableInChanges(t *testing.T) {
	doInRepoWithCommit(func(git *git_testing.GitTesting) {
		git.RemoveFile("a.txt")
		git.AddAndcommit("a.txt", "Deleted this file. After all, it only had lorem-ipsum content.")
		assert.Equal(t, 0, len(RepoLocatedAt(git.Root()).additionsInLastCommit()),
			"There should be no additions because there is only an outgoing deletion")
	})
}

func TestDiffContainingBinaryFileChangesDoesNotBlowUp(t *testing.T) {
	doInRepoWithCommit(func(git *git_testing.GitTesting) {
		repo := RepoLocatedAt(git.Root())
		exec.Command("cp", "./pixel.jpg", repo.root).Run()
		git.AddAndcommit("pixel.jpg", "Testing binary diff.")
		assert.Len(t, repo.additionsInLastCommit(), 1)
		assert.Equal(t, "pixel.jpg", string(repo.additionsInLastCommit()[0].Name))
	})
}

func TestStagedAdditionsIncludeStagedFiles(t *testing.T) {
	doInRepoWithCommit(func(git *git_testing.GitTesting) {
		git.OverwriteFileContent("a.txt", "New content.\n")
		git.Add("a.txt")

		git.AppendFileContent("a.txt", "More new content\n")
		git.AppendFileContent("alice/bob/b.txt", "New content to b\n")

		stagedAdditions := RepoLocatedAt(git.Root()).StagedAdditions()
		assert.Len(t, stagedAdditions, 1)
		assert.Equal(t, "a.txt", string(stagedAdditions[0].Name))
		assert.Equal(t, "New content.\n", string(stagedAdditions[0].Data))
	})
}

func TestStagedAdditionsIncludeStagedNewFiles(t *testing.T) {
	doInRepoWithCommit(func(git *git_testing.GitTesting) {
		git.CreateFileWithContents("new.txt", "New content.\n")
		git.Add("new.txt")

		stagedAdditions := RepoLocatedAt(git.Root()).StagedAdditions()
		assert.Len(t, stagedAdditions, 1)
		assert.Equal(t, "new.txt", string(stagedAdditions[0].Name))
		assert.Equal(t, "New content.\n", string(stagedAdditions[0].Data))
	})
}

func TestStagedAdditionsShouldNotIncludeDeletedFiles(t *testing.T) {
	doInRepoWithCommit(func(git *git_testing.GitTesting) {
		git.RemoveFile("a.txt")
		git.Add(".")

		stagedAdditions := RepoLocatedAt(git.Root()).StagedAdditions()
		assert.Len(t, stagedAdditions, 0)
	})
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

func doInRepoWithCommit(gitOperation git_testing.GitOperation) {
	git_testing.DoInTempGitRepo(func(git *git_testing.GitTesting) {
		git.SetupBaselineFiles("a.txt", filepath.Join("alice", "bob", "b.txt"))
		git.SetupBaselineFiles("c.txt", filepath.Join("folder b", "c.txt"))
		gitOperation(git)
	})
}

func TestGetDiffForStagedFilesWithRenamedFile(t *testing.T) {
	doInRepoWithCommit(func(git *git_testing.GitTesting) {
		root := git.Root()

		// A rename is only reported as one "diff --git a/old b/new" header when
		// the two versions are similar enough, so start from a file with enough
		// lines for one appended line to leave the similarity index high.
		lines := make([]string, 0, 60)
		for i := 0; i < 60; i++ {
			lines = append(lines, fmt.Sprintf("line %d of ordinary content", i))
		}
		git.CreateFileWithContents("big.txt", strings.Join(lines, "\n"))
		git.AddAndcommit("big.txt", "add a file worth renaming")

		runGitIn(t, root, "mv", "big.txt", "renamed.txt")

		f, err := os.OpenFile(filepath.Join(root, "renamed.txt"), os.O_APPEND|os.O_WRONLY, 0o644)
		assert.NoError(t, err)
		_, err = f.WriteString("AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMIK7MDENGbPxRfiCYEXAMPLEKEY\n")
		assert.NoError(t, err)
		assert.NoError(t, f.Close())

		runGitIn(t, root, "add", "-A")

		t.Log("staged diff:\n" + runGitIn(t, root, "diff", "--staged", "--src-prefix=a/", "--dst-prefix=b/"))

		additions := RepoLocatedAt(root).GetDiffForStagedFiles()
		for _, a := range additions {
			t.Logf("addition path=%q data=%q", string(a.Path), string(a.Data))
		}

		var renamed *Addition
		for i := range additions {
			if string(additions[i].Path) == "renamed.txt" {
				renamed = &additions[i]
			}
		}
		if assert.NotNil(t, renamed, "the renamed file produced no addition, so its new content is never scanned") {
			assert.Contains(t, string(renamed.Data), "AWS_SECRET_ACCESS_KEY")
		}
	})
}

func runGitIn(t *testing.T, root string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err, string(out))
	return string(out)
}

func TestMatchGitDiffLine(t *testing.T) {
	tests := []struct {
		name         string
		line         string
		expectMatch  bool
		expectedFile string
	}{
		{"modified file", "diff --git a/a.txt b/a.txt", true, "a.txt"},
		{"path with a space", "diff --git a/folder b/c.txt b/folder b/c.txt", true, "folder b/c.txt"},
		{"renamed file", "diff --git a/big.txt b/renamed.txt", true, "renamed.txt"},
		{"renamed file with differing name lengths", "diff --git a/x.txt b/a-much-longer-name.txt", true, "a-much-longer-name.txt"},
		{"rename whose source contains the separator", "diff --git a/old b/x.txt b/new.txt", true, "new.txt"},
		{"not a diff header", "+++ b/a.txt", false, ""},
		{"index line", "index 0000000..1111111 100644", false, ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			match, filename := MatchGitDiffLine(test.line)
			assert.Equal(t, test.expectMatch, match)
			assert.Equal(t, test.expectedFile, filename)
		})
	}
}
