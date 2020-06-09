package pattern

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"talisman/talismanrc"
	"testing"
)

var (
	testRegexpPasswordPattern = `(?i)(['"_]?password['"]? *[:=][^,;]{8,})`
	testRegexpPassword        = regexp.MustCompile(testRegexpPasswordPattern)
	testRegexpPwPattern       = `(?i)(['"_]?pw['"]? *[:=][^,;]{8,})`
	testRegexpPw              = regexp.MustCompile(testRegexpPwPattern)
)

func TestShouldReturnEmptyStringWhenDoesNotMatchAnyRegex(t *testing.T) {
	assert.Equal(t, "", NewPatternMatcher([]*regexp.Regexp{testRegexpPassword}).check("safeString")[0])
}

func TestShouldReturnStringWhenMatchedPasswordPattern(t *testing.T) {
	assert.Equal(t, []string{"password\" :  123456789"}, NewPatternMatcher([]*regexp.Regexp{testRegexpPassword}).check("password\" :  123456789"))
	assert.Equal(t, []string{"pw\"  :  123456789"}, NewPatternMatcher([]*regexp.Regexp{testRegexpPw}).check("pw\"  :  123456789"))
}

func TestShouldAddGoodPatternToMatcher(t *testing.T) {
	pm := NewPatternMatcher([]*regexp.Regexp{})
	pm.add(talismanrc.PatternString(testRegexpPwPattern))
	assert.Equal(t, []string{"pw\"  :  123456789"}, NewPatternMatcher([]*regexp.Regexp{testRegexpPw}).check("pw\"  :  123456789"))
}


func TestShouldNotAddBadPatternToMatcher(t *testing.T) {
	pm := NewPatternMatcher([]*regexp.Regexp{})
	pm.add(`*a(crappy|regex`)
	assert.Equal(t, 0, len(pm.regexes))
}
