package detector

import (
	"strings"
	"talisman/gitrepo"
	"talisman/talismanrc"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	customPatterns []talismanrc.PatternString
)

func TestShouldDetectPasswordPatterns(t *testing.T) {
	filename := "secret.txt"
	values := [7]string {"password","secret", "key", "pwd","pass","pword","passphrase"}
	for i := 0; i < len(values); i++ {
		shouldPassDetectionOfSecretPattern(filename, []byte(strings.ToTitle(values[i])+":UnsafeString"), t)
		shouldPassDetectionOfSecretPattern(filename, []byte(values[i]+"=UnsafeString"), t)
		shouldPassDetectionOfSecretPattern(filename, []byte("\"SERVER_"+strings.ToUpper(values[i])+"\" : UnsafeString"), t)
		shouldPassDetectionOfSecretPattern(filename, []byte(values[i]+"2-string : UnsafeString"), t)
		shouldPassDetectionOfSecretPattern(filename, []byte("<"+values[i]+" data=123> randomStringGoesHere </"+values[i]+">"), t)
		shouldPassDetectionOfSecretPattern(filename, []byte("<admin "+values[i]+"> randomStringGoesHere </my"+values[i]+">"), t)
	}

	shouldPassDetectionOfSecretPattern(filename, []byte("\"pw\" : UnsafeString"), t)
	shouldPassDetectionOfSecretPattern(filename, []byte("Pw=UnsafeString"), t)

	shouldPassDetectionOfSecretPattern(filename, []byte("<ConsumerKey>alksjdhfkjaklsdhflk12345adskjf</ConsumerKey>"), t)
	shouldPassDetectionOfSecretPattern(filename, []byte("AWS key :"), t)
	shouldPassDetectionOfSecretPattern(filename, []byte(`BEGIN RSA PRIVATE KEY-----
	aghjdjadslgjagsfjlsgjalsgjaghjldasja
	-----END RSA PRIVATE KEY`), t)
	shouldPassDetectionOfSecretPattern(filename, []byte(`PWD=appropriate`), t)
	shouldPassDetectionOfSecretPattern(filename, []byte(`pass=appropriate`), t)
	shouldPassDetectionOfSecretPattern(filename, []byte(`adminpwd=appropriate`), t)

	shouldFailDetectionOfSecretPattern(filename, []byte("\"pAsSWoRD\" :1234567"), t)
	shouldFailDetectionOfSecretPattern(filename, []byte(`setPassword("12345678")`), t)
	shouldFailDetectionOfSecretPattern(filename, []byte(`setenv(password, "12345678")`), t)
	shouldFailDetectionOfSecretPattern(filename, []byte(`random=12345678)`), t)
}

func TestShouldIgnorePasswordPatterns(t *testing.T) {
	results := NewDetectionResults()
	content := []byte("\"password\" : UnsafePassword")
	filename := "secret.txt"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}
	fileIgnoreConfig := talismanrc.FileIgnoreConfig{filename, "833b6c24c8c2c5c7e1663226dc401b29c005492dc76a1150fc0e0f07f29d4cc3", []string{"filecontent"}}
	ignores := &talismanrc.TalismanRC{FileIgnoreConfig: []talismanrc.FileIgnoreConfig{fileIgnoreConfig}}

	NewPatternDetector(customPatterns).Test(ChecksumCompare{calculator: nil, talismanRC: talismanrc.NewTalismanRC(nil)}, additions, ignores, results)
	assert.True(t, results.Successful(), "Expected file %s to be ignored by pattern", filename)
}

func DetectionOfSecretPattern(filename string, content []byte) (*DetectionResults, []gitrepo.Addition, string) {
	results := NewDetectionResults()
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}
	NewPatternDetector(customPatterns).Test(ChecksumCompare{calculator: nil, talismanRC: talismanrc.NewTalismanRC(nil)}, additions, talismanRC, results)
	expected := "Potential secret pattern : " + string(content)
	return results, additions, expected
}

func shouldPassDetectionOfSecretPattern(filename string, content []byte, t *testing.T) {
	results, additions, expected := DetectionOfSecretPattern(filename, content)
	assert.Equal(t, expected, getFailureMessage(results, additions))
	assert.Len(t, results.Results, 1)
}

func shouldFailDetectionOfSecretPattern(filename string, content []byte, t *testing.T) {
	results, additions, expected := DetectionOfSecretPattern(filename, content)
	assert.NotEqual(t, expected, getFailureMessage(results, additions))
	assert.Len(t, results.Results, 0)
}

func getFailureMessage(results *DetectionResults, additions []gitrepo.Addition) string {
	failureMessages := []string{}
	for _, failureDetails := range results.GetFailures(additions[0].Path) {
		failureMessages = append(failureMessages, failureDetails.Message)
	}
	if len(failureMessages) == 0 {
		return ""
	}
	return failureMessages[0]
}
