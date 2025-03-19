package filesize

import (
	"talisman/detector/helpers"
	"talisman/detector/severity"
	"testing"

	"talisman/gitrepo"
	"talisman/talismanrc"

	"github.com/stretchr/testify/assert"
)

var talismanRC = &talismanrc.TalismanRC{}
var defaultChecksumCompareUtility = *helpers.BuildCC("default", talismanRC, gitrepo.RepoLocatedAt("."))

func checksumCompareWithTalismanRC(tRC *talismanrc.TalismanRC) *helpers.ChecksumCompare {
	return helpers.BuildCC("default", tRC, gitrepo.RepoLocatedAt("."))
}

func TestShouldFlagLargeFiles(t *testing.T) {
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	content := []byte("more than one byte")
	additions := []gitrepo.Addition{gitrepo.NewAddition("filename", content)}
	NewFileSizeDetector(2).Test(defaultChecksumCompareUtility, additions, talismanRC, results, func() {})
	assert.True(t, results.HasFailures(), "Expected file to fail the check against file size detector.")
}

func TestShouldNotFlagLargeFilesIfThresholdIsBelowSeverity(t *testing.T) {
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	content := []byte("more than one byte")
	talismanRCWithThreshold := &talismanrc.TalismanRC{Threshold: severity.High}
	additions := []gitrepo.Addition{gitrepo.NewAddition("filename", content)}
	NewFileSizeDetector(2).Test(defaultChecksumCompareUtility, additions, talismanRCWithThreshold, results, func() {})
	assert.False(t, results.HasFailures(), "Expected file to not fail the check against file size detector.")
	assert.True(t, results.HasWarnings(), "Expected file to have warnings against file size detector.")
}

func TestShouldNotFlagSmallFiles(t *testing.T) {
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	content := []byte("m")
	additions := []gitrepo.Addition{gitrepo.NewAddition("filename", content)}
	NewFileSizeDetector(2).Test(defaultChecksumCompareUtility, additions, talismanRC, results, func() {})
	assert.False(t, results.HasFailures(), "Expected file to not fail the check against file size detector.")
}

func TestShouldNotFlagIgnoredLargeFiles(t *testing.T) {
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	content := []byte("more than one byte")

	filename := "filename"
	fileIgnoreConfig := &talismanrc.FileIgnoreConfig{
		FileName:        filename,
		IgnoreDetectors: []string{"filesize"},
	}
	talismanRC := &talismanrc.TalismanRC{
		IgnoreConfigs: []talismanrc.IgnoreConfig{fileIgnoreConfig},
	}

	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}
	NewFileSizeDetector(2).Test(*checksumCompareWithTalismanRC(talismanRC), additions, talismanRC, results, func() {})
	assert.True(t, results.Successful(), "expected file %s to be ignored by file size detector", filename)
}
