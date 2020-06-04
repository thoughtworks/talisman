package detector

import (
	"talisman/gitrepo"
	"talisman/talismanrc"
	"talisman/utility"
)

type ChecksumCompare struct {
	additions    []gitrepo.Addition
	ignoreConfig *talismanrc.TalismanRC
}

//NewChecksumCompare returns new instance of the ChecksumCompare
func NewChecksumCompare(gitAdditions []gitrepo.Addition, talismanRCConfig *talismanrc.TalismanRC) *ChecksumCompare {
	cc := ChecksumCompare{additions: gitAdditions, ignoreConfig: talismanRCConfig}
	return &cc
}

func (cc *ChecksumCompare) IsScanNotRequired(addition gitrepo.Addition) bool {
	currentCollectiveChecksum := utility.CollectiveSHA256Hash([]string{string(addition.Path)})
	declaredCheckSum := ""
	for _, ignore := range cc.ignoreConfig.FileIgnoreConfig {
		if addition.Matches(ignore.FileName) {
			currentCollectiveChecksum = utility.CollectiveSHA256Hash([]string{ignore.FileName})
			declaredCheckSum = ignore.Checksum
		}
	}
	return currentCollectiveChecksum == declaredCheckSum
}
