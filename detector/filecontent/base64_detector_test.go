package filecontent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64DetectorShouldNotDetectSafeText(t *testing.T) {
	s := "pretty safe"
	bd := Base64Detector{}
	bd.initBase64Map()

	res := bd.CheckBase64Encoding(s)
	assert.Equal(t, "", res)
}

func TestBase64DetectorShouldDetectBase64Text(t *testing.T) {
	s := "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	bd := Base64Detector{}
	bd.initBase64Map()

	res := bd.CheckBase64Encoding(s)
	assert.Equal(t, s, res)
}

func TestBase64DetectorShouldNotDetectLongMethodNamesEvenWithHighEntropy(t *testing.T) {
	s := "TestBase64DetectorShouldNotDetectLongMethodNamesEvenWithRidiculousHighEntropyWordsMightExist"
	bd := Base64Detector{}
	bd.initBase64Map()

	res := bd.CheckBase64Encoding(s)
	assert.Equal(t, "", res)
}
