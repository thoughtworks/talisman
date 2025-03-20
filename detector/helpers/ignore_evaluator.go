package helpers

import (
	"os"
	"talisman/checksumcalculator"
	"talisman/gitrepo"
	"talisman/talismanrc"
	"talisman/utility"
)

type IgnoreEvaluator interface {
	ShouldIgnore(addition gitrepo.Addition, detectorType string) bool
}

type scanAllAdditions struct{}

// Returns an IgnoreEvaluator that forces all files to be scanned, such as when scanning the history of a repo
func ScanHistoryEvaluator() IgnoreEvaluator {
	return &scanAllAdditions{}
}

// Returns false so that all additions are scanned
func (ie *scanAllAdditions) ShouldIgnore(gitrepo.Addition, string) bool {
	return false
}

type ignoreEvaluator struct {
	calculator checksumcalculator.ChecksumCalculator
	talismanRC *talismanrc.TalismanRC
}

// Returns an IgnoreEvaluator around the rules defined in the current .talismanrc file
func BuildIgnoreEvaluator(hasherMode string, talismanRC *talismanrc.TalismanRC, repo gitrepo.GitRepo) IgnoreEvaluator {
	wd, _ := os.Getwd()
	hasher := utility.MakeHasher(hasherMode, wd)
	allTrackedFiles := append(repo.TrackedFilesAsAdditions(), repo.StagedAdditions()...)
	calculator := checksumcalculator.NewChecksumCalculator(hasher, allTrackedFiles)
	return &ignoreEvaluator{calculator: calculator, talismanRC: talismanRC}
}

// ShouldIgnore returns true if the talismanRC indicates that a Detector should ignore an Addition
func (ie *ignoreEvaluator) ShouldIgnore(addition gitrepo.Addition, detectorType string) bool {
	return ie.talismanRC.Deny(addition, detectorType) || ie.isScanNotRequired(addition)
}

// isScanNotRequired returns true if an Addition's checksum matches one ignored by the .talismanrc file
func (ie *ignoreEvaluator) isScanNotRequired(addition gitrepo.Addition) bool {
	for _, ignore := range ie.talismanRC.IgnoreConfigs {
		if addition.Matches(ignore.GetFileName()) {
			currentCollectiveChecksum := ie.calculator.CalculateCollectiveChecksumForPattern(ignore.GetFileName())
			return ignore.ChecksumMatches(currentCollectiveChecksum)
		}
	}
	return false
}
