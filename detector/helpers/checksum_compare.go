package helpers

import (
	"talisman/checksumcalculator"
	"talisman/gitrepo"
	"talisman/talismanrc"
	"talisman/utility"
)

type ChecksumCompare struct {
	calculator checksumcalculator.ChecksumCalculator
	hasher     utility.SHA256Hasher
	talismanRC *talismanrc.TalismanRC
}

//NewChecksumCompare returns new instance of the ChecksumCompare
func NewChecksumCompare(calculator checksumcalculator.ChecksumCalculator, hasher utility.SHA256Hasher, talismanRCConfig *talismanrc.TalismanRC) ChecksumCompare {
	return ChecksumCompare{calculator: calculator, hasher: hasher, talismanRC: talismanRCConfig}
}

func (cc *ChecksumCompare) IsScanNotRequired(addition gitrepo.Addition) bool {
	currentCollectiveChecksum := cc.hasher.CollectiveSHA256Hash([]string{string(addition.Path)})
	for _, ignore := range cc.talismanRC.IgnoreConfigs {
		if addition.Matches(ignore.GetFileName()) {
			currentCollectiveChecksum = cc.calculator.CalculateCollectiveChecksumForPattern(ignore.GetFileName())
			if ignore.ChecksumMatches(currentCollectiveChecksum) {
				return true
			}
		}
	}
	return false
}
