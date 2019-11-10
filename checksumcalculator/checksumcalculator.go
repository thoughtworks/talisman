package checksumcalculator

import (
	"fmt"
	"os"
	"talisman/gitrepo"
	"talisman/talismanrc"
	"talisman/utility"

	yaml "gopkg.in/yaml.v2"
)

type ChecksumCalculator struct {
	fileNamePatterns []string
}

//NewChecksumCalculator returns new instance of the CheckSumDetector
func NewChecksumCalculator(patterns []string) *ChecksumCalculator {
	cc := ChecksumCalculator{fileNamePatterns: patterns}
	return &cc
}

//SuggestTalismanRC returns the suggestion for .talismanrc format
func (cc *ChecksumCalculator) SuggestTalismanRC() string {
	wd, _ := os.Getwd()
	repo := gitrepo.RepoLocatedAt(wd)
	gitTrackedFilesAsAdditions := repo.TrackedFilesAsAdditions()
	//Adding staged files for calculation
	gitTrackedFilesAsAdditions = append(gitTrackedFilesAsAdditions, repo.StagedAdditions()...)
	var fileIgnoreConfigs []talismanrc.FileIgnoreConfig
	result := ""
	for _, pattern := range cc.fileNamePatterns {
		collectiveChecksum := cc.calculateCollectiveChecksumForPattern(pattern, gitTrackedFilesAsAdditions)
		if collectiveChecksum != "" {
			fileIgnoreConfig := talismanrc.FileIgnoreConfig{FileName: pattern, Checksum: collectiveChecksum, IgnoreDetectors: []string{}}
			fileIgnoreConfigs = append(fileIgnoreConfigs, fileIgnoreConfig)
		}
	}
	if len(fileIgnoreConfigs) != 0 {
		result = result + fmt.Sprintf("\n\x1b[33m.talismanrc format for given file names / patterns\x1b[0m\n")
		talismanRCIgnoreConfig := talismanrc.TalismanRCIgnore{FileIgnoreConfig: fileIgnoreConfigs}
		m, _ := yaml.Marshal(&talismanRCIgnoreConfig)
		result = result + string(m)
	}
	return result
}

func (cc *ChecksumCalculator) calculateCollectiveChecksumForPattern(fileNamePattern string, additions []gitrepo.Addition) string {
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
