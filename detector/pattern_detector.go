package detector

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"sync"
	"talisman/gitrepo"
)

type PatternDetector struct {
	secretsPattern *PatternMatcher
}

type match struct {
	name       gitrepo.FileName
	path       gitrepo.FilePath
	commits    []string
	detections []string
}

//Test tests the contents of the Additions to ensure that they don't look suspicious
func (detector PatternDetector) Test(additions []gitrepo.Addition, ignoreConfig TalismanRCIgnore, result *DetectionResults) {
	cc := NewChecksumCompare(additions, ignoreConfig)
	matches := make(chan match, 512)
	ignoredFilePaths := make(chan gitrepo.FilePath, 512)
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(len(additions))
	for _, addition := range additions {
		go func(addition gitrepo.Addition) {
			defer waitGroup.Done()
			if ignoreConfig.Deny(addition, "filecontent") || cc.IsScanNotRequired(addition) {
				ignoredFilePaths <- addition.Path
				return
			}
			detections := detector.secretsPattern.check(string(addition.Data))
			matches <- match{name: addition.Name, path: addition.Path, detections: detections, commits: addition.Commits}
		}(addition)
	}
	go func() {
		waitGroup.Wait()
		close(matches)
		close(ignoredFilePaths)
	}()
	for i := 0; i < len(additions); i++ {
		select {
		case match := <-matches:
			detector.processMatch(match, result)
		case ignore := <-ignoredFilePaths:
			detector.processIgnore(ignore, result)
		}
	}
}

func (detector PatternDetector) processIgnore(ignoredFilePath gitrepo.FilePath, result *DetectionResults) {
	log.WithFields(log.Fields{
		"filePath": ignoredFilePath,
	}).Info("Ignoring addition as it was specified to be ignored.")
	result.Ignore(ignoredFilePath, "filecontent")
}

func (detector PatternDetector) processMatch(match match, result *DetectionResults) {
	for _, detection := range match.detections {
		if detection != "" {
			if string(match.name) == DefaultRCFileName {
				log.WithFields(log.Fields{
					"filePath": match.path,
					"pattern":  detection,
				}).Warn("Warning file as it matched pattern.")
				result.Warn(match.path, "filecontent", fmt.Sprintf("Potential secret pattern : %s", detection), match.commits)
			} else {
				log.WithFields(log.Fields{
					"filePath": match.path,
					"pattern":  detection,
				}).Info("Failing file as it matched pattern.")
				result.Fail(match.path, "filecontent", fmt.Sprintf("Potential secret pattern : %s", detection), match.commits)
			}
		}
	}
}

//NewPatternDetector returns a PatternDetector that tests Additions against the pre-configured patterns
func NewPatternDetector() *PatternDetector {
	patternStrings := []string{
		"(?i)(['|\"|_]?password['|\"]? *[:|=][^,|;|\n]{8,})",
		"(?i)(['|\"|_]?pw['|\"]? *[:|=][^,|;|\n]{8,})",
		"(?i)(['|\"|_]?pwd['|\"]? *[:|=][^,|;|\n]{8,})",
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
		"(?s)(BEGIN RSA PRIVATE KEY.*END RSA PRIVATE KEY)",
	}

	return &PatternDetector{NewSecretsPatternDetector(patternStrings)}
}
