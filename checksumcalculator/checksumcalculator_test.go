package checksumcalculator

import (
	"github.com/stretchr/testify/assert"
	"talisman/gitrepo"
	"talisman/utility"
	"testing"
)

func TestNewChecksumCalculator(t *testing.T) {
	t.Run("should return empty CollectiveChecksum when non existing file name pattern is sent", func(t *testing.T) {
		defaultSHA256Hasher := utility.DefaultSHA256Hasher{}
		gitAdditions := []gitrepo.Addition{
			{
				Path: "GitRepoPath1",
				Name: "GitRepoName1",
			},
		}
		expectedCC := ""
		fileNamePattern := "*NonExistenceFileNamePattern"
		cc := NewChecksumCalculator(defaultSHA256Hasher, gitAdditions)

		actualCC := cc.CalculateCollectiveChecksumForPattern(fileNamePattern)

		assert.Equal(t, expectedCC, actualCC)
	})

	t.Run("should return  CollectiveChecksum when existing file name pattern is sent", func(t *testing.T) {
		defaultSHA256Hasher := utility.DefaultSHA256Hasher{}
		gitAdditions := []gitrepo.Addition{
			{
				Path: "GitRepoPath1",
				Name: "GitRepoName1",
			},
		}
		expectedCC := "54bbf09e5c906e2d7cc0808729f8120cfa3c4bad3fb6a85689ae23ca00e5a3c8"
		fileNamePattern := "*RepoName1"
		cc := NewChecksumCalculator(defaultSHA256Hasher, gitAdditions)

		actualCC := cc.CalculateCollectiveChecksumForPattern(fileNamePattern)

		assert.Equal(t, expectedCC, actualCC)
	})
}

func TestDefaultChecksumCalculator_SuggestTalismanRC(t *testing.T) {
	t.Run("should return no suggestion for .talismanrc format when no matching file name patterns is sent", func(t *testing.T) {
		defaultSHA256Hasher := utility.DefaultSHA256Hasher{}
		gitAdditions := []gitrepo.Addition{
			{
				Path: "GitRepoPath1",
				Name: "GitRepoName1",
			}, {
				Path: "GitRepoPath2",
				Name: "GitRepoName2",
			},
		}
		expectedCC := ""
		fileNamePatterns := []string{"*NonExistenceFileNamePattern1", "*NonExistenceFileNamePattern2", "*NonExistenceFileNamePattern3"}
		cc := NewChecksumCalculator(defaultSHA256Hasher, gitAdditions)

		actualCC := cc.SuggestTalismanRC(fileNamePatterns)

		assert.Equal(t, expectedCC, actualCC)
	})

	t.Run("should return suggestion for .talismanrc format when matching file name patterns is sent", func(t *testing.T) {
		defaultSHA256Hasher := utility.DefaultSHA256Hasher{}
		gitAdditions := []gitrepo.Addition{
			{
				Path: "GitRepoPath1",
				Name: "GitRepoName1",
			}, {
				Path: "GitRepoPath2",
				Name: "GitRepoName2",
			},
		}
		expectedCC := "\n\x1b[33m.talismanrc format for given file names / patterns\x1b[0m\nfileignoreconfig:\n- filename: '*1'\n  checksum: 54bbf09e5c906e2d7cc0808729f8120cfa3c4bad3fb6a85689ae23ca00e5a3c8\n- filename: Git*2\n  checksum: d2103425cbcc10118556d7cd63dd198b3e771bcc0359d6b196ddafcf5b87dac5\nversion: \"1.0\"\n"
		fileNamePatterns := []string{"*1", "Git*2", "*NonExistenceFileNamePattern3"}
		cc := NewChecksumCalculator(defaultSHA256Hasher, gitAdditions)

		actualCC := cc.SuggestTalismanRC(fileNamePatterns)

		assert.Equal(t, expectedCC, actualCC)
	})
}
