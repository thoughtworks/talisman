package filecontent

import (
	"fmt"
	"regexp"
	"strings"
	"talisman/detector/helpers"
	"talisman/detector/severity"
	"talisman/gitrepo"
	mock "talisman/internal/mock/checksumcalculator"
	"talisman/talismanrc"
	"talisman/utility"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
)

var emptyTalismanRC = &talismanrc.TalismanRC{IgnoreConfigs: []talismanrc.IgnoreConfig{}}
var defaultChecksumCompareUtility = helpers.
	NewChecksumCompare(nil, utility.MakeHasher("default", "."), emptyTalismanRC)
var dummyCallback = func() {}
var filename = "filename"

func TestShouldNotFlagSafeText(t *testing.T) {
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, []byte("prettySafe"))}

	NewFileContentDetector(emptyTalismanRC).
		Test(defaultChecksumCompareUtility, additions, emptyTalismanRC, results, dummyCallback)
	assert.False(t, results.HasFailures(), "Expected file to not contain base64 encoded texts")
}

func TestShouldIgnoreFileIfNeeded(t *testing.T) {
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, []byte("prettySafe"))}
	talismanRCIWithFilenameIgnore := &talismanrc.TalismanRC{
		IgnoreConfigs: []talismanrc.IgnoreConfig{
			&talismanrc.FileIgnoreConfig{FileName: filename},
		},
	}
	mockChecksumCalculator := mock.NewMockChecksumCalculator(gomock.NewController(t))
	mockChecksumCalculator.EXPECT().
		CalculateCollectiveChecksumForPattern("filename").
		Return("mock-checksum-for-filename")
	checksumCompare := helpers.
		NewChecksumCompare(mockChecksumCalculator, utility.MakeHasher("default", "."), talismanRCIWithFilenameIgnore)

	NewFileContentDetector(talismanRCIWithFilenameIgnore).
		Test(checksumCompare, additions, talismanRCIWithFilenameIgnore, results, dummyCallback)

	assert.True(t, results.Successful(), "Expected file %s to be ignored by pattern", filename)
}

func TestShouldNotFlag4CharSafeText(t *testing.T) {
	/*This only tell that an input could have been a b64 encoded value, but it does not tell whether or not the
		input is actually a b64 encoded value. In other words, abcd will match, but it is not necessarily represent
		 the encoded value of i· rather just a plain abcd input see
	stackoverflow.com/questions/8571501/how-to-check-whether-the-string-is-base64-encoded-or-not#comment23919648_8571649*/
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, []byte("abcd"))}

	NewFileContentDetector(emptyTalismanRC).
		Test(defaultChecksumCompareUtility, additions, emptyTalismanRC, results, dummyCallback)
	assert.False(t, results.HasFailures(), "Expected file to not contain base64 encoded texts")
}

func TestShouldNotFlagLowEntropyBase64Text(t *testing.T) {
	const lowEntropyString string = "YWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWEK"
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	content := []byte(lowEntropyString)
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}

	NewFileContentDetector(emptyTalismanRC).
		Test(defaultChecksumCompareUtility, additions, emptyTalismanRC, results, dummyCallback)
	assert.False(t, results.HasFailures(), "Expected file to not contain base64 encoded texts")
}

func TestShouldFlagPotentialAWSSecretKeys(t *testing.T) {
	const awsSecretAccessKey string = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, []byte(awsSecretAccessKey))}
	filePath := additions[0].Path

	NewFileContentDetector(emptyTalismanRC).
		Test(defaultChecksumCompareUtility, additions, emptyTalismanRC, results, dummyCallback)

	expectedMessage := fmt.
		Sprintf("Expected file to not contain base64 encoded texts such as: %s", awsSecretAccessKey)
	assert.True(t, results.HasFailures(), "Expected file to not contain base64 encoded texts")
	assert.Equal(t, expectedMessage, getFailureMessages(results, filePath)[0])
	assert.Len(t, results.Results, 1)
}

