package filecontent

import (
	"fmt"
	"strings"
	"talisman/detector/helpers"
	"talisman/gitrepo"
	"talisman/talismanrc"
	"talisman/utility"
	"testing"

	"github.com/stretchr/testify/assert"
)

var talismanRCContents = `
fileignoreconfig:
  - filename    : filename
`

var talismanRC = &talismanrc.TalismanRC{}

func TestShouldNotFlagSafeText(t *testing.T) {
	results := helpers.NewDetectionResults(talismanrc.Hook)
	content := []byte("prettySafe")
	filename := "filename"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}

	NewFileContentDetector(talismanRC).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additions, &talismanrc.TalismanRC{}, results, func() {})
	assert.False(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
}

func TestShouldIgnoreFileIfNeeded(t *testing.T) {
	results := helpers.NewDetectionResults(talismanrc.Hook)
	content := []byte("prettySafe")
	filename := "filename"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}

	NewFileContentDetector(talismanRC).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additions, talismanrc.NewTalismanRC([]byte(talismanRCContents)), results, func() {})
	assert.True(t, results.Successful(), "Expected file %s to be ignored by pattern", filename)
}

func TestShouldNotFlag4CharSafeText(t *testing.T) {
	/*This only tell that an input could have been a b64 encoded value, but it does not tell whether or not the
	input is actually a b64 encoded value. In other words, abcd will match, but it is not necessarily represent
	 the encoded value of i· rather just a plain abcd input
	 see stackoverflow.com/questions/8571501/how-to-check-whether-the-string-is-base64-encoded-or-not#comment23919648_8571649*/
	results := helpers.NewDetectionResults(talismanrc.Hook)
	content := []byte("abcd")
	filename := "filename"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}

	NewFileContentDetector(talismanRC).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additions, talismanRC, results, func() {})
	assert.False(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
}

func TestShouldNotFlagLowEntropyBase64Text(t *testing.T) {
	const lowEntropyString string = "YWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWEK"
	results := helpers.NewDetectionResults(talismanrc.Hook)
	content := []byte(lowEntropyString)
	filename := "filename"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}

	NewFileContentDetector(talismanRC).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additions, talismanRC, results, func() {})
	assert.False(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
}

func TestShouldFlagPotentialAWSSecretKeys(t *testing.T) {
	const awsSecretAccessKey string = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	results := helpers.NewDetectionResults(talismanrc.Hook)
	content := []byte(awsSecretAccessKey)
	filename := "filename"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}
	filePath := additions[0].Path

	NewFileContentDetector(talismanRC).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additions, talismanRC, results, func() {})
	expectedMessage := fmt.Sprintf("Expected file to not to contain base64 encoded texts such as: %s", awsSecretAccessKey)
	assert.True(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
	assert.Equal(t, expectedMessage, getFailureMessages(results, filePath)[0])
	assert.Len(t, results.Results, 1)
}

func TestShouldFlagPotentialSecretWithoutTrimmingWhenLengthLessThan50Characters(t *testing.T) {
	const secret string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9asdfa"
	results := helpers.NewDetectionResults(talismanrc.Hook)
	content := []byte(secret)
	filename := "filename"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}
	filePath := additions[0].Path

	NewFileContentDetector(talismanRC).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additions, talismanRC, results, func() {})
	expectedMessage := fmt.Sprintf("Expected file to not to contain base64 encoded texts such as: %s", secret)
	assert.True(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
	assert.Equal(t, expectedMessage, getFailureMessages(results, filePath)[0])
	assert.Len(t, results.Results, 1)
}

func TestShouldFlagPotentialJWT(t *testing.T) {
	const jwt string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzY290Y2guaW8iLCJleHAiOjEzMDA4MTkzODAsIm5hbWUiOiJDaHJpcyBTZXZpbGxlamEiLCJhZG1pbiI6dHJ1ZX0.03f329983b86f7d9a9f5fef85305880101d5e302afafa20154d094b229f757"
	results := helpers.NewDetectionResults(talismanrc.Hook)
	content := []byte(jwt)
	filename := "filename"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}
	filePath := additions[0].Path

	NewFileContentDetector(talismanRC).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additions, talismanRC, results, func() {})
	expectedMessage := fmt.Sprintf("Expected file to not to contain base64 encoded texts such as: %s", jwt[:47]+"...")
	assert.True(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
	assert.Equal(t, expectedMessage, getFailureMessages(results, filePath)[0])
	assert.Len(t, results.Results, 1)
}

func TestShouldFlagPotentialSecretsWithinJavaCode(t *testing.T) {
	const dangerousJavaCode string = "public class HelloWorld {\r\n\r\n    public static void main(String[] args) {\r\n        // Prints \"Hello, World\" to the terminal window.\r\n        accessKey=\"wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY\";\r\n        System.out.println(\"Hello, World\");\r\n    }\r\n\r\n}"
	results := helpers.NewDetectionResults(talismanrc.Hook)
	content := []byte(dangerousJavaCode)
	filename := "filename"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}
	filePath := additions[0].Path

	NewFileContentDetector(talismanRC).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additions, talismanRC, results, func() {})
	expectedMessage := "Expected file to not to contain base64 encoded texts such as: accessKey=\"wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPL..."
	assert.True(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
	assert.Equal(t, expectedMessage, getFailureMessages(results, filePath)[0])
	assert.Len(t, results.Results, 1)
}

