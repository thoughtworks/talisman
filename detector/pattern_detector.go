package detector

import (
	"fmt"
	"regexp"
	"sync"
	"talisman/gitrepo"
	"talisman/talismanrc"

	log "github.com/Sirupsen/logrus"
)

type PatternDetector struct {
	secretsPattern *PatternMatcher
}

var (
	detectorPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(['"_]?password['"]? *[:=][^,;\n]{8,})`),
		regexp.MustCompile(`(?i)(['"_]?pw['"]? *[:=][^,;\n]{8,})`),
		regexp.MustCompile(`(?i)(['"_]?pwd['"]? *[:=][^,;\n]{8,})`),
		regexp.MustCompile(`(?i)(['"_]?pass['"]? *[:=][^,;\n]{8,})`),
		regexp.MustCompile(`(?i)(['"_]?pword['"]? *[:=][^,;\n]{8,})`),
		regexp.MustCompile(`(?i)(['"_]?adminPassword['"]? *[:=\n][^,;]{8,})`),
		regexp.MustCompile(`(?i)(['"_]?passphrase['"]? *[:=\n][^,;]{8,})`),
		regexp.MustCompile(`(<[^(><.)]?password[^(><.)]*?>[^(><.)]+</[^(><.)]?password[^(><.)]*?>)`),
		regexp.MustCompile(`(<[^(><.)]?passphrase[^(><.)]*?>[^(><.)]+</[^(><.)]?passphrase[^(><.)]*?>)`),
		regexp.MustCompile(`(?i)(<ConsumerKey>\S*</ConsumerKey>)`),
		regexp.MustCompile(`(?i)(<ConsumerSecret>\S*</ConsumerSecret>)`),
		regexp.MustCompile(`(?i)(AWS[ \w]+key[ \w]+[:=])`),
		regexp.MustCompile(`(?i)(AWS[ \w]+secret[ \w]+[:=])`),
		regexp.MustCompile(`(?s)(BEGIN RSA PRIVATE KEY.*END RSA PRIVATE KEY)`),
	}
)

type match struct {
	name       gitrepo.FileName
	path       gitrepo.FilePath
	commits    []string
	detections []string
}

//Test tests the contents of the Additions to ensure that they don't look suspicious
func (detector PatternDetector) Test(additions []gitrepo.Addition, ignoreConfig *talismanrc.TalismanRC, result *DetectionResults) {
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
	for ignoredChanHasMore, matchChanHasMore := true, true; ignoredChanHasMore || matchChanHasMore; {
		select {
		case match, hasMore := <-matches:
			if !hasMore {
				matchChanHasMore = false
				continue
			}
			detector.processMatch(match, result)
		case ignore, hasMore := <-ignoredFilePaths:
			if !hasMore {
				ignoredChanHasMore = false
				continue
			}
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
			if string(match.name) == talismanrc.DefaultRCFileName {
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
func NewPatternDetector(custom []talismanrc.PatternString) *PatternDetector {
	matcher := NewPatternMatcher(detectorPatterns)
	for _, pattern := range custom {
		matcher.add(pattern)
	}
	return &PatternDetector{matcher}
}
