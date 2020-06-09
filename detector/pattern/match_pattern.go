package pattern

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"regexp"
	"talisman/talismanrc"
)

type PatternMatcher struct {
	regexes []*regexp.Regexp
}

func (pm *PatternMatcher) check(content string) []string {
	var detected []string
	for _, regex := range pm.regexes {
		logrus.Debugf("checking for pattern %v", regex)
		matches := regex.FindAllString(content, -1)
		if matches != nil {
			detected = append(detected, matches...)
		}
	}
	if detected != nil {
		return detected
	}
	return []string{""}
}

func (pm *PatternMatcher) add(ps talismanrc.PatternString) {
	re, err := regexp.Compile(fmt.Sprintf("(%s)", string(ps)))
	if err != nil {
		logrus.Warnf("ignoring invalid pattern '%s'", ps)
		return
	}
	logrus.Infof("added custom pattern '%s'", ps)
	pm.regexes = append(pm.regexes, re)
}

func NewPatternMatcher(patterns []*regexp.Regexp) *PatternMatcher {
	return &PatternMatcher{patterns}
}
