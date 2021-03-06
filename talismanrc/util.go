package talismanrc

import (
	"regexp"
	"github.com/spf13/afero"
)

var (
	emptyStringPattern = regexp.MustCompile(`^\s*$`)
	fs                 = afero.NewOsFs()
	currentRCFileName  = DefaultRCFileName
)

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
