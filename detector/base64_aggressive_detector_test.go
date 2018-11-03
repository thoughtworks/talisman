package detector

import (
	"testing"

	"talisman/git_repo"

	"github.com/stretchr/testify/assert"
)

func TestShouldFlagPotentialAWSAccessKeysInAggressiveMode(t *testing.T) {
	const awsAccessKeyIDExample string = "AKIAIOSFODNN7EXAMPLE\n"
	results := NewDetectionResults()
	content := []byte(awsAccessKeyIDExample)
	filename := "filename"
	additions := []git_repo.Addition{git_repo.NewAddition(filename, content)}

	NewFileContentDetector().AggressiveMode().Test(additions, NewIgnores(), results)
	assert.True(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
}

func TestShouldFlagPotentialAWSAccessKeysAtPropertyDefinitionInAggressiveMode(t *testing.T) {
	const awsAccessKeyIDExample string = "accessKey=AKIAIOSFODNN7EXAMPLE"
	results := NewDetectionResults()
	content := []byte(awsAccessKeyIDExample)
	filename := "filename"
	additions := []git_repo.Addition{git_repo.NewAddition(filename, content)}

	NewFileContentDetector().AggressiveMode().Test(additions, NewIgnores(), results)
	assert.True(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
}

func TestShouldNotFlagPotentialSecretsWithinSafeJavaCodeEvenInAggressiveMode(t *testing.T) {
	const awsAccessKeyIDExample string = "public class HelloWorld {\r\n\r\n    public static void main(String[] args) {\r\n        // Prints \"Hello, World\" to the terminal window.\r\n        System.out.println(\"Hello, World\");\r\n    }\r\n\r\n}"
	results := NewDetectionResults()
	content := []byte(awsAccessKeyIDExample)
	filename := "filename"
	additions := []git_repo.Addition{git_repo.NewAddition(filename, content)}

	NewFileContentDetector().AggressiveMode().Test(additions, NewIgnores(), results)
	if results == nil {
		additions = nil
	}
	assert.False(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
}
