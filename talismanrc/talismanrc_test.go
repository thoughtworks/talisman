package talismanrc

import (
	"io/ioutil"
	"testing"

	"talisman/detector/severity"
	"talisman/gitrepo"

	logr "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	logr.SetOutput(ioutil.Discard)
}

func TestShouldIgnoreEmptyLinesInTheFile(t *testing.T) {
	for _, s := range []string{"", " ", "  "} {
		assert.True(t, NewTalismanRC([]byte(s)).AcceptsAll(), "Expected '%s' to result in no ignore patterns.", s)
	}
}

func TestShouldIgnoreUnformattedFiles(t *testing.T) {
	for _, s := range []string{"#", "#monkey", "# this monkey likes bananas  "} {
		assert.True(t, NewTalismanRC([]byte(s)).AcceptsAll(), "Expected commented line '%s' to result in no ignore patterns", s)
	}
}

func TestShouldConvertThresholdToValue(t *testing.T) {
	talismanRCContents := []byte("threshold: high")
	assert.Equal(t, NewTalismanRC(talismanRCContents).Threshold, severity.High)
}

func TestDirectoryPatterns(t *testing.T) {
	assertAccepts("foo/", "", "bar", t)
	assertAccepts("foo/", "", "foo", t)
	assertDenies("foo/", "filename", "foo/bar", t)
	assertDenies("foo/", "filename", "foo/bar.txt", t)
	assertDenies("foo/", "filename", "foo/bar/baz.txt", t)
}

func TestIgnoreAdditionsByScope(t *testing.T) {
	file1 := testAddition("yarn.lock")
	file2 := testAddition("similaryarn.lock")
	file3 := testAddition("java.lock")
	file4 := testAddition("Gopkg.lock")
	file5 := testAddition("vendors/abc")
	additions := []gitrepo.Addition{file1, file2, file3, file4, file5}

	scopesToIgnore := []string{"node", "go"}
	talismanRCConfig := CreatetalismanRCWithScopeIgnore(scopesToIgnore)

	nodeIgnores := []string{"node.lock", "*yarn.lock"}
	javaIgnores := []string{"java.lock"}
	goIgnores := []string{"go.lock", "Gopkg.lock", "vendors/"}
	scopesMap := map[string][]string{"node": nodeIgnores, "java": javaIgnores, "go": goIgnores}
	knownScopes = scopesMap
	filteredAdditions := talismanRCConfig.FilterAdditions(additions)

	assert.NotContains(t, filteredAdditions, file1)
	assert.NotContains(t, filteredAdditions, file2)
	assert.Contains(t, filteredAdditions, file3)
	assert.NotContains(t, filteredAdditions, file4)
	assert.NotContains(t, filteredAdditions, file5)
}

func TestIgnoringDetectors(t *testing.T) {
	assertDeniesDetector("foo", "someDetector", "foo", "someDetector", t)
	assertAcceptsDetector("foo", "someDetector", "foo", "someOtherDetector", t)
}

func TestAddIgnoreFiles(t *testing.T) {
	talismanRCConfig := CreatetalismanRCWithScopeIgnore([]string{})
	talismanRCConfig.AddIgnores(Hook, []IgnoreConfig{&FileIgnoreConfig{"Foo", "SomeCheckSum", []string{}, []string{}}})
	talismanRCConfig = Get()
	assert.Equal(t, 1, len(talismanRCConfig.FileIgnoreConfig))
}

func assertDenies(line, ignoreDetector string, path string, t *testing.T) {
	assertDeniesDetector(line, ignoreDetector, path, "filename", t)
}

func assertDeniesDetector(line, ignoreDetector string, path string, detectorName string, t *testing.T) {
	assert.True(t, CreatetalismanRCWithFileName(line, ignoreDetector).Deny(testAddition(path), detectorName), "%s is expected to deny a file named %s.", line, path)
}

func assertAccepts(line, ignoreDetector string, path string, t *testing.T, detectorNames ...string) {
	assertAcceptsDetector(line, ignoreDetector, path, "someDetector", t)
}

func assertAcceptsDetector(line, ignoreDetector string, path string, detectorName string, t *testing.T) {
	assert.True(t, CreatetalismanRCWithFileName(line, ignoreDetector).Accept(testAddition(path), detectorName), "%s is expected to accept a file named %s.", line, path)
}

func testAddition(path string) gitrepo.Addition {
	return gitrepo.NewAddition(path, make([]byte, 0))
}

func CreatetalismanRCWithFileName(filename string, detector string) *TalismanRC {
	fileIgnoreConfig := FileIgnoreConfig{}
	fileIgnoreConfig.FileName = filename
	if detector != "" {
		fileIgnoreConfig.IgnoreDetectors = make([]string, 1)
		fileIgnoreConfig.IgnoreDetectors[0] = detector
	}
	talismanRC := TalismanRC{}
	talismanRC.FileIgnoreConfig = make([]FileIgnoreConfig, 1)
	talismanRC.FileIgnoreConfig[0] = fileIgnoreConfig
	return &talismanRC
}

func CreatetalismanRCWithScopeIgnore(scopesToIgnore []string) *TalismanRC {
	var scopeConfigs []ScopeConfig
	for _, scopeIgnore := range scopesToIgnore {
		scopeIgnoreConfig := ScopeConfig{}
		scopeIgnoreConfig.ScopeName = scopeIgnore
		scopeConfigs = append(scopeConfigs, scopeIgnoreConfig)
	}

	talismanRC := TalismanRC{ScopeConfig: scopeConfigs}
	return &talismanRC
}