func TestShouldNotFlagPotentialSecretsWithinSafeJavaCode(t *testing.T) {
	const safeJavaCode string = "public class HelloWorld {\r\n\r\n    public static void main(String[] args) {\r\n        // Prints \"Hello, World\" to the terminal window.\r\n        System.out.println(\"Hello, World\");\r\n    }\r\n\r\n}"
	results := helpers.NewDetectionResults(talismanrc.Hook)
	content := []byte(safeJavaCode)
	filename := "filename"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}

	NewFileContentDetector(talismanRC).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additions, talismanRC, results, func() {})
	assert.False(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
}

func TestShouldNotFlagPotentialSecretsWithinSafeLongMethodName(t *testing.T) {
	const safeLongMethodName string = "TestBase64DetectorShouldNotDetectLongMethodNamesEvenWithRidiculousHighEntropyWordsMightExist"
	results := helpers.NewDetectionResults(talismanrc.Hook)
	content := []byte(safeLongMethodName)
	filename := "filename"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}

	NewFileContentDetector(talismanRC).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additions, talismanRC, results, func() {})
	assert.False(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
}

func TestShouldFlagPotentialSecretsEncodedInHex(t *testing.T) {
	const hex string = "68656C6C6F20776F726C6421"
	results := helpers.NewDetectionResults(talismanrc.Hook)
	content := []byte(hex)
	filename := "filename"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}
	filePath := additions[0].Path

	NewFileContentDetector(talismanRC).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additions, talismanRC, results, func() {})
	expectedMessage := "Expected file to not to contain hex encoded texts such as: " + hex
	assert.Equal(t, expectedMessage, getFailureMessages(results, filePath)[0])
	assert.Len(t, results.Results, 1)
}

func TestShouldNotFlagPotentialCreditCardNumberIfAboveThreshold(t *testing.T) {
	const creditCardNumber string = "340000000000009"
	results := helpers.NewDetectionResults(talismanrc.Hook)
	content := []byte(creditCardNumber)
	filename := "filename"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}

	var talismanRCContents = "threshold: high"
	talismanRCWithThreshold := talismanrc.NewTalismanRC([]byte(talismanRCContents))

	NewFileContentDetector(talismanRC).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanRCWithThreshold), additions, talismanRCWithThreshold, results, func() {})

	assert.False(t, results.HasFailures(), "Expected file to not flag base64 encoded texts if threshold is higher")
}

func TestResultsShouldContainHexTextsIfHexAndBase64ExistInFile(t *testing.T) {
	const hex string = "68656C6C6F20776F726C6421"
	const base64 string = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	const hexAndBase64 = hex + "\n" + base64
	results := helpers.NewDetectionResults(talismanrc.Hook)
	content := []byte(hexAndBase64)
	filename := "filename"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}
	filePath := additions[0].Path

	NewFileContentDetector(talismanRC).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additions, talismanRC, results, func() {})
	expectedMessage := "Expected file to not to contain hex encoded texts such as: " + hex
	messageReceived := strings.Join(getFailureMessages(results, filePath), " ")
	assert.Regexp(t, expectedMessage, messageReceived, "Should contain hex detection message")
	assert.Len(t, results.Results, 1)
}

func TestResultsShouldContainBase64TextsIfHexAndBase64ExistInFile(t *testing.T) {
	const hex string = "68656C6C6F20776F726C6421"
	const base64 string = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	const hexAndBase64 = hex + "\n" + base64
	results := helpers.NewDetectionResults(talismanrc.Hook)
	content := []byte(hexAndBase64)
	filename := "filename"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}
	filePath := additions[0].Path

	NewFileContentDetector(talismanRC).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additions, talismanRC, results, func() {})
	expectedMessage := "Expected file to not to contain base64 encoded texts such as: " + base64
	messageReceived := strings.Join(getFailureMessages(results, filePath), " ")
	assert.Regexp(t, expectedMessage, messageReceived, "Should contain base64 detection message")
	assert.Len(t, results.Results, 1)
}

func TestResultsShouldContainCreditCardNumberIfCreditCardNumberExistInFile(t *testing.T) {
	const creditCardNumber string = "340000000000009"
	results := helpers.NewDetectionResults(talismanrc.Hook)
	content := []byte(creditCardNumber)
	filename := "filename"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}
	filePath := additions[0].Path

	NewFileContentDetector(talismanRC).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additions, talismanRC, results, func() {})
	expectedMessage := "Expected file to not to contain credit card numbers such as: " + creditCardNumber
	assert.Equal(t, expectedMessage, getFailureMessages(results, filePath)[0])
	assert.Len(t, results.Results, 1)
}

func getFailureMessages(results *helpers.DetectionResults, filePath gitrepo.FilePath) []string {
	failureMessages := []string{}
	for _, failureDetails := range results.GetFailures(filePath) {
		failureMessages = append(failureMessages, failureDetails.Message)
	}
	return failureMessages
}
