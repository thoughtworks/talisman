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
