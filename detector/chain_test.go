package detector

import (
	"github.com/stretchr/testify/assert"
	"talisman/detector/helpers"
	"talisman/gitrepo"
	"talisman/talismanrc"
	"testing"
)

type FailingDetection struct{}

func (v FailingDetection) Test(comparator helpers.ChecksumCompare, currentAdditions []gitrepo.Addition, ignoreConfig *talismanrc.TalismanRC, result *helpers.DetectionResults) {
	result.Fail("some_file", "filecontent", "FAILED BY DESIGN", []string{})
}

type PassingDetection struct{}

func (p PassingDetection) Test(comparator helpers.ChecksumCompare, currentAdditions []gitrepo.Addition, ignoreConfig *talismanrc.TalismanRC, result *helpers.DetectionResults) {
}

func TestEmptyValidationChainPassesAllValidations(t *testing.T) {
	v := NewChain()
	results := helpers.NewDetectionResults()
	v.Test(nil, &talismanrc.TalismanRC{}, results)
	assert.False(t, results.HasFailures(), "Empty validation chain is expected to always pass")
}

func TestValidationChainWithFailingValidationAlwaysFails(t *testing.T) {
	v := NewChain()
	v.AddDetector(PassingDetection{})
	v.AddDetector(FailingDetection{})
	results := helpers.NewDetectionResults()
	v.Test(nil, &talismanrc.TalismanRC{}, results)

	assert.False(t, results.Successful(), "Expected validation chain with a failure to fail.")
}
