package talismanrc

import (
	"io/ioutil"
	"os"
	"testing"

	"talisman/detector/severity"
	"talisman/gitrepo"

	logr "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	logr.SetOutput(ioutil.Discard)
}

func TestShouldIgnoreEmptyLinesInTheFile(t *testing.T) {
	defaultRepoFileReader := repoFileReader()
	for _, s := range []string{"", " ", "  ", "\t", " \t", "\t\t \t"} {
		setRepoFileReader(func(string) ([]byte, error) {
			return []byte(s), nil
		})

		talismanRC := For(HookMode)
		assert.True(t, talismanRC.AcceptsAll(), "Expected '%s' to result in no ignore patterns.", s)
		talismanRC = For(ScanMode)
		assert.True(t, talismanRC.AcceptsAll(), "Expected '%s' to result in no ignore patterns.", s)
	}
	setRepoFileReader(defaultRepoFileReader)
}

func TestShouldIgnoreUnformattedFiles(t *testing.T) {
	defaultRepoFileReader := repoFileReader()
	for _, s := range []string{"#", "#monkey", "# this monkey likes bananas  "} {
		setRepoFileReader(func(string) ([]byte, error) {
			return []byte(s), nil
		})

		talismanRC := For(HookMode)
		assert.True(t, talismanRC.AcceptsAll(), "Expected commented line '%s' to result in no ignore patterns", s)
		talismanRC = For(ScanMode)
		assert.True(t, talismanRC.AcceptsAll(), "Expected commented line '%s' to result in no ignore patterns", s)
	}
	setRepoFileReader(defaultRepoFileReader)
}

func TestShouldConvertThresholdToValue(t *testing.T) {
	talismanRCContents := []byte("threshold: high")
	assert.Equal(t, newPersistedRC(talismanRCContents).Threshold, severity.High)
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
	file6 := testAddition("imgJpeg.jpeg")
	file7 := testAddition("imgJpg.jpg")
	file8 := testAddition("imgPng.png")
	file9 := testAddition("build.bzl")
	additions := []gitrepo.Addition{file1, file2, file3, file4, file5, file6, file7, file8, file9}

	scopesToIgnore := []string{"node", "go", "images", "bazel"}
	talismanRCConfig := createTalismanRCWithScopeIgnores(scopesToIgnore)

	nodeIgnores := []string{"node.lock", "*yarn.lock"}
	javaIgnores := []string{"java.lock"}
	goIgnores := []string{"go.lock", "Gopkg.lock", "vendors/"}
	imageIgnores := []string{"*.jpeg", "*.jpg", "*.png"}
	bazelIgnores := []string{"*.bzl"}
	scopesMap := map[string][]string{"node": nodeIgnores, "java": javaIgnores, "go": goIgnores, "images": imageIgnores, "bazel": bazelIgnores}
	knownScopes = scopesMap
	filteredAdditions := talismanRCConfig.FilterAdditions(additions)

	assert.NotContains(t, filteredAdditions, file1)
	assert.NotContains(t, filteredAdditions, file2)
	assert.Contains(t, filteredAdditions, file3)
	assert.NotContains(t, filteredAdditions, file4)
	assert.NotContains(t, filteredAdditions, file5)
	assert.NotContains(t, filteredAdditions, file6)
	assert.NotContains(t, filteredAdditions, file7)
	assert.NotContains(t, filteredAdditions, file8)
	assert.NotContains(t, filteredAdditions, file9)
}

func TestIgnoringDetectors(t *testing.T) {
	assertDeniesDetector("foo", "someDetector", "foo", "someDetector", t)
	assertAcceptsDetector("foo", "someDetector", "foo", "someOtherDetector", t)
}

func TestMakeWithFileIgnores(t *testing.T) {
	ignoreConfigs := []FileIgnoreConfig{}
	builtConfig := MakeWithFileIgnores(ignoreConfigs)
	assert.Equal(t, builtConfig.FileIgnoreConfig, ignoreConfigs)
	assert.Equal(t, builtConfig.Version, DefaultRCVersion)
}

func TestBuildIgnoreConfig(t *testing.T) {
	ignoreConfig := BuildIgnoreConfig(HookMode, "filename", "asdfasdfasdfasdfasdf", nil)
	assert.IsType(t, &FileIgnoreConfig{}, ignoreConfig)
}

func TestAddIgnoreFilesInHookMode(t *testing.T) {
	ignoreConfig := &FileIgnoreConfig{
		FileName:        "Foo",
		Checksum:        "SomeCheckSum",
		IgnoreDetectors: []string{},
		AllowedPatterns: []string{}}
	os.Remove(DefaultRCFileName)
	talismanRCConfig := createTalismanRCWithScopeIgnores([]string{})
	talismanRCConfig.base.AddIgnores(HookMode, []IgnoreConfig{ignoreConfig})
	talismanRCConfigFromFile := ConfigFromFile()
	assert.Equal(t, 1, len(talismanRCConfigFromFile.FileIgnoreConfig))
	os.Remove(DefaultRCFileName)
}

func assertDenies(line, ignoreDetector string, path string, t *testing.T) {
	assertDeniesDetector(line, ignoreDetector, path, "filename", t)
}

