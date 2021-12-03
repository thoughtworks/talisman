package detector

import (
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

	log "github.com/sirupsen/logrus"
)

//Chain represents a chain of Detectors.
//It is itself a detector.
type Chain struct {
	detectors []detector.Detector
	mode      string
}

//NewChain returns an empty DetectorChain
//It is itself a detector, but it tests nothing.
func NewChain(hooktype string) *Chain {
	result := Chain{[]detector.Detector{}, hooktype}
	return &result
}

//DefaultChain returns a DetectorChain with pre-configured detectors
func DefaultChain(tRC *talismanrc.TalismanRC, runMode string) *Chain {
	chain := NewChain(runMode)
	chain.AddDetector(filename.DefaultFileNameDetector(tRC.Threshold))
	chain.AddDetector(filecontent.NewFileContentDetector(tRC))
	chain.AddDetector(pattern.NewPatternDetector(tRC.CustomPatterns))
	return chain
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
	hasher := utility.MakeHasher(dc.mode, wd)
	hasher.Start()
	calculator := checksumcalculator.NewChecksumCalculator(hasher, append(allAdditions, currentAdditions...))
	cc := helpers.NewChecksumCompare(calculator, hasher, talismanRC)
	log.Printf("Number of files to scan: %d\n", len(currentAdditions))
	log.Printf("Number of detectors: %d\n", len(dc.detectors))
	total := len(currentAdditions) * len(dc.detectors)
	progressBar := utility.GetProgressBar(os.Stdout, "Talisman Scan")
	progressBar.Start(total)
	for _, v := range dc.detectors {
		v.Test(cc, currentAdditions, talismanRC, result, func() {
			progressBar.Increment()
		})
	}
	progressBar.Finish()
	hasher.Shutdown()
}
