package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thoughtworks/talisman/git_repo"
)

func TestEmptyValidationChainPassesAllValidations(t *testing.T) {
	v := NewDetectorChain()
	results := NewDetectionResults()
	v.Test(nil, NewIgnores(), results)
	assert.False(t, results.HasFailures(), "Empty validation chain is expected to always pass")
}

func TestValidationChainWithFailingValidationAlwaysFails(t *testing.T) {
	v := NewDetectorChain()
	v.AddDetector(PassingDetection{})
	v.AddDetector(FailingDetection{})
	results := NewDetectionResults()
	v.Test(nil, NewIgnores(), results)

	assert.False(t, results.Successful(), "Expected validation chain with a failure to fail.")
}

type FailingDetection struct{}

func (v FailingDetection) Test(additions []git_repo.Addition, ignores Ignores, result *DetectionResults) {
	result.Fail("some_file", "FAILED BY DESIGN")
}

type PassingDetection struct{}

func (p PassingDetection) Test(additions []git_repo.Addition, ignores Ignores, result *DetectionResults) {
}
