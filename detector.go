package main

import (
	"fmt"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/thoughtworks/talisman/git_repo"
)

//Detector represents a single kind of test to be performed against a set of Additions
//Detectors are expected to honor the ignores that are passed in and log them in the results
//Detectors are expected to signal any errors to the results
type Detector interface {
	Test(additions []git_repo.Addition, ignores Ignores, result *DetectionResults)
}

//FileNameDetector represents tests performed against the fileName of the Additions.
//The Paths of the supplied Additions are tested against the configured patterns and if any of them match, it is logged as a failure during the run
type FileNameDetector struct {
	flagPatterns []*regexp.Regexp
}

//DefaultFileNameDetector returns a FileNameDetector that tests Additions against the pre-configured patterns
func DefaultFileNameDetector() Detector {
	return NewFileNameDetector("^.+_rsa$",
		"^.+_dsa$",
		"^.+_ed25519$",
		"^.+_ecdsa$",
		"^\\.\\w+_history$",
		"^.+\\.pem$",
		"^.+\\.ppk$",
		"^.+\\.key(pair)?$",
		"^.+\\.pkcs12$",
		"^.+\\.pfx$",
		"^.+\\.p12$",
		"^.+\\.asc$",
		"^\\.?htpasswd$",
		"^\\.?netrc$",
		"^.*\\.tblk$",
		"^.*\\.ovpn$",
		"^.*\\.kdb$",
		"^.*\\.agilekeychain$",
		"^.*\\.keychain$",
		"^.*\\.key(store|ring)$",
		"^jenkins\\.plugins\\.publish_over_ssh\\.BapSshPublisherPlugin.xml$",
		"^credentials\\.xml$",
		"^.*\\.pubxml(\\.user)?$",
		"^\\.?s3cfg$",
		"^.*\\.ovpn$",
		"^\\.gitrobrc$",
		"^\\.?(bash|zsh)rc$",
		"^\\.?(bash_|zsh_)?profile$",
		"^\\.?(bash_|zsh_)?aliases$",
		"^secret_token.rb$",
		"^omniauth.rb$",
		"^carrierwave.rb$",
		"^schema.rb$",
		"^database.yml$",
		"^settings.py$",
		"^.*(config)(\\.inc)?\\.php$",
		"^LocalSettings.php$",
		"\\.?env",
		"\\bdump|dump\\b",
		"\\bsql|sql\\b",
		"\\bdump|dump\\b",
		"password",
		"backup",
		"private.*key",
		"(oauth).*(token)",
		"^.*\\.log$",
		"^\\.?kwallet$",
		"^\\.?gnucash$")
}

//NewFileNameDetector returns a FileNameDetector that tests Additions against the supplied patterns
func NewFileNameDetector(patternStrings ...string) Detector {
	var patterns = make([]*regexp.Regexp, len(patternStrings))
	for i, p := range patternStrings {
		patterns[i], _ = regexp.Compile(p)
	}
	return FileNameDetector{patterns}
}

//Test tests the fileNames of the Additions to ensure that they don't look suspicious
func (fd FileNameDetector) Test(additions []git_repo.Addition, ignores Ignores, result *DetectionResults) {
	for _, addition := range additions {
		if ignores.Deny(addition) {
			log.WithFields(log.Fields{
				"filePath": addition.Path,
			}).Info("Ignoring addition as it was specified to be ignored.")
			result.Ignore(addition.Path, fmt.Sprintf("%s was ignored by .talismanignore", addition.Path))
			continue
		}
		for _, pattern := range fd.flagPatterns {
			if addition.Name.Matches(pattern) {
				log.WithFields(log.Fields{
					"filePath": addition.Path,
					"pattern":  pattern,
				}).Info("Failing file as it matched pattern.")
				result.Fail(addition.Path, fmt.Sprintf("The file name %q failed checks against the pattern %s", addition.Path, pattern))
			}
		}
	}
}

//DetectorChain represents a chain of Detectors.
//It is itself a detector.
type DetectorChain struct {
	detectors []Detector
}

//NewDetectorChain returns an empty DetectorChain
//It is itself a detector, but it tests nothing.
func NewDetectorChain() *DetectorChain {
	result := DetectorChain{make([]Detector, 0)}
	return &result
}

//DefaultDetectorChain returns a DetectorChain with pre-configured detectors
func DefaultDetectorChain() *DetectorChain {
	result := NewDetectorChain()
	result.AddDetector(DefaultFileNameDetector())
	return result
}

//AddDetector adds the detector that is passed in to the chain
func (dc *DetectorChain) AddDetector(d Detector) *DetectorChain {
	dc.detectors = append(dc.detectors, d)
	return dc
}

//Test validates the additions against each detector in the chain.
//The results are passed in from detector to detector and thus collect all errors from all detectors
func (dc *DetectorChain) Test(additions []git_repo.Addition, ignores Ignores, result *DetectionResults) {
	for _, v := range dc.detectors {
		v.Test(additions, ignores, result)
	}
}
