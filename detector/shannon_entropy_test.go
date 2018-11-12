package detector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntropyCandidatesShouldBeFoundForGivenSet(t *testing.T) {
	entropy := Entropy{}
	dc := Base64Detector{}
	dc.initBase64Map()
	candidatesWithinWord := entropy.GetEntropyCandidatesWithinWord("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", 20, dc.base64Map)
	assert.Equal(t, 1, len(candidatesWithinWord))
}

func TestEntropyCandidatesShouldBeEmptyForShorterWords(t *testing.T) {
	entropy := Entropy{}
	dc := Base64Detector{}
	dc.initBase64Map()
	candidatesWithinWord := entropy.GetEntropyCandidatesWithinWord("abc", 4, dc.base64Map)
	assert.Equal(t, 0, len(candidatesWithinWord))
}

func TestEntropyValueOfSecretShouldBeHigherThanFour(t *testing.T) {
	entropy := Entropy{}
	shannonEntropy := entropy.GetShannonEntropy("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", BASE64_CHARS)
	assert.True(t, 4 < shannonEntropy)
}

func TestEntropyValueOfEmptyStringShouldBeZero(t *testing.T) {
	entropy := Entropy{}
	shannonEntropy := entropy.GetShannonEntropy("", BASE64_CHARS)
	assert.Equal(t, float64(0), shannonEntropy)
}
