package helpers

import (
	"os"
	"talisman/checksumcalculator"
	"talisman/gitrepo"
	"talisman/talismanrc"
	"talisman/utility"
)

type IgnoreEvaluator struct {
	calculator checksumcalculator.ChecksumCalculator
	talismanRC *talismanrc.TalismanRC
}

func BuildIgnoreEvaluator(hasherMode string, talismanRC *talismanrc.TalismanRC, repo gitrepo.GitRepo) *IgnoreEvaluator {
	wd, _ := os.Getwd()
	hasher := utility.MakeHasher(hasherMode, wd)
	allTrackedFiles := append(repo.TrackedFilesAsAdditions(), repo.StagedAdditions()...)
	calculator := checksumcalculator.NewChecksumCalculator(hasher, allTrackedFiles)
	return &IgnoreEvaluator{calculator: calculator, talismanRC: talismanRC}
}

// isScanNotRequired returns true if an Addition's checksum matches one ignored by the .talismanrc file
func (ie *IgnoreEvaluator) isScanNotRequired(addition gitrepo.Addition) bool {
	for _, ignore := range ie.talismanRC.IgnoreConfigs {
		if addition.Matches(ignore.GetFileName()) {
			currentCollectiveChecksum := ie.calculator.CalculateCollectiveChecksumForPattern(ignore.GetFileName())
			return ignore.ChecksumMatches(currentCollectiveChecksum)
		}
	}
	return false
}

// ShouldIgnore returns true if the talismanRC indicates that a Detector should ignore an Addition
func (ie *IgnoreEvaluator) ShouldIgnore(addition gitrepo.Addition, detectorType string) bool {
	return ie.talismanRC.Deny(addition, detectorType) || ie.isScanNotRequired(addition)
}
