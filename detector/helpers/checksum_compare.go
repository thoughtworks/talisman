package helpers

import (
	"os"
	"talisman/checksumcalculator"
	"talisman/gitrepo"
	"talisman/talismanrc"
	"talisman/utility"
)

type ChecksumCompare struct {
	calculator checksumcalculator.ChecksumCalculator
	talismanRC *talismanrc.TalismanRC
}

func BuildCC(hasherMode string, talismanRC *talismanrc.TalismanRC, repo gitrepo.GitRepo) *ChecksumCompare {
	wd, _ := os.Getwd()
	hasher := utility.MakeHasher(hasherMode, wd)
	allTrackedFiles := append(repo.TrackedFilesAsAdditions(), repo.StagedAdditions()...)
	calculator := checksumcalculator.NewChecksumCalculator(hasher, allTrackedFiles)
	return &ChecksumCompare{calculator: calculator, talismanRC: talismanRC}
}

// IsScanNotRequired returns true if an Addition's checksum matches one ignored by the .talismanrc file
func (cc *ChecksumCompare) IsScanNotRequired(addition gitrepo.Addition) bool {
	for _, ignore := range cc.talismanRC.IgnoreConfigs {
		if addition.Matches(ignore.GetFileName()) {
			currentCollectiveChecksum := cc.calculator.CalculateCollectiveChecksumForPattern(ignore.GetFileName())
			return ignore.ChecksumMatches(currentCollectiveChecksum)
		}
	}
	return false
}

func (cc *ChecksumCompare) ShouldIgnore(addition gitrepo.Addition, detectorType string) bool {
	return cc.talismanRC.Deny(addition, detectorType) || cc.IsScanNotRequired(addition)
}
