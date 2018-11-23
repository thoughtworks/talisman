package detector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHexDetectorShouldNotDetectSafeText(t *testing.T) {
	s := "pretty safe"
	hd := HexDetector{}
	hd.initHexMap()

	res := hd.checkHexEncoding(s)
	assert.Equal(t, "", res)
}

func TestHexDetectorShouldDetectBase64Text(t *testing.T) {
	s := "6A6176617375636B73676F726F636B7368616861"
	hd := HexDetector{}
	hd.initHexMap()

	res := hd.checkHexEncoding(s)
	assert.Equal(t, s, res)
}
