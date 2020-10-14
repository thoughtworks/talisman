package detector

import (
	log "github.com/Sirupsen/logrus"
	"github.com/cheggaaa/pb/v3"
	"os"
	"talisman/checksumcalculator"
	"talisman/detector/detector"
	"talisman/detector/filecontent"
	"talisman/detector/filename"
	"talisman/detector/helpers"
	"talisman/detector/pattern"
	"talisman/gitrepo"
	"talisman/talismanrc"
	"talisman/utility"
)

//Chain represents a chain of Detectors.
//It is itself a detector.
type Chain struct {
	detectors []detector.Detector
}

//NewChain returns an empty DetectorChain
//It is itself a detector, but it tests nothing.
func NewChain() *Chain {
	result := Chain{make([]detector.Detector, 0)}
	return &result
}

//DefaultChain returns a DetectorChain with pre-configured detectors
func DefaultChain(tRC *talismanrc.TalismanRC) *Chain {
	result := NewChain()
	result.AddDetector(filename.DefaultFileNameDetector(tRC.Threshold))
	result.AddDetector(filecontent.NewFileContentDetector(tRC))
	result.AddDetector(pattern.NewPatternDetector(tRC.CustomPatterns))
	return result
}

//AddDetector adds the detector that is passed in to the chain
func (dc *Chain) AddDetector(d detector.Detector) *Chain {
	dc.detectors = append(dc.detectors, d)
	return dc
}

//Test validates the additions against each detector in the chain.
//The results are passed in from detector to detector and thus collect all errors from all detectors
func (dc *Chain) Test(currentAdditions []gitrepo.Addition, talismanRC *talismanrc.TalismanRC, result *helpers.DetectionResults) {
	wd, _ := os.Getwd()
	repo := gitrepo.RepoLocatedAt(wd)
	allAdditions := repo.TrackedFilesAsAdditions()
	hasher := utility.DefaultSHA256Hasher{}
	calculator := checksumcalculator.NewChecksumCalculator(hasher, append(allAdditions, currentAdditions...))
	cc := helpers.NewChecksumCompare(calculator, hasher, talismanRC)
	log.Printf("Number of files to scan: %d\n", len(currentAdditions))
	log.Printf("Number of detectors: %d\n", len(dc.detectors))
	total := len(currentAdditions) * len(dc.detectors)
	bar := pb.StartNew(total)
	for _, v := range dc.detectors {
		v.Test(cc, currentAdditions, talismanRC, result, func() {
			bar.Increment()
		})
	}
	bar.Finish()
}
