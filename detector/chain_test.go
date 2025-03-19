package detector

import (
	"io/ioutil"
	"talisman/detector/filecontent"
	"talisman/detector/filename"
	"talisman/detector/helpers"
	"talisman/detector/pattern"
	"talisman/detector/severity"
	"talisman/gitrepo"
	"talisman/talismanrc"
	"testing"

	logr "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	logr.SetOutput(ioutil.Discard)
}

type FailingDetection struct{}

func (v FailingDetection) Test(comparator helpers.IgnoreEvaluator, currentAdditions []gitrepo.Addition, ignoreConfig *talismanrc.TalismanRC, result *helpers.DetectionResults, additionCompletionCallback func()) {
	result.Fail("some_file", "filecontent", "FAILED BY DESIGN", []string{}, severity.Low)
}

type PassingDetection struct{}

func (p PassingDetection) Test(comparator helpers.IgnoreEvaluator, currentAdditions []gitrepo.Addition, ignoreConfig *talismanrc.TalismanRC, result *helpers.DetectionResults, additionCompletionCallback func()) {
}

func TestEmptyValidationChainPassesAllValidations(t *testing.T) {
	ie := helpers.BuildIgnoreEvaluator("pre-push", nil, gitrepo.RepoLocatedAt("."))
	v := NewChain(ie)
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	v.Test(nil, &talismanrc.TalismanRC{}, results)
	assert.False(t, results.HasFailures(), "Empty validation chain is expected to always pass")
}

func TestValidationChainWithFailingValidationAlwaysFails(t *testing.T) {
	ie := helpers.BuildIgnoreEvaluator("pre-push", nil, gitrepo.RepoLocatedAt("."))
	v := NewChain(ie)
	v.AddDetector(PassingDetection{})
	v.AddDetector(FailingDetection{})
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	v.Test(nil, &talismanrc.TalismanRC{}, results)

	assert.False(t, results.Successful(), "Expected validation chain with a failure to fail.")
}

func TestDefaultChainShouldCreateChainSpecifiedModeAndPresetDetectors(t *testing.T) {
	talismanRC := &talismanrc.TalismanRC{
		Threshold:      severity.Medium,
		CustomPatterns: []talismanrc.PatternString{"AKIA*"},
	}
	ie := helpers.BuildIgnoreEvaluator("pre-push", talismanRC, gitrepo.RepoLocatedAt("."))
	v := DefaultChain(talismanRC, ie)
	assert.Equal(t, 3, len(v.detectors))

	defaultFileNameDetector := filename.DefaultFileNameDetector(talismanRC.Threshold)
	assert.Equal(t, defaultFileNameDetector, v.detectors[0])

	expectedFileContentDetector := filecontent.NewFileContentDetector(talismanRC)
	assert.Equal(t, expectedFileContentDetector, v.detectors[1])

	expectedPatternDetector := pattern.NewPatternDetector(talismanRC.CustomPatterns)
	assert.Equal(t, expectedPatternDetector, v.detectors[2])
}
