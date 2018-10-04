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

func TestShouldIgnoreFileIfNeeded(t *testing.T) {
	results := NewDetectionResults()
	content := []byte("prettySafe")
	filename := "filename"
	additions := []git_repo.Addition{git_repo.NewAddition(filename, content)}
	ignores := NewIgnores(filename)

	NewFileContentDetector().Test(additions, ignores, results)
	assert.True(t, results.Successful(), "Expected file %s to be ignored by pattern", filename)
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

func TestShouldNotFlagLowEntropyBase64Text(t *testing.T) {
	const lowEntropyString string = "YWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWEK"
	results := NewDetectionResults()
	content := []byte(lowEntropyString)
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

func TestShouldFlagPotentialSecretsWithinJavaCode(t *testing.T) {
	const dangerousJavaCode string = "public class HelloWorld {\r\n\r\n    public static void main(String[] args) {\r\n        // Prints \"Hello, World\" to the terminal window.\r\n        accessKey=\"wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY\";\r\n        System.out.println(\"Hello, World\");\r\n    }\r\n\r\n}"
	results := NewDetectionResults()
	content := []byte(dangerousJavaCode)
	filename := "filename"
	additions := []git_repo.Addition{git_repo.NewAddition(filename, content)}

	NewFileContentDetector().Test(additions, NewIgnores(), results)
	assert.True(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
}

func TestShouldNotFlagPotentialSecretsWithinSafeJavaCode(t *testing.T) {
	const safeJavaCode string = "public class HelloWorld {\r\n\r\n    public static void main(String[] args) {\r\n        // Prints \"Hello, World\" to the terminal window.\r\n        System.out.println(\"Hello, World\");\r\n    }\r\n\r\n}"
	results := NewDetectionResults()
	content := []byte(safeJavaCode)
	filename := "filename"
	additions := []git_repo.Addition{git_repo.NewAddition(filename, content)}

	NewFileContentDetector().Test(additions, NewIgnores(), results)
	assert.False(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
}

func TestShouldNotFlagPotentialSecretsWithinSafeLongMethodName(t *testing.T) {
	const safeLongMethodName string = "TestBase64DetectorShouldNotDetectLongMethodNamesEvenWithRidiculousHighEntropyWordsMightExist"
	results := NewDetectionResults()
	content := []byte(safeLongMethodName)
	filename := "filename"
	additions := []git_repo.Addition{git_repo.NewAddition(filename, content)}

	NewFileContentDetector().Test(additions, NewIgnores(), results)
	assert.False(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
}


func TestShouldFlagPotentialSecretsEncodedInHex(t *testing.T) {
	const hex string = "68656C6C6F20776F726C6421"
	results := NewDetectionResults()
	content := []byte(hex)
	filename := "filename"
	additions := []git_repo.Addition{git_repo.NewAddition(filename, content)}

	NewFileContentDetector().Test(additions, NewIgnores(), results)
	assert.True(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
}

func TestResultsShouldContainHexTextsIfHexAndBase64ExistInFile(t *testing.T) {
	const hex string = "68656C6C6F20776F726C6421"
	const base64 string = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	const hexAndBase64 = hex + "\n" + base64
	results := NewDetectionResults()
	content := []byte(hexAndBase64)
	filename := "filename"
	additions := []git_repo.Addition{git_repo.NewAddition(filename, content)}
	filePath := additions[0].Path

	NewFileContentDetector().Test(additions, NewIgnores(), results)
	expectedMsg := "Expected file to not to contain base64 or hex encoded texts such as: " + hex
	assert.Equal(t, expectedMsg, results.Failures(filePath)[0])
}

func TestResultsShouldContainBase64TextsIfHexAndBase64ExistInFile(t *testing.T) {
	const hex string = "68656C6C6F20776F726C6421"
	const base64 string = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	const hexAndBase64 = hex + "\n" + base64
	results := NewDetectionResults()
	content := []byte(hexAndBase64)
	filename := "filename"
	additions := []git_repo.Addition{git_repo.NewAddition(filename, content)}
	filePath := additions[0].Path

	NewFileContentDetector().Test(additions, NewIgnores(), results)
	expectedMsg := "Expected file to not to contain base64 or hex encoded texts such as: " + base64
	assert.Equal(t, expectedMsg, results.Failures(filePath)[1])
}
