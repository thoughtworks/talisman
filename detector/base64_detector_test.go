package detector

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBase64DetectorShouldNotDetectSafeText(t *testing.T) {
	s := "pretty safe"
	bd := Base64Detector{}
	bd.initBase64Map()

	res := bd.checkBase64Encoding(s)
	assert.Equal(t, "", res)
}

func TestBase64DetectorShouldDetectBase64Text(t *testing.T) {
	s := "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	bd := Base64Detector{}
	bd.initBase64Map()

	res := bd.checkBase64Encoding(s)
	assert.Equal(t, s, res)
}