func TestShouldFlagPotentialSecretWithoutTrimmingWhenLengthLessThan50Characters(t *testing.T) {
	const secret string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9asdfa"
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, []byte(secret))}
	filePath := additions[0].Path

	NewFileContentDetector(emptyTalismanRC).
		Test(defaultChecksumCompareUtility, additions, emptyTalismanRC, results, dummyCallback)

	expectedMessage := fmt.Sprintf("Expected file to not contain base64 encoded texts such as: %s", secret)
	assert.True(t, results.HasFailures(), "Expected file to not contain base64 encoded texts")
	assert.Equal(t, expectedMessage, getFailureMessages(results, filePath)[0])
	assert.Len(t, results.Results, 1)
}

func TestShouldFlagPotentialJWT(t *testing.T) {
	const jwt string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzY290Y2guaW8iLCJleHAiOjEzMDA4MTkzODAsIm5hbWUi" +
		"OiJDaHJpcyBTZXZpbGxlamEiLCJhZG1pbiI6dHJ1ZX0.03f329983b86f7d9a9f5fef85305880101d5e302afafa20154d094b229f757"
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	content := []byte(jwt)
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}
	filePath := additions[0].Path

	NewFileContentDetector(emptyTalismanRC).
		Test(defaultChecksumCompareUtility, additions, emptyTalismanRC, results, dummyCallback)

	expectedMessage := fmt.
		Sprintf("Expected file to not contain base64 encoded texts such as: %s", jwt[:47]+"...")
	assert.True(t, results.HasFailures(), "Expected file to not contain base64 encoded texts")
	assert.Equal(t, expectedMessage, getFailureMessages(results, filePath)[0])
	assert.Len(t, results.Results, 1)
}

func TestShouldFlagPotentialSecretsWithinJavaCode(t *testing.T) {
	const dangerousJavaCode string = "public class HelloWorld {\r\n\r\n" +
		"    public static void main(String[] args) {\r\n        " +
		"		// Prints \"Hello, World\" to the terminal window.\r\n  " +
		"		String accessKey=\"wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY\";\r\n        " +
		"		System.out.println(\"Hello, World\");\r\n    " +
		"		}\r\n\r\n" +
		"}"
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	content := []byte(dangerousJavaCode)
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}
	filePath := additions[0].Path

	NewFileContentDetector(emptyTalismanRC).
		Test(defaultChecksumCompareUtility, additions, emptyTalismanRC, results, dummyCallback)
	expectedMessage := "Expected file to not contain base64 encoded texts such as: " +
		"accessKey=\"wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPL..."
	assert.True(t, results.HasFailures(), "Expected file to not contain base64 encoded texts")
	assert.Equal(t, expectedMessage, getFailureMessages(results, filePath)[0])
	assert.Len(t, results.Results, 1)
}

func TestShouldNotFlagPotentialSecretsWithinSafeJavaCode(t *testing.T) {
	const safeJavaCode string = "public class HelloWorld {\r\n\r\n" +
		"    public static void main(String[] args) {\r\n        " +
		"		//Prints \"Hello, World\" to the terminal window.\r\n        " +
		"   	System.out.println(\"Hello, World\");\r\n    " +
		"	}\r\n\r\n" +
		"}"
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, []byte(safeJavaCode))}

	NewFileContentDetector(emptyTalismanRC).
		Test(defaultChecksumCompareUtility, additions, emptyTalismanRC, results, dummyCallback)
	assert.False(t, results.HasFailures(), "Expected file to not contain base64 encoded texts")
}

func TestShouldNotFlagPotentialSecretsWithinSafeLongMethodName(t *testing.T) {
	safeLongMethodName := "TestBase64DetectorShouldNotDetectLongMethodNamesEvenWithRidiculousHighEntropyWordsMightExist"
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, []byte(safeLongMethodName))}

	NewFileContentDetector(emptyTalismanRC).
		Test(defaultChecksumCompareUtility, additions, emptyTalismanRC, results, dummyCallback)
	assert.False(t, results.HasFailures(), "Expected file to not contain base64 encoded texts")
}

func TestShouldFlagPotentialSecretsEncodedInHex(t *testing.T) {
	const hex string = "68656C6C6F20776F726C6421"
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, []byte(hex))}
	filePath := additions[0].Path

	NewFileContentDetector(emptyTalismanRC).
		Test(defaultChecksumCompareUtility, additions, emptyTalismanRC, results, dummyCallback)
	expectedMessage := "Expected file to not contain hex encoded texts such as: " + hex
	assert.Equal(t, expectedMessage, getFailureMessages(results, filePath)[0])
	assert.Len(t, results.Results, 1)
}

