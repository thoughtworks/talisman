package talismanrc

import (
	"regexp"
	"talisman/detector/severity"
)

type PatternString string

type CustomSeverityConfig struct {
	Detector string            `yaml:"detector"`
	Severity severity.Severity `yaml:"severity"`
}

type IgnoreConfig interface {
	isEffective(string) bool
	GetFileName() string
	GetAllowedPatterns() []*regexp.Regexp
	ChecksumMatches(checksum string) bool
}

type FileIgnoreConfig struct {
	compiledPatterns []*regexp.Regexp
	FileName         string   `yaml:"filename"`
	Checksum         string   `yaml:"checksum,omitempty"`
	IgnoreDetectors  []string `yaml:"ignore_detectors,omitempty"`
	AllowedPatterns  []string `yaml:"allowed_patterns,omitempty"`
}

func (i *FileIgnoreConfig) isEffective(detectorName string) bool {
	return !isEmptyString(i.FileName) &&
		contains(i.IgnoreDetectors, detectorName)
}

func (i *FileIgnoreConfig) GetFileName() string {
	return i.FileName
}

func (i *FileIgnoreConfig) ChecksumMatches(incomingChecksum string) bool {
	return i.Checksum == incomingChecksum
}

func (i *FileIgnoreConfig) GetAllowedPatterns() []*regexp.Regexp {
	if i.compiledPatterns == nil {
		i.compiledPatterns = make([]*regexp.Regexp, len(i.AllowedPatterns))
		for idx, p := range i.AllowedPatterns {
			i.compiledPatterns[idx] = regexp.MustCompile(p)
		}
	}
	return i.compiledPatterns
}

type ScanFileIgnoreConfig struct {
	compiledPatterns []*regexp.Regexp
	FileName         string   `yaml:"filename"`
	Checksums        []string `yaml:"checksums,omitempty"`
	IgnoreDetectors  []string `yaml:"ignore_detectors,omitempty"`
	AllowedPatterns  []string `yaml:"allowed_patterns,omitempty"`
}

func (i *ScanFileIgnoreConfig) isEffective(detectorName string) bool {
	return !isEmptyString(i.FileName) &&
		contains(i.IgnoreDetectors, detectorName)
}

func (i *ScanFileIgnoreConfig) GetFileName() string {
	return i.FileName
}

func (i *ScanFileIgnoreConfig) GetAllowedPatterns() []*regexp.Regexp {
	if i.compiledPatterns == nil {
		i.compiledPatterns = make([]*regexp.Regexp, len(i.AllowedPatterns))
		for idx,p := range i.AllowedPatterns {
			i.compiledPatterns[idx] = regexp.MustCompile(p)
		}
	}
	return i.compiledPatterns
}

func (i *ScanFileIgnoreConfig) ChecksumMatches(incomingChecksum string) bool {
	for _, v := range i.Checksums {
		if v == incomingChecksum {
			return true
		}
	}
	return false
}

type ScopeConfig struct {
	ScopeName string `yaml:"scope"`
}

type ExperimentalConfig struct {
	Base64EntropyThreshold float64 `yaml:"base64EntropyThreshold,omitempty"`
}