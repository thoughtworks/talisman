package detector

import (
	"fmt"
	"talisman/git_repo"

	log "github.com/Sirupsen/logrus"
)

type PatternDetector struct {
	secretsPattern *PatternMatcher
}

//Test tests the contents of the Additions to ensure that they don't look suspicious
func (detector PatternDetector) Test(additions []git_repo.Addition, ignoreConfig TalismanRCIgnore, result *DetectionResults) {
	cc := NewChecksumCompare(additions, ignoreConfig)
	for _, addition := range additions {
		if ignoreConfig.Deny(addition, "filecontent") || cc.IsScanNotRequired(addition) {
			log.WithFields(log.Fields{
				"filePath": addition.Path,
			}).Info("Ignoring addition as it was specified to be ignored.")
			result.Ignore(addition.Path, "filecontent")
			continue
		}
		detections := detector.secretsPattern.check(string(addition.Data))
		for _, detection := range detections {
			if detection != "" {
				if string(addition.Name) == DefaultRCFileName {
					log.WithFields(log.Fields{
						"filePath": addition.Path,
						"pattern":  detection,
					}).Warn("Warning file as it matched pattern.")
					result.Warn(addition.Path, fmt.Sprintf("Potential secret pattern : %s", detection), addition.Commits)
				} else {
					log.WithFields(log.Fields{
						"filePath": addition.Path,
						"pattern":  detection,
					}).Info("Failing file as it matched pattern.")
					result.Fail(addition.Path, fmt.Sprintf("Potential secret pattern : %s", detection), addition.Commits)
				}
			}
		}
	}
}

//NewPatternDetector returns a PatternDetector that tests Additions against the pre-configured patterns
func NewPatternDetector() *PatternDetector {
	patternStrings := []string{

		"(?i)(['|\"|_]?password['|\"]? *[:|=][^,|;|\n]{8,})",
		"(?i)(['|\"|_]?pw['|\"]? *[:|=][^,|;|\n]{8,})",
		"(?i)(['|\"|_]?pass['|\"]? *[:|=][^,|;|\n]{8,})",
		"(?i)(['|\"|_]?pword['|\"]? *[:|=][^,|;|\n]{8,})",
		"(?i)(['|\"|_]?adminPassword['|\"]? *[:|=|\n][^,|;]{8,})",
		"(?i)(['|\"|_]?passphrase['|\"]? *[:|=|\n][^,|;]{8,})",
		"(<[^(><.)]?password[^(><.)]*?>[^(><.)]+</[^(><.)]?password[^(><.)]*?>)",
		"(<[^(><.)]?passphrase[^(><.)]*?>[^(><.)]+</[^(><.)]?passphrase[^(><.)]*?>)",
		"(?i)(<ConsumerKey>\\S*<\\/ConsumerKey>)",
		"(?i)(<ConsumerSecret>\\S*<\\/ConsumerSecret>)",
		"(?i)(AWS[ |\\w]+key[ |\\w]+[:|=])",
		"(?i)(AWS[ |\\w]+secret[ |\\w]+[:|=])",
		"(?s)(BEGIN RSA PRIVATE KEY.*END RSA PRIVATE KEY)"}

	return &PatternDetector{NewSecretsPatternDetector(patternStrings)}
}