func TestShouldNotFlagPotentialCreditCardNumberIfAboveThreshold(t *testing.T) {
	const creditCardNumber string = "340000000000009"
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, []byte(creditCardNumber))}
	talismanRCWithThreshold := &talismanrc.TalismanRC{Threshold: severity.High}
	checksumCompareWithThreshold := helpers.
		NewChecksumCompare(nil, utility.MakeHasher("default", "."), talismanRCWithThreshold)

	NewFileContentDetector(emptyTalismanRC).
		Test(checksumCompareWithThreshold, additions, talismanRCWithThreshold, results, dummyCallback)

	assert.False(t, results.HasFailures(), "Expected no base64 detection when threshold is higher")
}

func TestShouldNotFlagPotentialSecretsIfIgnored(t *testing.T) {
	const hex string = "68656C6C6F20776F726C6421"
	talismanRCWithIgnores := &talismanrc.TalismanRC{
		AllowedPatterns: []*regexp.Regexp{regexp.MustCompile("[0-9a-fA-F]*")}}
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, []byte(hex))}

	NewFileContentDetector(emptyTalismanRC).
		Test(defaultChecksumCompareUtility, additions, talismanRCWithIgnores, results, dummyCallback)

	assert.False(t, results.HasFailures(), "Expected file ignore allowed pattern for hex text")
}

func TestResultsShouldNotFlagCreditCardNumberIfSpecifiedInFileIgnores(t *testing.T) {
	const creditCardNumber string = "340000000000009"
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	fileIgnoreConfig := &talismanrc.FileIgnoreConfig{
		FileName: filename, Checksum: "",
		AllowedPatterns: []string{creditCardNumber},
	}
	talismanRCWithFileIgnore := &talismanrc.TalismanRC{
		IgnoreConfigs:   []talismanrc.IgnoreConfig{fileIgnoreConfig},
	}
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, []byte(creditCardNumber))}

	NewFileContentDetector(emptyTalismanRC).
		Test(defaultChecksumCompareUtility, additions, talismanRCWithFileIgnore, results, dummyCallback)
	
	assert.False(t, results.HasFailures(), "Expected the creditcard number to be ignored based on talisman RC")

}


func TestResultsShouldContainHexTextsIfHexAndBase64ExistInFile(t *testing.T) {
	const hex string = "68656C6C6F20776F726C6421"
	const base64 string = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	const hexAndBase64 = hex + "\n" + base64
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, []byte(hexAndBase64))}
	filePath := additions[0].Path

	NewFileContentDetector(emptyTalismanRC).
		Test(defaultChecksumCompareUtility, additions, emptyTalismanRC, results, dummyCallback)
	expectedMessage := "Expected file to not contain hex encoded texts such as: " + hex
	messageReceived := strings.Join(getFailureMessages(results, filePath), " ")
	assert.Regexp(t, expectedMessage, messageReceived, "Should contain hex detection message")
	assert.Len(t, results.Results, 1)
}

func TestResultsShouldContainBase64TextsIfHexAndBase64ExistInFile(t *testing.T) {
	const hex string = "68656C6C6F20776F726C6421"
	const base64 string = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	const hexAndBase64 = hex + "\n" + base64
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, []byte(hexAndBase64))}
	filePath := additions[0].Path

	NewFileContentDetector(emptyTalismanRC).
		Test(defaultChecksumCompareUtility, additions, emptyTalismanRC, results, dummyCallback)

	expectedMessage := "Expected file to not contain base64 encoded texts such as: " + base64
	messageReceived := strings.Join(getFailureMessages(results, filePath), " ")
	assert.Regexp(t, expectedMessage, messageReceived, "Should contain base64 detection message")
	assert.Len(t, results.Results, 1)
}

func TestResultsShouldContainCreditCardNumberIfCreditCardNumberExistInFile(t *testing.T) {
	const creditCardNumber string = "340000000000009"
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, []byte(creditCardNumber))}
	filePath := additions[0].Path

	NewFileContentDetector(emptyTalismanRC).
		Test(defaultChecksumCompareUtility, additions, emptyTalismanRC, results, dummyCallback)

	expectedMessage := "Expected file to not contain credit card numbers such as: " + creditCardNumber
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
