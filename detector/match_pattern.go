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

func initPattern(patternStrings []string) *PatternMatcher {
	var patterns = make([]*regexp.Regexp, len(patternStrings))
	for i, pattern := range patternStrings {
		patterns[i], _ = regexp.Compile(pattern)
	}
	return &PatternMatcher{patterns}
}

func NewSecretsPatternDetector(patterns []string) *PatternMatcher {
	return initPattern(patterns)
}
