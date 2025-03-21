package checksumcalculator

import (
	"strings"
	"talisman/gitrepo"
	"talisman/talismanrc"
	"talisman/utility"
)

type ChecksumCalculator interface {
	SuggestTalismanRC(fileNamePatterns []string) string
	CalculateCollectiveChecksumForPattern(fileNamePattern string) string
}

type checksumCalculator struct {
	allTrackedFiles []gitrepo.Addition
	hasher          utility.SHA256Hasher
}

// NewChecksumCalculator returns new instance of the CheckSumDetector
func NewChecksumCalculator(hasher utility.SHA256Hasher, gitAdditions []gitrepo.Addition) ChecksumCalculator {
	return &checksumCalculator{hasher: hasher, allTrackedFiles: gitAdditions}
}

// SuggestTalismanRC returns the suggestion for .talismanrc format
func (cc *checksumCalculator) SuggestTalismanRC(fileNamePatterns []string) string {
	var fileIgnoreConfigs []talismanrc.FileIgnoreConfig
	result := strings.Builder{}
	for _, pattern := range fileNamePatterns {
		collectiveChecksum := cc.CalculateCollectiveChecksumForPattern(pattern)
		if collectiveChecksum != "" {
			fileIgnoreConfigs = append(fileIgnoreConfigs, talismanrc.IgnoreFileWithChecksum(pattern, collectiveChecksum))
		}
	}
	if len(fileIgnoreConfigs) != 0 {
		result.WriteString("\n\x1b[33m.talismanrc format for given file names / patterns\x1b[0m\n")
		result.Write([]byte(talismanrc.SuggestRCFor(fileIgnoreConfigs)))
	}
	return result.String()
}

// CalculateCollectiveChecksumForPattern calculates and returns the checksum for files matching the input pattern
func (cc *checksumCalculator) CalculateCollectiveChecksumForPattern(fileNamePattern string) string {
	var patternPaths []string
	currentCollectiveChecksum := ""
	for _, file := range cc.allTrackedFiles {
		if file.Matches(fileNamePattern) {
			patternPaths = append(patternPaths, string(file.Path))
		}
	}
	// Calculate current collective checksum
	patternPaths = utility.UniqueItems(patternPaths)
	if len(patternPaths) != 0 {
		currentCollectiveChecksum = cc.hasher.CollectiveSHA256Hash(patternPaths)
	}
	return currentCollectiveChecksum
}
