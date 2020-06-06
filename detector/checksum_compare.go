package detector

import (
	"talisman/checksumcalculator"
	"talisman/gitrepo"
	"talisman/talismanrc"
	"talisman/utility"
)

type ChecksumCompare struct {
	calculator   *checksumcalculator.ChecksumCalculator
	ignoreConfig *talismanrc.TalismanRC
}

//NewChecksumCompare returns new instance of the ChecksumCompare
func NewChecksumCompare(calculator *checksumcalculator.ChecksumCalculator, talismanRCConfig *talismanrc.TalismanRC) *ChecksumCompare {
	cc := ChecksumCompare{calculator: calculator, ignoreConfig: talismanRCConfig}
	return &cc
}

func (cc *ChecksumCompare) IsScanNotRequired(addition gitrepo.Addition) bool {
	currentCollectiveChecksum := utility.CollectiveSHA256Hash([]string{string(addition.Path)})
	declaredCheckSum := ""
	for _, ignore := range cc.ignoreConfig.FileIgnoreConfig {
		if addition.Matches(ignore.FileName) {
			currentCollectiveChecksum = cc.calculator.CalculateCollectiveChecksumForPattern(ignore.FileName)
			declaredCheckSum = ignore.Checksum
		}
	}
	return currentCollectiveChecksum == declaredCheckSum
}
