package pattern

import (
	"regexp"
	"talisman/detector/severity"
	"talisman/talismanrc"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testRegexpPasswordPattern = `(?i)(['"_]?password['"]? *[:=][^,;]{8,})`
	testRegexpPassword        = regexp.MustCompile(testRegexpPasswordPattern)
	testRegexpPwPattern       = `(?i)(['"_]?pw['"]? *[:=][^,;]{8,})`
	testRegexpPw              = regexp.MustCompile(testRegexpPwPattern)
)

func TestShouldReturnEmptyStringWhenDoesNotMatchAnyRegex(t *testing.T) {
	detections := NewPatternMatcher([]*severity.PatternSeverity{{Pattern: testRegexpPassword, Severity: severity.Low()}}).check("safeString", severity.LowSeverity)
	assert.Equal(t, []DetectionsWithSeverity(nil), detections)
}

func TestShouldReturnStringWhenMatchedPasswordPattern(t *testing.T) {
	detections1 := NewPatternMatcher([]*severity.PatternSeverity{{Pattern: testRegexpPassword, Severity: severity.Low()}}).check("password\" :  123456789", severity.LowSeverity)
	detections2 := NewPatternMatcher([]*severity.PatternSeverity{{Pattern: testRegexpPw, Severity: severity.Medium()}}).check("pw\"  :  123456789", severity.LowSeverity)
	assert.Equal(t, []DetectionsWithSeverity{{detections: []string{"password\" :  123456789"}, severity: severity.Low()}}, detections1)
	assert.Equal(t, []DetectionsWithSeverity{{detections: []string{"pw\"  :  123456789"}, severity: severity.Medium()}}, detections2)
}

func TestShouldAddGoodPatternWithHighSeverityToMatcher(t *testing.T) {
	pm := NewPatternMatcher([]*severity.PatternSeverity{})
	pm.add(talismanrc.PatternString(testRegexpPwPattern))
	detections := pm.check("pw\"  :  123456789", severity.LowSeverity)
	assert.Equal(t, []DetectionsWithSeverity{{detections: []string{"pw\"  :  123456789"}, severity: severity.High()}}, detections)
}

func TestShouldNotAddBadPatternToMatcher(t *testing.T) {
	pm := NewPatternMatcher([]*severity.PatternSeverity{})
	pm.add(`*a(crappy|regex`)
	assert.Equal(t, 0, len(pm.regexes))
}