func assertDeniesDetector(line, ignoreDetector string, path string, detectorName string, t *testing.T) {
	assert.True(t, createTalismanRCWithFileIgnores(line, ignoreDetector).Deny(testAddition(path), detectorName), "%s is expected to deny a file named %s.", line, path)
}

func assertAccepts(line, ignoreDetector string, path string, t *testing.T, detectorNames ...string) {
	assertAcceptsDetector(line, ignoreDetector, path, "someDetector", t)
}

func assertAcceptsDetector(line, ignoreDetector string, path string, detectorName string, t *testing.T) {
	assert.True(t, createTalismanRCWithFileIgnores(line, ignoreDetector).Accept(testAddition(path), detectorName), "%s is expected to accept a file named %s.", line, path)
}

func testAddition(path string) gitrepo.Addition {
	return gitrepo.NewAddition(path, make([]byte, 0))
}

func createTalismanRCWithFileIgnores(filename string, detector string) *TalismanRC {
	fileIgnoreConfig := &FileIgnoreConfig{}
	fileIgnoreConfig.FileName = filename
	if detector != "" {
		fileIgnoreConfig.IgnoreDetectors = []string{detector}
	}

	return &TalismanRC{IgnoreConfigs: []IgnoreConfig{fileIgnoreConfig}}
}

func createTalismanRCWithScopeIgnores(scopesToIgnore []string) *TalismanRC {
	var scopeConfigs []ScopeConfig
	for _, scopeIgnore := range scopesToIgnore {
		scopeIgnoreConfig := ScopeConfig{}
		scopeIgnoreConfig.ScopeName = scopeIgnore
		scopeConfigs = append(scopeConfigs, scopeIgnoreConfig)
	}

	return &TalismanRC{ScopeConfig: scopeConfigs}
}

func TestFileIgnoreConfig_ChecksumMatches(t *testing.T) {
	fileIgnoreConfig := &FileIgnoreConfig{
		FileName:        "some_filename",
		Checksum:        "some_checksum",
		IgnoreDetectors: nil,
		AllowedPatterns: nil,
	}

	assert.True(t, fileIgnoreConfig.ChecksumMatches("some_checksum"))
	assert.False(t, fileIgnoreConfig.ChecksumMatches("some_other_checksum"))
}

func TestFileIgnoreConfig_GetAllowedPatterns(t *testing.T) {
	fileIgnoreConfig := &FileIgnoreConfig{
		FileName:        "some_filename",
		Checksum:        "some_checksum",
		IgnoreDetectors: nil,
		AllowedPatterns: nil,
	}

	//No allowed patterns specified
	allowedPatterns := fileIgnoreConfig.GetAllowedPatterns()
	assert.Equal(t, 0, len(allowedPatterns))

	fileIgnoreConfig.compiledPatterns = nil
	fileIgnoreConfig.AllowedPatterns = []string{"[Ff]ile[nN]ame"}
	allowedPatterns = fileIgnoreConfig.GetAllowedPatterns()
	assert.Equal(t, 1, len(allowedPatterns))
	assert.Regexp(t, allowedPatterns[0], "fileName")
}

func TestSuggestRCFor(t *testing.T) {
	t.Run("should suggest proper RC when ignore configs are valid", func(t *testing.T) {
		fileIgnoreConfigs := []IgnoreConfig{
			&FileIgnoreConfig{
				FileName: "some_filename",
				Checksum: "some_checksum",
			},
		}
		expectedRC := `fileignoreconfig:
- filename: some_filename
  checksum: some_checksum
version: ""
`
		str := SuggestRCFor(fileIgnoreConfigs)
		assert.Equal(t, expectedRC, str)
	})

	t.Run("should ignore invalid configs", func(t *testing.T) {
		fileIgnoreConfigs := []IgnoreConfig{
			&FileIgnoreConfig{
				FileName: "some_filename",
				Checksum: "some_checksum",
			},
		}
		expectedRC := `fileignoreconfig:
- filename: some_filename
  checksum: some_checksum
version: ""
`
		str := SuggestRCFor(fileIgnoreConfigs)
		assert.Equal(t, expectedRC, str)
	})
}

func TestFor(t *testing.T) {
	var repoFileReader = func(string) ([]byte, error) {
		return []byte(`fileignoreconfig:
- filename: testfile_1.yml
  checksum: file1_checksum
- filename: testfile_2.yml
  checksum: file2_checksum
- filename: testfile_3.yml
  checksum: file3_checksum`), nil
	}
	t.Run("talismanrc.For(mode) should read multiple entries in rc file correctly", func(t *testing.T) {
		setRepoFileReader(repoFileReader)
		rc := For(HookMode)
		assert.Equal(t, 3, len(rc.IgnoreConfigs))

		assert.Equal(t, rc.IgnoreConfigs[0].GetFileName(), "testfile_1.yml")
		assert.True(t, rc.IgnoreConfigs[0].ChecksumMatches("file1_checksum"))
		assert.Equal(t, rc.IgnoreConfigs[1].GetFileName(), "testfile_2.yml")
		assert.True(t, rc.IgnoreConfigs[1].ChecksumMatches("file2_checksum"))
		assert.Equal(t, rc.IgnoreConfigs[2].GetFileName(), "testfile_3.yml")
		assert.True(t, rc.IgnoreConfigs[2].ChecksumMatches("file3_checksum"))

	})
}
