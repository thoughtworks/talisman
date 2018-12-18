package detector

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"talisman/git_repo"
)

type PatternDetector struct {
	secretsPattern    *PatternMatcher
}

func (detector PatternDetector) Test(additions []git_repo.Addition, ignores Ignores, result *DetectionResults) {
	for _, addition := range additions {
		if ignores.Deny(addition, "filecontent") {
			log.WithFields(log.Fields{
				"filePath": addition.Path,
			}).Info("Ignoring addition as it was specified to be ignored.")
			result.Ignore(addition.Path, "filename")
			continue
		}
		detections := detector.secretsPattern.check(string(addition.Data))
		for _, detection := range detections {
			if detection != "" {
				log.WithFields(log.Fields{
					"filePath": addition.Path,
					"pattern":  detection,
				}).Info("Failing file as it matched pattern.")
				result.Fail(addition.Path, fmt.Sprintf("Potential secret pattern : %s", detection))
			}
		}
	}
}

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
