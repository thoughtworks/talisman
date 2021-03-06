package pattern

import (
	"fmt"
	"regexp"
	"sync"
	"talisman/detector/helpers"
	"talisman/detector/severity"
	"talisman/gitrepo"
	"talisman/talismanrc"

	log "github.com/Sirupsen/logrus"
)

type PatternDetector struct {
	secretsPattern *PatternMatcher
}

var (
	detectorPatterns = []*severity.PatternSeverity{
		{Pattern: regexp.MustCompile(`(?i)((.*)(password|passphrase|secret|key|pwd|pword|pass)(.*) *[:=>][^,;\n]{8,})`), Severity: severity.SeverityConfiguration["PasswordPhrasePattern"]},
		{Pattern: regexp.MustCompile(`(?i)((:)(password|passphrase|secret|key|pwd|pword|pass)(.*) *[ ][^,;\n]{8,})`), Severity: severity.SeverityConfiguration["PasswordPhrasePattern"]},
		{Pattern: regexp.MustCompile(`(?i)(['"_]?pw['"]? *[:=][^,;\n]{8,})`), Severity: severity.SeverityConfiguration["PasswordPhrasePattern"]},
		{Pattern: regexp.MustCompile(`(?i)(<ConsumerKey>\S*</ConsumerKey>)`), Severity: severity.SeverityConfiguration["ConsumerKeyPattern"]},
		{Pattern: regexp.MustCompile(`(?i)(<ConsumerSecret>\S*</ConsumerSecret>)`), Severity: severity.SeverityConfiguration["ConsumerSecretParrern"]},
		{Pattern: regexp.MustCompile(`(?i)(AWS[ \w]+key[ \w]+[:=])`), Severity: severity.SeverityConfiguration["AWSKeyPattern"]},
		{Pattern: regexp.MustCompile(`(?i)(AWS[ \w]+secret[ \w]+[:=])`), Severity: severity.SeverityConfiguration["AWSSecretPattern"]},
		{Pattern: regexp.MustCompile(`(?s)(BEGIN RSA PRIVATE KEY.*END RSA PRIVATE KEY)`), Severity: severity.SeverityConfiguration["RSAKeyPattern"]},
	}
)

type match struct {
	name       gitrepo.FileName
	path       gitrepo.FilePath
	commits    []string
	detections []DetectionsWithSeverity
}

//Test tests the contents of the Additions to ensure that they don't look suspicious
func (detector PatternDetector) Test(comparator helpers.ChecksumCompare, currentAdditions []gitrepo.Addition, ignoreConfig *talismanrc.TalismanRC, result *helpers.DetectionResults, additionCompletionCallback func()) {
	matches := make(chan match, 512)
	ignoredFilePaths := make(chan gitrepo.FilePath, 512)
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(len(currentAdditions))
	for _, addition := range currentAdditions {
		go func(addition gitrepo.Addition) {
			defer waitGroup.Done()
			defer additionCompletionCallback()
			if ignoreConfig.Deny(addition, "filecontent") || comparator.IsScanNotRequired(addition) {
				ignoredFilePaths <- addition.Path
				return
			}
			detections := detector.secretsPattern.check(processAllowedPatterns(addition, ignoreConfig), ignoreConfig.Threshold)
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
			detector.processMatch(match, result, ignoreConfig.Threshold)
		case ignore, hasMore := <-ignoredFilePaths:
			if !hasMore {
				ignoredChanHasMore = false
				continue
			}
			detector.processIgnore(ignore, result)
		}
	}
}

var knownPatterns = make(map[string]*regexp.Regexp)

func replaceAllStrings(str string, pattern string) string {
	var pat *regexp.Regexp
	var ok bool

	if pat, ok = knownPatterns[pattern]; !ok {
		pat = regexp.MustCompile(fmt.Sprintf("(?i)%s", pattern))
		knownPatterns[pattern] = pat
	}

	return pat.ReplaceAllString(str, "")
}

func processAllowedPatterns(addition gitrepo.Addition, tRC *talismanrc.TalismanRC) string {
	// Processing global allowed patterns
	for _, pattern := range tRC.AllowedPatterns {
		addition.Data = []byte(replaceAllStrings(string(addition.Data), pattern))
	}

	// Processing allowed patterns based on file path
	for _, fileignoreconfig := range tRC.FileIgnoreConfig {
		if fileignoreconfig.FileName == string(addition.Path) {
			for _, pattern := range fileignoreconfig.AllowedPatterns {
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

func (detector PatternDetector) processMatch(match match, result *helpers.DetectionResults, threshold severity.Severity) {
	for _, detectionWithSeverity := range match.detections {
		for _, detection := range detectionWithSeverity.detections {
			if detection != "" {
				if string(match.name) == talismanrc.DefaultRCFileName || !detectionWithSeverity.severity.ExceedsThreshold(threshold) {
					log.WithFields(log.Fields{
						"filePath": match.path,
						"pattern":  detection,
					}).Warn("Warning file as it matched pattern.")
					result.Warn(match.path, "filecontent", fmt.Sprintf("Potential secret pattern : %s", detection), match.commits, detectionWithSeverity.severity)
				} else {
					log.WithFields(log.Fields{
						"filePath": match.path,
						"pattern":  detection,
					}).Info("Failing file as it matched pattern.")
					result.Fail(match.path, "filecontent", fmt.Sprintf("Potential secret pattern : %s", detection), match.commits, detectionWithSeverity.severity)
				}
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
