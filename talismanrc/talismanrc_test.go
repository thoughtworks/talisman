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
	additions := []gitrepo.Addition{file1, file2, file3, file4, file5, file6, file7, file8}

	scopesToIgnore := []string{"node", "go", "images"}
	talismanRCConfig := createTalismanRCWithScopeIgnores(scopesToIgnore)

	nodeIgnores := []string{"node.lock", "*yarn.lock"}
	javaIgnores := []string{"java.lock"}
	goIgnores := []string{"go.lock", "Gopkg.lock", "vendors/"}
	imageIgnores := []string{"*.jpeg", "*.jpg", "*.png"}
	scopesMap := map[string][]string{"node": nodeIgnores, "java": javaIgnores, "go": goIgnores, "images": imageIgnores}
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
	var ignoreConfig IgnoreConfig
	ignoreConfig = BuildIgnoreConfig(HookMode, "filename", "asdfasdfasdfasdfasdf", nil)
	assert.IsType(t, &FileIgnoreConfig{}, ignoreConfig)

	ignoreConfig = BuildIgnoreConfig(ScanMode, "filename", "asdfasdfasdfasdfasdf", nil)
	assert.IsType(t, &ScanFileIgnoreConfig{}, ignoreConfig)
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

func TestAddIgnoreFilesInScanMode(t *testing.T) {
	ignoreConfig := &ScanFileIgnoreConfig{
		FileName:        "Foo",
		Checksums:       []string{"SomeCheckSum"},
		IgnoreDetectors: []string{},
		AllowedPatterns: []string{}}
	os.Remove(DefaultRCFileName)
	talismanRCConfig := createTalismanRCWithScopeIgnores([]string{})
	talismanRCConfig.base.AddIgnores(ScanMode, []IgnoreConfig{ignoreConfig})
	talismanRCScanConfigFromFile := ConfigFromFile()
	assert.Equal(t, 1, len(talismanRCScanConfigFromFile.ScanConfig.FileIgnoreConfig))
	assert.Equal(t, 0, len(talismanRCScanConfigFromFile.FileIgnoreConfig))
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

func TestScanFileIgnoreConfig_ChecksumMatches(t *testing.T) {
	fileIgnoreConfig := &ScanFileIgnoreConfig{
		FileName:        "some_filename",
		Checksums:       []string{"some_checksum", "some_other_checksum"},
		IgnoreDetectors: nil,
		AllowedPatterns: nil,
	}

	assert.True(t, fileIgnoreConfig.ChecksumMatches("some_checksum"))
	assert.True(t, fileIgnoreConfig.ChecksumMatches("some_other_checksum"))
	assert.False(t, fileIgnoreConfig.ChecksumMatches("some_different_checksum"))
}

func TestScanFileIgnoreConfig_isEffective(t *testing.T) {
	fileIgnoreConfig := &ScanFileIgnoreConfig{
		FileName:        "some_filename",
		Checksums:       []string{"some_checksum", "some_other_checksum"},
		IgnoreDetectors: nil,
		AllowedPatterns: nil,
	}
	//Ignore config does not apply when detector not ignored
	assert.False(t, fileIgnoreConfig.isEffective("filename"))

	//Ignore config applies when detector explicitly ignored
	fileIgnoreConfig.IgnoreDetectors = []string{"filename"}
	assert.True(t, fileIgnoreConfig.isEffective("filename"))

	//Ignore config does not apply when filename not set
	fileIgnoreConfig.FileName = ""
	assert.False(t, fileIgnoreConfig.isEffective("filename"))
}

func TestScanFileIgnoreConfig_GetAllowedPatterns(t *testing.T) {
	fileIgnoreConfig := &ScanFileIgnoreConfig{
		FileName:        "some_filename",
		Checksums:       nil,
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
			&ScanFileIgnoreConfig{
				FileName:  "some_other_filename",
				Checksums: []string{"some_other_checksum"},
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
