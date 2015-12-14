package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDetectionResultsAreSuccessful(t *testing.T) {
	results := NewDetectionResults()
	assert.True(t, results.Successful(), "New detection result is always expected to succeed")
	assert.False(t, results.HasFailures(), "New detection result is not expected to fail")
}

func TestCallingFailOnDetectionResultsFails(t *testing.T) {
	results := NewDetectionResults()
	results.Fail("some_filename", "Bomb")
	assert.False(t, results.Successful(), "Calling fail on a result should not make it succeed")
	assert.True(t, results.HasFailures(), "Calling fail on a result should make it fail")
}

func TestCanRecordMultipleErrorsAgainstASingleFile(t *testing.T) {
	results := NewDetectionResults()
	results.Fail("some_filename", "Bomb")
	results.Fail("some_filename", "Complete & utter failure")
	results.Fail("another_filename", "Complete & utter failure")
	assert.Len(t, results.Failures("some_filename"), 2, "Expected two errors against some_filename.")
	assert.Len(t, results.Failures("another_filename"), 1, "Expected one error against another_filename")
}

func TestResultsReportsFailures(t *testing.T) {
	results := NewDetectionResults()
	results.Fail("some_filename", "Bomb")
	results.Fail("some_filename", "Complete & utter failure")
	results.Fail("another_filename", "Complete & utter failure")

	actualErrorReport := results.ReportFileFailures("some_filename")
	assert.Regexp(t, "The following errors were detected in some_filename", actualErrorReport, "Error report does not contain expected output")
	assert.Regexp(t, "Bomb", actualErrorReport, "Error report does not contain expected output")
	assert.Regexp(t, "Complete & utter failure", actualErrorReport, "Error report does not contain expected output")
}

func TestLoggingIgnoredFilesDoesNotCauseFailure(t *testing.T) {
	results := NewDetectionResults()
	results.Ignore("some_file", "Ignoring this file, just because")
	results.Ignore("some/other_file", "Ignoring this file too")
	assert.True(t, results.Successful(), "Calling ignore should keep the result successful.")
	assert.True(t, results.HasIgnores(), "Calling ignore should be logged.")
	assert.False(t, results.HasFailures(), "Calling ignore should not cause a result to fail.")

	assert.Regexp(t, "Ignoring this file, just because", results.Report(), "foo")
	assert.Regexp(t, "Ignoring this file too", results.Report(), "foo")
}
