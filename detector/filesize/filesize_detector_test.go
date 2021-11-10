package filesize

import (
	"talisman/detector/helpers"
	"talisman/detector/severity"
	"talisman/utility"
	"testing"

	"talisman/gitrepo"
	"talisman/talismanrc"

	"github.com/stretchr/testify/assert"
)

var talismanRC = &talismanrc.TalismanRC{}
var defaultHasher = utility.MakeHasher("default", ".")

func TestShouldFlagLargeFiles(t *testing.T) {
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	content := []byte("more than one byte")
	additions := []gitrepo.Addition{gitrepo.NewAddition("filename", content)}
	NewFileSizeDetector(2).Test(helpers.NewChecksumCompare(nil, defaultHasher, talismanRC), additions, talismanRC, results, func() {})
	assert.True(t, results.HasFailures(), "Expected file to fail the check against file size detector.")
}

func TestShouldNotFlagLargeFilesIfThresholdIsBelowSeverity(t *testing.T) {
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	content := []byte("more than one byte")
	talismanRCWithThreshold := &talismanrc.TalismanRC{Threshold: severity.High}
	additions := []gitrepo.Addition{gitrepo.NewAddition("filename", content)}
	NewFileSizeDetector(2).Test(helpers.NewChecksumCompare(nil, defaultHasher, talismanRCWithThreshold), additions, talismanRCWithThreshold, results, func() {})
	assert.False(t, results.HasFailures(), "Expected file to not to fail the check against file size detector.")
	assert.True(t, results.HasWarnings(), "Expected file to have warnings against file size detector.")
}

func TestShouldNotFlagSmallFiles(t *testing.T) {
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	content := []byte("m")
	additions := []gitrepo.Addition{gitrepo.NewAddition("filename", content)}
	NewFileSizeDetector(2).Test(helpers.NewChecksumCompare(nil, defaultHasher, talismanRC), additions, talismanRC, results, func() {})
	assert.False(t, results.HasFailures(), "Expected file to not to fail the check against file size detector.")
}

func TestShouldNotFlagIgnoredLargeFiles(t *testing.T) {
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	content := []byte("more than one byte")

	filename := "filename"
	fileIgnoreConfig := &talismanrc.FileIgnoreConfig{}
	fileIgnoreConfig.FileName = filename
	fileIgnoreConfig.IgnoreDetectors = make([]string, 1)
	fileIgnoreConfig.IgnoreDetectors[0] = "filesize"
	talismanRC := &talismanrc.TalismanRC{}
	talismanRC.IgnoreConfigs = make([]talismanrc.IgnoreConfig, 1)
	talismanRC.IgnoreConfigs[0] = fileIgnoreConfig

	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}
	NewFileSizeDetector(2).Test(helpers.NewChecksumCompare(nil, defaultHasher, talismanRC), additions, talismanRC, results, func() {})
	assert.True(t, results.Successful(), "expected file %s to be ignored by file size detector", filename)
}
