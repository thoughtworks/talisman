package filecontent

import (
	"talisman/detector/helpers"
	"talisman/utility"
	"testing"

	"talisman/gitrepo"
	"talisman/talismanrc"

	"github.com/stretchr/testify/assert"
)

func TestShouldFlagPotentialAWSAccessKeysInAggressiveMode(t *testing.T) {
	const awsAccessKeyIDExample string = "AKIAIOSFODNN7EXAMPLE\n"
	results := helpers.NewDetectionResults(talismanrc.Hook)
	content := []byte(awsAccessKeyIDExample)
	filename := "filename"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}

	NewFileContentDetector(talismanRC).AggressiveMode().Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additions, talismanRC, results, func() {})
	assert.True(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
}

func TestShouldFlagPotentialAWSAccessKeysAtPropertyDefinitionInAggressiveMode(t *testing.T) {
	const awsAccessKeyIDExample string = "accessKey=AKIAIOSFODNN7EXAMPLE"
	results := helpers.NewDetectionResults(talismanrc.Hook)
	content := []byte(awsAccessKeyIDExample)
	filename := "filename"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}

	NewFileContentDetector(talismanRC).AggressiveMode().Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additions, talismanRC, results, func() {})
	assert.True(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
}

func TestShouldNotFlagPotentialSecretsWithinSafeJavaCodeEvenInAggressiveMode(t *testing.T) {
	const awsAccessKeyIDExample string = "public class HelloWorld {\r\n\r\n    public static void main(String[] args) {\r\n        // Prints \"Hello, World\" to the terminal window.\r\n        System.out.println(\"Hello, World\");\r\n    }\r\n\r\n}"
	results := helpers.NewDetectionResults(talismanrc.Hook)
	content := []byte(awsAccessKeyIDExample)
	filename := "filename"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}

	NewFileContentDetector(talismanRC).AggressiveMode().Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additions, talismanRC, results, func() {})
	if results == nil {
		additions = nil
	}
	assert.False(t, results.HasFailures(), "Expected file to not to contain base64 encoded texts")
}
