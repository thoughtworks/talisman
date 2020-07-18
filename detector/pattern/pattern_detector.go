package pattern

import (
	"fmt"
	"regexp"
	"sync"
	"talisman/detector/helpers"
	"talisman/gitrepo"
	"talisman/talismanrc"

	log "github.com/Sirupsen/logrus"
)

type PatternDetector struct {
	secretsPattern *PatternMatcher
}

var (
	detectorPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)((.*)(password|passphrase|secret|key|pwd|pword|pass)(.*) *[:=>][^,;\n]{8,})`),
		regexp.MustCompile(`(?i)(['"_]?pw['"]? *[:=][^,;\n]{8,})`),
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
func (detector PatternDetector) Test(comparator helpers.ChecksumCompare, currentAdditions []gitrepo.Addition, ignoreConfig *talismanrc.TalismanRC, result *helpers.DetectionResults) {
	matches := make(chan match, 512)
	ignoredFilePaths := make(chan gitrepo.FilePath, 512)
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(len(currentAdditions))
	for _, addition := range currentAdditions {
		go func(addition gitrepo.Addition) {
			defer waitGroup.Done()
			if ignoreConfig.Deny(addition, "filecontent") || comparator.IsScanNotRequired(addition) {
				ignoredFilePaths <- addition.Path
				return
			}
			detections := detector.secretsPattern.check(processAllowedPatterns(addition, ignoreConfig))
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

func replaceAllStrings(str string, pattern string) string {
	re := regexp.MustCompile(fmt.Sprintf("(?i)%s", pattern))
	return re.ReplaceAllString(str, "TALISMAN_ALLOWED_PATTERN")
}

func processAllowedPatterns(addition gitrepo.Addition, tRC *talismanrc.TalismanRC) string {
	// Processing global allowed patterns
	for _, pattern := range tRC.AllowedPatterns {
		addition.Data = []byte(replaceAllStrings(string(addition.Data), pattern))
	}

	// Processing allowed patterns based on file path
	for k, v := range tRC.FileIgnoreConfig {
		if v.FileName == string(addition.Path) {
			for _, pattern := range v.AllowedPatterns {
				addition.Data = []byte(replaceAllStrings(string(addition.Data), pattern))
			}
		}
	}
	return string(addition.Data)
}

func (detector PatternDetector) processIgnore(ignoredFilePath gitrepo.FilePath, result *helpers.DetectionResults) {
	log.WithFields(log.Fields{
		"filePath": ignoredFilePath,
	}).Info("Ignoring addition as it was specified to be ignored.")
	result.Ignore(ignoredFilePath, "filecontent")
}

func (detector PatternDetector) processMatch(match match, result *helpers.DetectionResults) {
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
