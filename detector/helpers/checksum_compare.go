package helpers

import (
	"talisman/checksumcalculator"
	"talisman/gitrepo"
	"talisman/talismanrc"
)

type ChecksumCompare struct {
	calculator checksumcalculator.ChecksumCalculator
	talismanRC *talismanrc.TalismanRC
}

// NewChecksumCompare returns new instance of the ChecksumCompare
func NewChecksumCompare(calculator checksumcalculator.ChecksumCalculator, talismanRCConfig *talismanrc.TalismanRC) ChecksumCompare {
	return ChecksumCompare{calculator: calculator, talismanRC: talismanRCConfig}
}

func (cc *ChecksumCompare) IsScanNotRequired(addition gitrepo.Addition) bool {
	for _, ignore := range cc.talismanRC.IgnoreConfigs {
		if addition.Matches(ignore.GetFileName()) {
			currentCollectiveChecksum := cc.calculator.CalculateCollectiveChecksumForPattern(ignore.GetFileName())
			return ignore.ChecksumMatches(currentCollectiveChecksum)
		}
	}
	return false
}
