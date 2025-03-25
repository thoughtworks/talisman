package talismanrc

import (
	"regexp"
	"talisman/detector/severity"

	logr "github.com/sirupsen/logrus"
)

type PatternString string

type Pattern struct {
	*regexp.Regexp
}

func (p Pattern) MarshalYAML() (interface{}, error) {
	return p.String(), nil
}

func (p *Pattern) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err != nil {
		logr.Errorf("Pattern.UmarshalYAML error: %v", err)
		return err
	}
	*p = Pattern{regexp.MustCompile(s)}
	return nil
}

type CustomSeverityConfig struct {
	Detector string            `yaml:"detector"`
	Severity severity.Severity `yaml:"severity"`
}

type FileIgnoreConfig struct {
	FileName        string   `yaml:"filename"`
	Checksum        string   `yaml:"checksum,omitempty"`
	IgnoreDetectors []string `yaml:"ignore_detectors,omitempty"`
	AllowedPatterns []string `yaml:"allowed_patterns,omitempty"`

	compiledPatterns []*regexp.Regexp
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

func IgnoreFileWithChecksum(filename, checksum string) FileIgnoreConfig {
	return FileIgnoreConfig{FileName: filename, Checksum: checksum}
}

type ScopeConfig struct {
	ScopeName string `yaml:"scope"`
}

type ExperimentalConfig struct {
	Base64EntropyThreshold float64 `yaml:"base64EntropyThreshold,omitempty"`
}
