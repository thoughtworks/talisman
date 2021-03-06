package checksumcalculator

import (
	"fmt"
	"talisman/gitrepo"
	"talisman/talismanrc"
	"talisman/utility"

	yaml "gopkg.in/yaml.v2"
)

type ChecksumCalculator interface {
	SuggestTalismanRC(fileNamePatterns []string) string
	CalculateCollectiveChecksumForPattern(fileNamePattern string) string
}

type DefaultChecksumCalculator struct {
	gitAdditions []gitrepo.Addition
	hasher       utility.SHA256Hasher
}

//NewChecksumCalculator returns new instance of the CheckSumDetector
func NewChecksumCalculator(hasher utility.SHA256Hasher, gitAdditions []gitrepo.Addition) ChecksumCalculator {
	cc := DefaultChecksumCalculator{hasher: hasher, gitAdditions: gitAdditions}
	return &cc
}

//SuggestTalismanRC returns the suggestion for .talismanrc format
func (cc *DefaultChecksumCalculator) SuggestTalismanRC(fileNamePatterns []string) string {
	var fileIgnoreConfigs []talismanrc.FileIgnoreConfig
	result := ""
	for _, pattern := range fileNamePatterns {
		collectiveChecksum := cc.CalculateCollectiveChecksumForPattern(pattern)
		if collectiveChecksum != "" {
			fileIgnoreConfig := talismanrc.FileIgnoreConfig{FileName: pattern, Checksum: collectiveChecksum, IgnoreDetectors: []string{}}
			fileIgnoreConfigs = append(fileIgnoreConfigs, fileIgnoreConfig)
		}
	}
	if len(fileIgnoreConfigs) != 0 {
		result = result + fmt.Sprintf("\n\x1b[33m.talismanrc format for given file names / patterns\x1b[0m\n")
		talismanRCConfig := talismanrc.MakeWithFileIgnores(fileIgnoreConfigs)
		m, _ := yaml.Marshal(&talismanRCConfig)
		result = result + string(m)
	}
	return result
}

func (cc *DefaultChecksumCalculator) CalculateCollectiveChecksumForPattern(fileNamePattern string) string {
	var patternPaths []string
	currentCollectiveChecksum := ""
	for _, addition := range cc.gitAdditions {
		if addition.Matches(fileNamePattern) {
			patternPaths = append(patternPaths, string(addition.Path))
		}
	}
	// Calculate current collective checksum
	patternPaths = utility.UniqueItems(patternPaths)
	if len(patternPaths) != 0 {
		currentCollectiveChecksum = cc.hasher.CollectiveSHA256Hash(patternPaths)
	}
	return currentCollectiveChecksum
}
