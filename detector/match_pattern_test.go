package detector

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

var (
	testRegexpPassword = regexp.MustCompile(`(?i)(['"_]?password['"]? *[:=][^,;]{8,})`)
	testRegexpPw       = regexp.MustCompile(`(?i)(['"_]?pw['"]? *[:=][^,;]{8,})`)
)

func TestShouldReturnEmptyStringWhenDoesNotMatchAnyRegex(t *testing.T) {
	assert.Equal(t, "", NewSecretsPatternDetector([]*regexp.Regexp{testRegexpPassword}).check("safeString")[0])
}

func TestShouldReturnStringWhenMatchedPasswordPattern(t *testing.T) {
	assert.Equal(t, []string{"password\" :  123456789"}, NewSecretsPatternDetector([]*regexp.Regexp{testRegexpPassword}).check("password\" :  123456789"))
	assert.Equal(t, []string{"pw\"  :  123456789"}, NewSecretsPatternDetector([]*regexp.Regexp{testRegexpPw}).check("pw\"  :  123456789"))
}
