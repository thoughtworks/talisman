package detector

import (
	"talisman/git_repo"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldDetectPasswordPatterns(t *testing.T) {
	filename := "secret.txt"

	shouldPassDetectionOfSecretPattern(filename, []byte("\"password\" : UnsafePassword"), t)
	shouldPassDetectionOfSecretPattern(filename, []byte("<password data=123> jdghfakjkdha</password>"), t)
	shouldPassDetectionOfSecretPattern(filename, []byte("<passphrase data=123> AasdfYlLKHKLasdKHAFKHSKmlahsdfLK</passphrase>"), t)
	shouldPassDetectionOfSecretPattern(filename, []byte("<ConsumerKey>alksjdhfkjaklsdhflk12345adskjf</ConsumerKey>"), t)
	shouldPassDetectionOfSecretPattern(filename, []byte("AWS key :"), t)
	shouldPassDetectionOfSecretPattern(filename, []byte(`BEGIN RSA PRIVATE KEY-----
aghjdjadslgjagsfjlsgjalsgjaghjldasja
-----END RSA PRIVATE KEY`), t)
}

func TestShouldIgnorePasswordPatterns(t *testing.T) {
	results := NewDetectionResults()
	content := []byte("\"password\" : UnsafePassword")
	filename := "secret.txt"
	additions := []git_repo.Addition{git_repo.NewAddition(filename, content)}
	fileIgnoreConfig := FileIgnoreConfig{filename, "833b6c24c8c2c5c7e1663226dc401b29c005492dc76a1150fc0e0f07f29d4cc3", []string{"filecontent"}}
	ignores := TalismanRCIgnore{[]FileIgnoreConfig{fileIgnoreConfig}}

	NewPatternDetector().Test(additions, ignores, results)
	assert.True(t, results.Successful(), "Expected file %s to be ignored by pattern", filename)
}

func shouldPassDetectionOfSecretPattern(filename string, content []byte, t *testing.T) {
	results := NewDetectionResults()
	additions := []git_repo.Addition{git_repo.NewAddition(filename, content)}
	NewPatternDetector().Test(additions, TalismanRCIgnore{}, results)
	expected := "Potential secret pattern : " + string(content)
	assert.Equal(t, expected, getFailureMessage(results, additions))
}

func getFailureMessage(results *DetectionResults, additions []git_repo.Addition) string {
	failureMessages := []string{}
	for failureMessage := range results.GetFailures(additions[0].Path).FailuresInCommits {
		failureMessages = append(failureMessages, failureMessage)
	}
	return failureMessages[0]
}
