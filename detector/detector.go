package detector

import (
	"os"
	"talisman/checksumcalculator"
	"talisman/gitrepo"
	"talisman/talismanrc"
	"talisman/utility"
)

//Detector represents a single kind of test to be performed against a set of Additions
//Detectors are expected to honor the ignores that are passed in and log them in the results
//Detectors are expected to signal any errors to the results
type Detector interface {
	Test(comparator ChecksumCompare, currentAdditions []gitrepo.Addition, ignoreConfig *talismanrc.TalismanRC, result *DetectionResults)
}

//Chain represents a chain of Detectors.
//It is itself a detector.
type Chain struct {
	detectors []Detector
}

//NewChain returns an empty DetectorChain
//It is itself a detector, but it tests nothing.
func NewChain() *Chain {
	result := Chain{make([]Detector, 0)}
	return &result
}

//DefaultChain returns a DetectorChain with pre-configured detectors
func DefaultChain(tRC *talismanrc.TalismanRC) *Chain {
	result := NewChain()
	result.AddDetector(DefaultFileNameDetector())
	result.AddDetector(NewFileContentDetector())
	result.AddDetector(NewPatternDetector(tRC.CustomPatterns))
	return result
}

//AddDetector adds the detector that is passed in to the chain
func (dc *Chain) AddDetector(d Detector) *Chain {
	dc.detectors = append(dc.detectors, d)
	return dc
}

//Test validates the additions against each detector in the chain.
//The results are passed in from detector to detector and thus collect all errors from all detectors
func (dc *Chain) Test(currentAdditions []gitrepo.Addition, talismanRC *talismanrc.TalismanRC, result *DetectionResults) {
	wd, _ := os.Getwd()
	repo := gitrepo.RepoLocatedAt(wd)
	allAdditions := repo.TrackedFilesAsAdditions()
	hasher := utility.DefaultSHA256Hasher{}
	calculator := checksumcalculator.NewChecksumCalculator(hasher, append(allAdditions, currentAdditions...))
	cc := NewChecksumCompare(calculator, hasher, talismanRC)
	for _, v := range dc.detectors {
		v.Test(cc, currentAdditions, talismanRC, result)
	}
}
