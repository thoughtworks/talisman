package detector

import "regexp"

type PatternMatcher struct {
	regexes []*regexp.Regexp
}

func (detector PatternMatcher) check(content string) []string {
	var detected []string
	for _, regex := range detector.regexes {
		matches := regex.FindStringSubmatch(content)
		if matches != nil {
			detected = append(detected, matches[1:]...)
		}
	}
	if detected != nil {
		return detected
	}
	return []string{""}
}

func NewSecretsPatternDetector(patterns []*regexp.Regexp) *PatternMatcher {
	return &PatternMatcher{patterns}
}
