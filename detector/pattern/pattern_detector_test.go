package pattern

import (
	"regexp"
	"strings"
	"talisman/detector/helpers"
	"talisman/detector/severity"
	"talisman/gitrepo"
	"talisman/talismanrc"
	"talisman/utility"
	"testing"

	"github.com/stretchr/testify/assert"
)

var talismanRC = &talismanrc.TalismanRC{}
var defaultChecksumCompare = helpers.NewChecksumCompare(nil, utility.MakeHasher("default", "."), talismanRC)
var dummyCallback = func() {}

var (
	customPatterns []talismanrc.PatternString
)

func TestShouldDetectPasswordPatterns(t *testing.T) {
	filename := "secret.txt"
	values := [7]string{"password", "secret", "key", "pwd", "pass", "pword", "passphrase"}
	for i := 0; i < len(values); i++ {
		shouldPassDetectionOfSecretPattern(filename, []byte(strings.ToTitle(values[i])+":UnsafeString"), t)
		shouldPassDetectionOfSecretPattern(filename, []byte(values[i]+"=UnsafeString"), t)
		shouldPassDetectionOfSecretPattern(filename, []byte("."+values[i]+"=randomStringGoesHere}"), t)
		shouldPassDetectionOfSecretPattern(filename, []byte(":"+values[i]+" randomStringGoesHere"), t)
		shouldPassDetectionOfSecretPattern(filename,
			[]byte("\"SERVER_"+strings.ToUpper(values[i])+"\" : UnsafeString"),
			t)
		shouldPassDetectionOfSecretPattern(filename, []byte(values[i]+"2-string : UnsafeString"), t)
		shouldPassDetectionOfSecretPattern(filename,
			[]byte("<"+values[i]+" data=123> randomStringGoesHere </"+values[i]+">"),
			t)
		shouldPassDetectionOfSecretPattern(filename,
			[]byte("<admin "+values[i]+"> randomStringGoesHere </my"+values[i]+">"),
			t)
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
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	content := []byte("\"password\" : UnsafePassword")
	filename := "secret.txt"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}
	fileIgnoreConfig := &talismanrc.FileIgnoreConfig{
		FileName:        filename,
		Checksum:        "833b6c24c8c2c5c7e1663226dc401b29c005492dc76a1150fc0e0f07f29d4cc3",
		IgnoreDetectors: []string{"filecontent"},
		AllowedPatterns: []string{}}
	ignores := &talismanrc.TalismanRC{IgnoreConfigs: []talismanrc.IgnoreConfig{fileIgnoreConfig}}

	NewPatternDetector(customPatterns).Test(defaultChecksumCompare, additions, ignores, results, dummyCallback)

	assert.True(t, results.Successful(), "Expected file %s to be ignored by pattern", filename)
}

func TestShouldIgnoreAllowedPattern(t *testing.T) {
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	content := []byte("\"key\" : \"This is an allowed keyword\"\npassword=y0uw1lln3v3rgu3ssmyP@55w0rd")
	filename := "allowed_pattern.txt"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}
	fileIgnoreConfig := &talismanrc.FileIgnoreConfig{
		FileName: filename, Checksum: "",
		IgnoreDetectors: []string{},
		AllowedPatterns: []string{"key"}}
	ignores := &talismanrc.TalismanRC{
		IgnoreConfigs:   []talismanrc.IgnoreConfig{fileIgnoreConfig},
		AllowedPatterns: []*regexp.Regexp{regexp.MustCompile("password")}}

	NewPatternDetector(customPatterns).Test(defaultChecksumCompare, additions, ignores, results, dummyCallback)

	assert.True(t,
		results.Successful(),
		"Expected keywords %v %v to be ignored by Talisman", fileIgnoreConfig.AllowedPatterns, ignores.AllowedPatterns)
}
func TestShouldOnlyWarnSecretPatternIfBelowThreshold(t *testing.T) {
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	content := []byte(`password=UnsafeString`)
	filename := "secret.txt"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}
	talismanRCWithThreshold := &talismanrc.TalismanRC{Threshold: severity.High}
	checksumCompare := helpers.NewChecksumCompare(nil, utility.MakeHasher("default", "."), talismanRCWithThreshold)

	NewPatternDetector(customPatterns).Test(checksumCompare, additions, talismanRCWithThreshold, results, dummyCallback)

	assert.False(t, results.HasFailures(), "Expected file %s to not have failures", filename)
	assert.True(t, results.HasWarnings(), "Expected file %s to have warnings", filename)
}

func DetectionOfSecretPattern(filename string, content []byte) (*helpers.DetectionResults, []gitrepo.Addition, string) {
	results := helpers.NewDetectionResults(talismanrc.HookMode)
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}
	NewPatternDetector(customPatterns).Test(defaultChecksumCompare, additions, talismanRC, results, dummyCallback)
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

func getFailureMessage(results *helpers.DetectionResults, additions []gitrepo.Addition) string {
	failureMessages := []string{}
	for _, failureDetails := range results.GetFailures(additions[0].Path) {
		failureMessages = append(failureMessages, failureDetails.Message)
	}
	if len(failureMessages) == 0 {
		return ""
	}
	return failureMessages[0]
}
