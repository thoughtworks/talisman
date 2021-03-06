package talismanrc

import (
	"talisman/detector/severity"
)

type CustomSeverityConfig struct {
	Detector string `yaml:"detector"`
	Severity severity.Severity `yaml:"severity"`
}

type IgnoreConfig interface {
	isEffective(string) bool
}

type FileIgnoreConfig struct {
	FileName        string   `yaml:"filename"`
	Checksum        string   `yaml:"checksum,omitempty"`
	IgnoreDetectors []string `yaml:"ignore_detectors,omitempty"`
	AllowedPatterns []string `yaml:"allowed_patterns,omitempty"`
}

func (i *FileIgnoreConfig) isEffective(detectorName string) bool {
	return !isEmptyString(i.FileName) &&
		contains(i.IgnoreDetectors, detectorName)
}

type ScanFileIgnoreConfig struct {
	FileName        string   `yaml:"filename"`
	Checksums       []string `yaml:"checksums,omitempty"`
	IgnoreDetectors []string `yaml:"ignore_detectors,omitempty"`
	AllowedPatterns []string `yaml:"allowed_patterns,omitempty"`
}

func (i *ScanFileIgnoreConfig) isEffective(detectorName string) bool {
	return !isEmptyString(i.FileName) &&
		contains(i.IgnoreDetectors, detectorName)
}

type ScopeConfig struct {
	ScopeName string `yaml:"scope"`
}

type ExperimentalConfig struct {
	Base64EntropyThreshold float64 `yaml:"base64EntropyThreshold,omitempty"`
}

type PatternString string

