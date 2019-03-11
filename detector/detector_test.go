package detector

import (
	"testing"

	"talisman/git_repo"

	"github.com/stretchr/testify/assert"
)

func TestEmptyValidationChainPassesAllValidations(t *testing.T) {
	v := NewChain()
	results := NewDetectionResults()
	v.Test(nil, TalismanRCIgnore{}, results)
	assert.False(t, results.HasFailures(), "Empty validation chain is expected to always pass")
}

func TestValidationChainWithFailingValidationAlwaysFails(t *testing.T) {
	v := NewChain()
	v.AddDetector(PassingDetection{})
	v.AddDetector(FailingDetection{})
	results := NewDetectionResults()
	v.Test(nil, TalismanRCIgnore{}, results)

	assert.False(t, results.Successful(), "Expected validation chain with a failure to fail.")
}

type FailingDetection struct{}

func (v FailingDetection) Test(additions []git_repo.Addition, ignoreConfig TalismanRCIgnore, result *DetectionResults) {
	result.Fail("some_file", "filecontent","FAILED BY DESIGN", []string{})
}

type PassingDetection struct{}

func (p PassingDetection) Test(additions []git_repo.Addition, ignoreConfig TalismanRCIgnore, result *DetectionResults) {
}
