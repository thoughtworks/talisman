package pattern

import (
	"fmt"
	"regexp"
	"talisman/detector/severity"
	"talisman/talismanrc"

	"github.com/Sirupsen/logrus"
)

type PatternMatcher struct {
	regexes []*severity.PatternSeverity
}

type DetectionsWithSeverity struct {
	detections []string
	severity   severity.Severity
}

func (pm *PatternMatcher) check(content string, thresholdValue severity.Severity) []DetectionsWithSeverity {
	var detectionsWithSeverity []DetectionsWithSeverity
	for _, pattern := range pm.regexes {
		var detected []string
		regex := pattern.Pattern
		logrus.Debugf("checking for pattern %v", regex)
		matches := regex.FindAllString(content, -1)
		if matches != nil {
			detected = append(detected, matches...)
			detectionsWithSeverity = append(detectionsWithSeverity, DetectionsWithSeverity{detections: detected, severity: pattern.Severity})
		}
	}
	return detectionsWithSeverity
}

func (pm *PatternMatcher) add(ps talismanrc.PatternString) {
	re, err := regexp.Compile(fmt.Sprintf("(%s)", string(ps)))
	if err != nil {
		logrus.Warnf("ignoring invalid pattern '%s'", ps)
		return
	}
	logrus.Infof("added custom pattern '%s' with high severity", ps)
	pm.regexes = append(pm.regexes, &severity.PatternSeverity{Pattern: re, Severity: severity.SeverityConfiguration["CustomPattern"]})
}

func NewPatternMatcher(patterns []*severity.PatternSeverity) *PatternMatcher {
	return &PatternMatcher{patterns}
}
