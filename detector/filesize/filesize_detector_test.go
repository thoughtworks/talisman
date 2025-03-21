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
var defaultIgnoreEvaluator = helpers.BuildIgnoreEvaluator("default", talismanRC, gitrepo.RepoLocatedAt("."))

func ignoreEvaluatorWithTalismanRC(tRC *talismanrc.TalismanRC) helpers.IgnoreEvaluator {
	return helpers.BuildIgnoreEvaluator("default", tRC, gitrepo.RepoLocatedAt("."))
}

func TestShouldFlagLargeFiles(t *testing.T) {
	results := helpers.NewDetectionResults()
	content := []byte("more than one byte")
	additions := []gitrepo.Addition{gitrepo.NewAddition("filename", content)}
	NewFileSizeDetector(2).Test(defaultIgnoreEvaluator, additions, talismanRC, results, func() {})
	assert.True(t, results.HasFailures(), "Expected file to fail the check against file size detector.")
}

func TestShouldNotFlagLargeFilesIfThresholdIsBelowSeverity(t *testing.T) {
	results := helpers.NewDetectionResults()
	content := []byte("more than one byte")
	talismanRCWithThreshold := &talismanrc.TalismanRC{Threshold: severity.High}
	additions := []gitrepo.Addition{gitrepo.NewAddition("filename", content)}
	NewFileSizeDetector(2).Test(defaultIgnoreEvaluator, additions, talismanRCWithThreshold, results, func() {})
	assert.False(t, results.HasFailures(), "Expected file to not fail the check against file size detector.")
	assert.True(t, results.HasWarnings(), "Expected file to have warnings against file size detector.")
}

func TestShouldNotFlagSmallFiles(t *testing.T) {
	results := helpers.NewDetectionResults()
	content := []byte("m")
	additions := []gitrepo.Addition{gitrepo.NewAddition("filename", content)}
	NewFileSizeDetector(2).Test(defaultIgnoreEvaluator, additions, talismanRC, results, func() {})
	assert.False(t, results.HasFailures(), "Expected file to not fail the check against file size detector.")
}

func TestShouldNotFlagIgnoredLargeFiles(t *testing.T) {
	results := helpers.NewDetectionResults()
	content := []byte("more than one byte")

	filename := "filename"
	fileIgnoreConfig := talismanrc.FileIgnoreConfig{
		FileName:        filename,
		IgnoreDetectors: []string{"filesize"},
	}
	talismanRC := &talismanrc.TalismanRC{
		FileIgnoreConfig: []talismanrc.FileIgnoreConfig{fileIgnoreConfig},
	}

	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}
	NewFileSizeDetector(2).Test(ignoreEvaluatorWithTalismanRC(talismanRC), additions, talismanRC, results, func() {})
	assert.True(t, results.Successful(), "expected file %s to be ignored by file size detector", filename)
}
