package detector

import (
	"github.com/stretchr/testify/assert"
	"github.com/thoughtworks/talisman/git_repo"
	"testing"
)

func TestShouldNotFlagSafeText(t *testing.T) {
	results := NewDetectionResults()
	content := []byte("prettySafe")
	filename := "filename"
	additions := []git_repo.Addition{git_repo.NewAddition(filename, content)}

	NewFileContentDetector().Test(additions, NewIgnores(), results)
	assert.False(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
}

func TestShouldFlagPotentialAWSAccessKeys(t *testing.T) {
	const awsAccessKeyIDExample string = "AKIAIOSFODNN7EXAMPLE"
	results := NewDetectionResults()
	content := []byte(awsAccessKeyIDExample)
	filename := "filename"
	additions := []git_repo.Addition{git_repo.NewAddition(filename, content)}

	NewFileContentDetector().Test(additions, NewIgnores(), results)
	assert.True(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
}

func TestShouldFlagPotentialAWSAccessKeysInPropertyDefinition(t *testing.T) {
	const awsAccessKeyIDExample string = "accessKey=AKIAIOSFODNN7EXAMPLE"
	results := NewDetectionResults()
	content := []byte(awsAccessKeyIDExample)
	filename := "filename"
	additions := []git_repo.Addition{git_repo.NewAddition(filename, content)}

	NewFileContentDetector().Test(additions, NewIgnores(), results)
	assert.True(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
}

func TestShouldNotFlag4CharSafeText(t *testing.T) {
	/*This only tell that an input could have been a b64 encoded value, but it does not tell whether or not the
	input is actually a b64 encoded value. In other words, abcd will match, but it is not necessarily represent
	 the encoded value of i· rather just a plain abcd input
	 see stackoverflow.com/questions/8571501/how-to-check-whether-the-string-is-base64-encoded-or-not#comment23919648_8571649*/
	results := NewDetectionResults()
	content := []byte("abcd")
	filename := "filename"
	additions := []git_repo.Addition{git_repo.NewAddition(filename, content)}

	NewFileContentDetector().Test(additions, NewIgnores(), results)
	assert.False(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
}

func TestShouldFlagPotentialAWSSecretKeys(t *testing.T) {
	const awsSecretAccessKey string = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	results := NewDetectionResults()
	content := []byte(awsSecretAccessKey)
	filename := "filename"
	additions := []git_repo.Addition{git_repo.NewAddition(filename, content)}

	NewFileContentDetector().Test(additions, NewIgnores(), results)
	assert.True(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")

}

func TestShouldFlagPotentialJWT(t *testing.T) {
	const jwt string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzY290Y2guaW8iLCJleHAiOjEzMDA4MTkzODAsIm5hbWUiOiJDaHJpcyBTZXZpbGxlamEiLCJhZG1pbiI6dHJ1ZX0.03f329983b86f7d9a9f5fef85305880101d5e302afafa20154d094b229f757"
	results := NewDetectionResults()
	content := []byte(jwt)
	filename := "filename"
	additions := []git_repo.Addition{git_repo.NewAddition(filename, content)}

	NewFileContentDetector().Test(additions, NewIgnores(), results)
	assert.True(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
}
