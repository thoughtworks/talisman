package detector

import (
	"testing"

	"talisman/git_repo"

	"github.com/stretchr/testify/assert"
)

func TestShouldFlagLargeFiles(t *testing.T) {
	results := NewDetectionResults()
	content := []byte("more than one byte")
	additions := []git_repo.Addition{git_repo.NewAddition("filename", content)}
	NewFileSizeDetector(2).Test(additions, TalismanRCIgnore{}, results)
	assert.True(t, results.HasFailures(), "Expected file to fail the check against file size detector.")
}

func TestShouldNotFlagSmallFiles(t *testing.T) {
	results := NewDetectionResults()
	content := []byte("m")
	additions := []git_repo.Addition{git_repo.NewAddition("filename", content)}
	NewFileSizeDetector(2).Test(additions, TalismanRCIgnore{}, results)
	assert.False(t, results.HasFailures(), "Expected file to not to fail the check against file size detector.")
}

func TestShouldNotFlagIgnoredLargeFiles(t *testing.T) {
	results := NewDetectionResults()
	content := []byte("more than one byte")

	filename := "filename"
	fileIgnoreConfig := FileIgnoreConfig{}
	fileIgnoreConfig.FileName = filename
	talismanRCIgnore := TalismanRCIgnore{}
	talismanRCIgnore.FileIgnoreConfig = make([]FileIgnoreConfig, 1)
	talismanRCIgnore.FileIgnoreConfig[0] = fileIgnoreConfig

	additions := []git_repo.Addition{git_repo.NewAddition(filename, content)}
	NewFileSizeDetector(2).Test(additions, talismanRCIgnore, results)
	assert.True(t, results.Successful(), "expected file %s to be ignored by file size detector", filename)
}
