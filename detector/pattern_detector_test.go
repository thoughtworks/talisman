package detector

import (
	"github.com/stretchr/testify/assert"
	"talisman/git_repo"
	"testing"
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
	ignores := NewIgnores(filename)

	NewPatternDetector().Test(additions, ignores, results)
	assert.True(t, results.Successful(), "Expected file %s to be ignored by pattern", filename)
}

func shouldPassDetectionOfSecretPattern(filename string, content []byte, t *testing.T) {
	results := NewDetectionResults()
	additions := []git_repo.Addition{git_repo.NewAddition(filename, content)}
	NewPatternDetector().Test(additions, NewIgnores(), results)
	assert.Equal(t, "Potential secret pattern : "+string(content), results.Failures(additions[0].Path)[0])
}