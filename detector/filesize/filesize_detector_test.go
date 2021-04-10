package filesize

import (
	"talisman/detector/helpers"
	"talisman/utility"
	"testing"

	"talisman/gitrepo"
	"talisman/talismanrc"

	"github.com/stretchr/testify/assert"
)

var talismanRC = &talismanrc.TalismanRC{}

func TestShouldFlagLargeFiles(t *testing.T) {
	results := helpers.NewDetectionResults(talismanrc.Hook)
	content := []byte("more than one byte")
	additions := []gitrepo.Addition{gitrepo.NewAddition("filename", content)}
	NewFileSizeDetector(2).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additions, talismanRC, results, func() {})
	assert.True(t, results.HasFailures(), "Expected file to fail the check against file size detector.")
}

func TestShouldNotFlagLargeFilesIfThresholdIsBelowSeverity(t *testing.T) {
	results := helpers.NewDetectionResults(talismanrc.Hook)
	content := []byte("more than one byte")
	var talismanRCContents = "threshold: high"
	talismanRCWithThreshold := talismanrc.NewTalismanRC([]byte(talismanRCContents))
	additions := []gitrepo.Addition{gitrepo.NewAddition("filename", content)}
	NewFileSizeDetector(2).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanRCWithThreshold), additions, talismanRCWithThreshold, results, func() {})
	assert.False(t, results.HasFailures(), "Expected file to not to fail the check against file size detector.")
	assert.True(t, results.HasWarnings(), "Expected file to have warnings against file size detector.")
}

func TestShouldNotFlagSmallFiles(t *testing.T) {
	results := helpers.NewDetectionResults(talismanrc.Hook)
	content := []byte("m")
	additions := []gitrepo.Addition{gitrepo.NewAddition("filename", content)}
	NewFileSizeDetector(2).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additions, talismanRC, results, func() {})
	assert.False(t, results.HasFailures(), "Expected file to not to fail the check against file size detector.")
}

func TestShouldNotFlagIgnoredLargeFiles(t *testing.T) {
	results := helpers.NewDetectionResults(talismanrc.Hook)
	content := []byte("more than one byte")

	filename := "filename"
	fileIgnoreConfig := talismanrc.FileIgnoreConfig{}
	fileIgnoreConfig.FileName = filename
	fileIgnoreConfig.IgnoreDetectors = make([]string, 1)
	fileIgnoreConfig.IgnoreDetectors[0] = "filesize"
	talismanRC := talismanRC
	talismanRC.FileIgnoreConfig = make([]talismanrc.FileIgnoreConfig, 1)
	talismanRC.FileIgnoreConfig[0] = fileIgnoreConfig

	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}
	NewFileSizeDetector(2).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additions, talismanRC, results, func() {})
	assert.True(t, results.Successful(), "expected file %s to be ignored by file size detector", filename)
}
