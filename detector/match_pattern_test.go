package detector

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldReturnEmptyStringWhenDoesNotMatchAnyRegex(t *testing.T) {
	assert.Equal(t, "", NewSecretsPatternDetector([]string{"(?i)(['|\"|_]?password['|\"]? *[:|=][^,|;]{8,})"}).check("safeString")[0])
}

func TestShouldReturnStringWhenMatchedPasswordPattern(t *testing.T) {
	assert.Equal(t, []string{"password\" :  123456789"}, NewSecretsPatternDetector([]string{"(?i)(['|\"|_]?password['|\"]? *[:|=][^,|;]{8,})"}).check("password\" :  123456789"))
	assert.Equal(t, []string{"pw\"  :  123456789"}, NewSecretsPatternDetector([]string{"(?i)(['|\"|_]?pw['|\"]? *[:|=][^,|;]{8,})"}).check("pw\"  :  123456789"))
}
