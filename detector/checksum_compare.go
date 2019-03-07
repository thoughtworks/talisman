package detector

import (
	"talisman/git_repo"
	"talisman/utility"
)

type ChecksumCompare struct {
	additions    []git_repo.Addition
	ignoreConfig TalismanRCIgnore
}

//NewChecksumCompare returns new instance of the ChecksumCompare
func NewChecksumCompare(gitAdditions []git_repo.Addition, talismanRCIgnoreConfig TalismanRCIgnore) *ChecksumCompare {
	cc := ChecksumCompare{additions: gitAdditions, ignoreConfig: talismanRCIgnoreConfig}
	return &cc
}

func (cc *ChecksumCompare) IsScanNotRequired(addition git_repo.Addition) bool {
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

//FilterIgnoresBasedOnChecksums filters the file ignores from the TalismanRCIgnore which doesn't have any checksum value or having mismatched checksum value from the .talsimanrc
func (cc *ChecksumCompare) FilterIgnoresBasedOnChecksums() TalismanRCIgnore {
	finalIgnores := []FileIgnoreConfig{}
	for _, ignore := range cc.ignoreConfig.FileIgnoreConfig {
		currentCollectiveChecksum := cc.calculateCollectiveChecksumForPattern(ignore.FileName, cc.additions)
		// Compare with previous checksum from FileIgnoreConfig
		if ignore.Checksum == currentCollectiveChecksum {
			finalIgnores = append(finalIgnores, ignore)
		}
	}
	rc := TalismanRCIgnore{}
	rc.FileIgnoreConfig = finalIgnores
	return rc
}

func (cc *ChecksumCompare) calculateCollectiveChecksumForPattern(fileNamePattern string, additions []git_repo.Addition) string {
	var patternpaths []string
	currentCollectiveChecksum := ""
	for _, addition := range additions {
		if addition.Matches(fileNamePattern) {
			patternpaths = append(patternpaths, string(addition.Path))
		}
	}
	// Calculate current collective checksum
	patternpaths = utility.UniqueItems(patternpaths)
	if len(patternpaths) != 0 {
		currentCollectiveChecksum = utility.CollectiveSHA256Hash(patternpaths)
	}
	return currentCollectiveChecksum
}
