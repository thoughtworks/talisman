package filecontent

import (
	"talisman/detector/helpers"
	"talisman/utility"
	"testing"

	"talisman/gitrepo"
	"talisman/talismanrc"

	"github.com/stretchr/testify/assert"
)

var _blankTalismanRC = &talismanrc.TalismanRC{}
var dummyCompletionCallbackFunc = func() {}
var aggressiveModeFileContentDetector = NewFileContentDetector(_blankTalismanRC).AggressiveMode()

func TestShouldFlagPotentialAWSAccessKeysInAggressiveMode(t *testing.T) {
	const awsAccessKeyIDExample string = "AKIAIOSFODNN7EXAMPLE\n"
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	filename := "filename"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, []byte(awsAccessKeyIDExample))}

	aggressiveModeFileContentDetector.
		Test(
			helpers.NewChecksumCompare(nil, utility.MakeHasher("default", "."), _blankTalismanRC),
			additions,
			_blankTalismanRC,
			results,
			dummyCompletionCallbackFunc)

	assert.True(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
}

func TestShouldFlagPotentialAWSAccessKeysAtPropertyDefinitionInAggressiveMode(t *testing.T) {
	const awsAccessKeyIDExample string = "accessKey=AKIAIOSFODNN7EXAMPLE"
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	filename := "filename"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, []byte(awsAccessKeyIDExample))}

	aggressiveModeFileContentDetector.
		Test(
			helpers.NewChecksumCompare(nil, utility.MakeHasher("default", "."), _blankTalismanRC),
			additions,
			_blankTalismanRC,
			results,
			dummyCompletionCallbackFunc)

	assert.True(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
}

func TestShouldNotFlagPotentialSecretsWithinSafeJavaCodeEvenInAggressiveMode(t *testing.T) {
	const awsAccessKeyIDExample string = "public class HelloWorld {\r\n\r\n" +
		"   public static void main(String[] args) {\r\n        " +
		"		// Prints \"Hello, World\" to the terminal window.\r\n        " +
		"		System.out.println(\"Hello, World\");\r\n    " +
		"	}\r\n\r\n" +
		"}"
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	filename := "filename"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, []byte(awsAccessKeyIDExample))}

	aggressiveModeFileContentDetector.
		Test(
			helpers.NewChecksumCompare(nil, utility.MakeHasher("default", "."), _blankTalismanRC),
			additions,
			_blankTalismanRC,
			results,
			dummyCompletionCallbackFunc)

	assert.False(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
}
