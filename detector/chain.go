package detector

import (
	"os"
	"talisman/detector/detector"
	"talisman/detector/filecontent"
	"talisman/detector/filename"
	"talisman/detector/helpers"
	"talisman/detector/pattern"
	"talisman/gitrepo"
	"talisman/talismanrc"
	"talisman/utility"

	log "github.com/sirupsen/logrus"
)

// Chain represents a chain of Detectors.
// It is itself a detector.
type Chain struct {
	detectors       []detector.Detector
	ignoreEvaluator helpers.IgnoreEvaluator
}

// NewChain returns an empty DetectorChain
// It is itself a detector, but it tests nothing.
func NewChain(ignoreEvaluator helpers.IgnoreEvaluator) *Chain {
	result := Chain{[]detector.Detector{}, ignoreEvaluator}
	return &result
}

// DefaultChain returns a DetectorChain with pre-configured detectors
func DefaultChain(tRC *talismanrc.TalismanRC, ignoreEvaluator helpers.IgnoreEvaluator) *Chain {
	chain := NewChain(ignoreEvaluator)
	chain.AddDetector(filename.DefaultFileNameDetector(tRC.Threshold))
	chain.AddDetector(filecontent.NewFileContentDetector(tRC))
	chain.AddDetector(pattern.NewPatternDetector(tRC.CustomPatterns))
	return chain
}

// AddDetector adds the detector that is passed in to the chain
func (dc *Chain) AddDetector(d detector.Detector) *Chain {
	dc.detectors = append(dc.detectors, d)
	return dc
}

// Test validates the additions against each detector in the chain.
// The results are passed in from detector to detector and thus collect all errors from all detectors
func (dc *Chain) Test(additions []gitrepo.Addition, talismanRC *talismanrc.TalismanRC, result *helpers.DetectionResults) {
	log.Printf("Number of files to scan: %d\n", len(additions))
	log.Printf("Number of detectors: %d\n", len(dc.detectors))
	total := len(additions) * len(dc.detectors)
	progressBar := utility.GetProgressBar(os.Stdout, "Talisman Scan")
	progressBar.Start(total)
	for _, v := range dc.detectors {
		v.Test(dc.ignoreEvaluator, additions, talismanRC, result, func() {
			progressBar.Increment()
		})
	}
	progressBar.Finish()
}
