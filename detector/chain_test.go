package detector

import (
	"io/ioutil"
	"talisman/detector/helpers"
	"talisman/detector/severity"
	"talisman/gitrepo"
	"talisman/talismanrc"
	"testing"

	"github.com/stretchr/testify/assert"
	logr "github.com/Sirupsen/logrus"
)


func init() {
	logr.SetOutput(ioutil.Discard)
}
type FailingDetection struct{}

func (v FailingDetection) Test(comparator helpers.ChecksumCompare, currentAdditions []gitrepo.Addition, ignoreConfig *talismanrc.TalismanRC, result *helpers.DetectionResults, additionCompletionCallback func()) {
	result.Fail("some_file", "filecontent", "FAILED BY DESIGN", []string{}, severity.Low)
}

type PassingDetection struct{}

func (p PassingDetection) Test(comparator helpers.ChecksumCompare, currentAdditions []gitrepo.Addition, ignoreConfig *talismanrc.TalismanRC, result *helpers.DetectionResults, additionCompletionCallback func()) {
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
