package talismanrc

import "regexp"

var emptyStringPattern = regexp.MustCompile(`^\s*$`)

func isEmptyString(str string) bool {
	return emptyStringPattern.MatchString(str)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
