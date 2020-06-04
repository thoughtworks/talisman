package detector

import (
	"talisman/checksumcalculator"
	"talisman/gitrepo"
	"talisman/talismanrc"
	"talisman/utility"
)

type ChecksumCompare struct {
	additions    []gitrepo.Addition
	ignoreConfig *talismanrc.TalismanRC
	allAdditions []gitrepo.Addition
}

//NewChecksumCompare returns new instance of the ChecksumCompare
func NewChecksumCompare(allAdditions []gitrepo.Addition, gitAdditions []gitrepo.Addition, talismanRCConfig *talismanrc.TalismanRC) *ChecksumCompare {
	cc := ChecksumCompare{allAdditions: allAdditions, additions: gitAdditions, ignoreConfig: talismanRCConfig}
	return &cc
}

func (cc *ChecksumCompare) IsScanNotRequired(addition gitrepo.Addition) bool {
	currentCollectiveChecksum := utility.CollectiveSHA256Hash([]string{string(addition.Path)})
	declaredCheckSum := ""
	for _, ignore := range cc.ignoreConfig.FileIgnoreConfig {
		if addition.Matches(ignore.FileName) {
			calculator := checksumcalculator.NewChecksumCalculator(cc.allAdditions)
			currentCollectiveChecksum = calculator.CalculateCollectiveChecksumForPattern(ignore.FileName)
			declaredCheckSum = ignore.Checksum
		}
	}
	return currentCollectiveChecksum == declaredCheckSum
}
