package checksumcalculator

import (
	"talisman/gitrepo"
	"talisman/utility"
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultSHA256Hasher utility.SHA256Hasher

func init() {
	defaultSHA256Hasher = utility.MakeHasher("default", ".")
}

func TestNewChecksumCalculator(t *testing.T) {
	t.Run("should return empty CollectiveChecksum when non existing file name pattern is sent", func(t *testing.T) {
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
		gitAdditions := []gitrepo.Addition{
			{
				Path: "GitRepoPath1/GitRepoName1",
				Name: "GitRepoName1",
			},
		}
		expectedCC := "19250a996e1200d33e91454bc662efae7682410e5347cfc56b0ff386dfbc10ae"
		fileNamePattern := "GitRepoPath1/"
		cc := NewChecksumCalculator(defaultSHA256Hasher, gitAdditions)

		actualCC := cc.CalculateCollectiveChecksumForPattern(fileNamePattern)

		assert.Equal(t, expectedCC, actualCC)
	})

	t.Run("should return the files own CollectiveChecksum when same file name is present in subfolders", func(t *testing.T) {
		gitAdditions := []gitrepo.Addition{
			{
				Path: "hello.txt",
				Name: "hello.txt",
			},
			{
				Path: "subfolder/hello.txt",
				Name: "hello.txt",
			},
		}
		hello_expectedCC := "9d30c2e4bcf181bba07374cc416f1892d89918038bce5172776475347c4d2d69"
		fileNamePattern1 := "hello.txt"
		cc1 := NewChecksumCalculator(defaultSHA256Hasher, gitAdditions)
		actualCC := cc1.CalculateCollectiveChecksumForPattern(fileNamePattern1)
		assert.Equal(t, hello_expectedCC, actualCC)

		subfolder_hello_expectedCC := "6c779c16bcc2e63c659be7649a531650210d6b96ae590a146f9ccdca383587f6"
		fileNamePattern2 := "subfolder/hello.txt"
		cc2 := NewChecksumCalculator(defaultSHA256Hasher, gitAdditions)
		subfolder_actualCC := cc2.CalculateCollectiveChecksumForPattern(fileNamePattern2)
		assert.Equal(t, subfolder_hello_expectedCC, subfolder_actualCC)

		all_txt_expectedCC := "aba77c9077539130e21a8f275fc5f1f43b7f0589e392bc89e4b96c578f0a9184"
		fileNamePattern3 := "*.txt"
		cc3 := NewChecksumCalculator(defaultSHA256Hasher, gitAdditions)
		all_txt_actualCC := cc3.CalculateCollectiveChecksumForPattern(fileNamePattern3)
		assert.Equal(t, all_txt_expectedCC, all_txt_actualCC)
	})
}

func TestDefaultChecksumCalculator_SuggestTalismanRC(t *testing.T) {
	t.Run("should return no suggestion for .talismanrc format when no matching file name patterns is sent", func(t *testing.T) {
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
