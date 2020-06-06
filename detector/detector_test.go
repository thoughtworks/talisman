package detector

import (
	"testing"

	"talisman/gitrepo"
	"talisman/talismanrc"

	"github.com/stretchr/testify/assert"
)

func TestEmptyValidationChainPassesAllValidations(t *testing.T) {
	v := NewChain()
	results := NewDetectionResults()
	v.Test(nil, &talismanrc.TalismanRC{}, results)
	assert.False(t, results.HasFailures(), "Empty validation chain is expected to always pass")
}

func TestValidationChainWithFailingValidationAlwaysFails(t *testing.T) {
	v := NewChain()
	v.AddDetector(PassingDetection{})
	v.AddDetector(FailingDetection{})
	results := NewDetectionResults()
	v.Test(nil, &talismanrc.TalismanRC{}, results)

	assert.False(t, results.Successful(), "Expected validation chain with a failure to fail.")
}

type FailingDetection struct{}

func (v FailingDetection) Test(comparator ChecksumCompare, currentAdditions []gitrepo.Addition, ignoreConfig *talismanrc.TalismanRC, result *DetectionResults) {
	result.Fail("some_file", "filecontent", "FAILED BY DESIGN", []string{})
}

type PassingDetection struct{}

func (p PassingDetection) Test(comparator ChecksumCompare, currentAdditions []gitrepo.Addition, ignoreConfig *talismanrc.TalismanRC, result *DetectionResults) {
}
