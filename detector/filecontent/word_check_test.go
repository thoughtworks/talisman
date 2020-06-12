package filecontent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordCheckWithWordsOnlyText(t *testing.T) {
	wc := WordCheck{}
	isWordsOnly := wc.containsWordsOnly("helloWorldGreetingsFromThoughtWorks")
	assert.True(t, isWordsOnly)
}

func TestWordCheckWithWordsOnlyLongText(t *testing.T) {
	wc := WordCheck{}
	isWordsOnly := wc.containsWordsOnly("TestBase64DetectorShouldNotDetectLongMethodNamesEvenWithRidiculousHighEntropyWordsMightExist")
	assert.True(t, isWordsOnly)
}

func TestWordCheckWithSingleWord(t *testing.T) {
	wc := WordCheck{}
	isWordsOnly := wc.containsWordsOnly("exception")
	assert.True(t, isWordsOnly)
}

func TestWordCheckWithSingleShortWord(t *testing.T) {
	wc := WordCheck{}
	isWordsOnly := wc.containsWordsOnly("for")
	assert.True(t, isWordsOnly)
}

func TestWordCheckWithSecret(t *testing.T) {
	wc := WordCheck{}
	isWordsOnly := wc.containsWordsOnly("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY")
	assert.False(t, isWordsOnly)
}

func TestWordCheckWithHex(t *testing.T) {
	wc := WordCheck{}
	isWordsOnly := wc.containsWordsOnly("68656C6C6F20776F726C6421")
	assert.False(t, isWordsOnly)
}

func TestWordCheckWithHalfWordsHalfHex(t *testing.T) {
	wc := WordCheck{}
	isWordsOnly := wc.containsWordsOnly("68656C6C6F20776F726C6421helloWorldGreetingsFromThoughtWorks")
	assert.False(t, isWordsOnly)
}

func TestWordCheckWithHalfWordsHalfSecret(t *testing.T) {
	wc := WordCheck{}
	isWordsOnly := wc.containsWordsOnly("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEYhelloWorldGreetingsFromThoughtWorks")
	assert.False(t, isWordsOnly)
}
