package detector

import (
	"testing"

	"talisman/gitrepo"
	"talisman/talismanrc"

	"github.com/stretchr/testify/assert"
)

func TestShouldFlagLargeFiles(t *testing.T) {
	results := NewDetectionResults()
	content := []byte("more than one byte")
	additions := []gitrepo.Addition{gitrepo.NewAddition("filename", content)}
	NewFileSizeDetector(2).Test(additions, talismanRCIgnore, results)
	assert.True(t, results.HasFailures(), "Expected file to fail the check against file size detector.")
}

func TestShouldNotFlagSmallFiles(t *testing.T) {
	results := NewDetectionResults()
	content := []byte("m")
	additions := []gitrepo.Addition{gitrepo.NewAddition("filename", content)}
	NewFileSizeDetector(2).Test(additions, talismanRCIgnore, results)
	assert.False(t, results.HasFailures(), "Expected file to not to fail the check against file size detector.")
}

func TestShouldNotFlagIgnoredLargeFiles(t *testing.T) {
	results := NewDetectionResults()
	content := []byte("more than one byte")

	filename := "filename"
	fileIgnoreConfig := talismanrc.FileIgnoreConfig{}
	fileIgnoreConfig.FileName = filename
	fileIgnoreConfig.IgnoreDetectors = make([]string, 1)
	fileIgnoreConfig.IgnoreDetectors[0] = "filesize"
	talismanRCIgnore := talismanRCIgnore
	talismanRCIgnore.FileIgnoreConfig = make([]talismanrc.FileIgnoreConfig, 1)
	talismanRCIgnore.FileIgnoreConfig[0] = fileIgnoreConfig

	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}
	NewFileSizeDetector(2).Test(additions, talismanRCIgnore, results)
	assert.True(t, results.Successful(), "expected file %s to be ignored by file size detector", filename)
}
