package talismanrc

import (
	"fmt"
	"io"
	"regexp"
	"testing"

	"talisman/detector/severity"
	"talisman/gitrepo"

	logr "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func init() {
	logr.SetOutput(io.Discard)
}

func TestShouldFilterAllowedPatternsFromAddition(t *testing.T) {
	const hex string = "68656C6C6F20776F726C6421"
	const fileContent string = "Prefix content" + hex
	gitRepoAddition1 := testAdditionWithData("file1", []byte(fileContent))
	talismanrc := &TalismanRC{AllowedPatterns: []*Pattern{{regexp.MustCompile(hex)}}}

	fileContentFiltered := talismanrc.FilterAllowedPatternsFromAddition(gitRepoAddition1)

	assert.Equal(t, fileContentFiltered, "Prefix content")
}

func TestShouldFilterAllowedPatternsFromAdditionBasedOnFileConfig(t *testing.T) {
	const hexContent string = "68656C6C6F20776F726C6421"
	const fileContent string = "Prefix content" + hexContent
	gitRepoAddition1 := testAdditionWithData("file1", []byte(fileContent))
	gitRepoAddition2 := testAdditionWithData("file2", []byte(fileContent))
	talismanrc := createTalismanRCWithFileIgnores("file1", "somedetector", []string{hexContent})

	fileContentFiltered1 := talismanrc.FilterAllowedPatternsFromAddition(gitRepoAddition1)
	fileContentFiltered2 := talismanrc.FilterAllowedPatternsFromAddition(gitRepoAddition2)

	assert.Equal(t, fileContentFiltered1, "Prefix content")
	assert.Equal(t, fileContentFiltered2, fileContent)
}

func TestObeysCustomSeverityLevelsAndThreshold(t *testing.T) {
	talismanRCContents := []byte(`threshold: high
custom_severities:
- detector: Base64Content
  severity: low
`)
	talismanRC, _ := newPersistedRC(talismanRCContents)
	assert.Equal(t, talismanRC.Threshold, severity.High)
	assert.Equal(t, len(talismanRC.CustomSeverities), 1)
	assert.Equal(t, talismanRC.CustomSeverities, []CustomSeverityConfig{{Detector: "Base64Content", Severity: severity.Low}})
}

func TestDirectoryPatterns(t *testing.T) {
	assertAccepts("foo/", "", "bar", t)
	assertAccepts("foo/", "", "foo", t)
	assertDenies("foo/", "filename", "foo/bar", t)
	assertDenies("foo/", "filename", "foo/bar.txt", t)
	assertDenies("foo/", "filename", "foo/bar/baz.txt", t)
}

func TestIgnoreAdditionsByScope(t *testing.T) {
	testTable := map[string][]gitrepo.Addition{
		"node": {
			testAddition("yarn.lock"),
			testAddition("pnpm-lock.yaml"),
			testAddition("package-lock.json"),
			testAddition("node_modules/module1/foo.js")},
		"go": {
			testAddition("Gopkg.lock"),
			testAddition("makefile"),
			testAddition("go.mod"), testAddition("go.sum"),
			testAddition("Gopkg.toml"), testAddition("Gopkg.lock"),
			testAddition("glide.yaml"), testAddition("glide.lock"),
		},
		"images": {
			testAddition("img.jpeg"),
			testAddition("img.jpg"),
			testAddition("img.png"),
			testAddition("img.tiff"),
			testAddition("img.bmp"),
		},
		"bazel": {testAddition("bazelfile.bzl")},
		"terraform": {
			testAddition(".terraform.lock.hcl"),
			testAddition("foo/.terraform.lock.hcl"),
			testAddition("foo/bar/.terraform.lock.hcl"),
		},
		"php": {
			testAddition("composer.lock"),
		},
		"python": {
			testAddition("poetry.lock"),
			testAddition("Pipfile.lock"),
			testAddition("requirements.txt"),
		},
	}

	for scopeName, additions := range testTable {
		t.Run(fmt.Sprintf("should ignore files for %s scope", scopeName), func(t *testing.T) {
			talismanRCConfig := createTalismanRCWithScopeIgnores([]string{scopeName})
			filteredAdditions := talismanRCConfig.FilterAdditions(additions)
			for _, addition := range additions {
				assert.NotContains(t, filteredAdditions, addition, "Expected %s to be ignored", addition.Name)
			}
		})
	}
}

func TestIgnoringDetectors(t *testing.T) {
	assertDeniesDetector("foo", "someDetector", "foo", "someDetector", t)
	assertAcceptsDetector("foo", "someDetector", "foo", "someOtherDetector", t)
}

func TestAddingFileIgnores(t *testing.T) {
	fs := afero.NewMemMapFs()
	file, err := afero.TempFile(fs, "", DefaultRCFileName)
	assert.NoError(t, err)
	ignoreFile := file.Name()

	SetFs__(fs)
	SetRcFilename__(ignoreFile)

	ignoreConfig := FileIgnoreConfig{
		FileName: "Foo",
		Checksum: "SomeCheckSum"}
	t.Run("When .talismanrc doesn't exist yet", func(t *testing.T) {

		initialRCConfig, _ := Load()
		initialRCConfig.AddIgnores([]FileIgnoreConfig{ignoreConfig})
		newRCConfig, _ := Load()
		assert.Equal(t, 1, len(newRCConfig.FileIgnoreConfig))
		_ = fs.Remove(ignoreFile)
	})

	t.Run("When there already is a .talismanrc", func(t *testing.T) {
		existingContent := `
fileignoreconfig:
- filename: existing.pem
  checksum: 123444ddssa75333b25b6275f97680604add51b84eb8f4a3b9dcbbc652e6f27ac
scopeconfig: [scope: go]
allowed_patterns:
- this-is-okay
- key={listOfThings.id}
custom_patterns:
- this-isn't-okay
threshold: medium
custom_severities:
- detector: HexContent
  severity: low
experimental:
  base64EntropyThreshold: 4.7
version: 1.0
`
		err = afero.WriteFile(fs, ignoreFile, []byte(existingContent), 0666)
		assert.NoError(t, err)

		initialRCConfig, _ := Load()
		initialRCConfig.AddIgnores([]FileIgnoreConfig{ignoreConfig})
		newRCConfig, _ := Load()
		expectedTalismanRC := &TalismanRC{
			FileIgnoreConfig: []FileIgnoreConfig{
				ignoreConfig,
				{FileName: "existing.pem", Checksum: "123444ddssa75333b25b6275f97680604add51b84eb8f4a3b9dcbbc652e6f27ac"}},
			ScopeConfig: []ScopeConfig{{"go"}},
			AllowedPatterns: []*Pattern{
				{regexp.MustCompile("this-is-okay")},
				{regexp.MustCompile("key={listOfThings.id}")}},
			CustomPatterns: []PatternString{"this-isn't-okay"},
			Threshold:      severity.Medium,
			CustomSeverities: []CustomSeverityConfig{
				{Detector: "HexContent", Severity: severity.Low}},
			Experimental: ExperimentalConfig{Base64EntropyThreshold: 4.7},
			Version:      "1.0",
		}
		assert.Equal(t, expectedTalismanRC, newRCConfig)
	})
}

func assertDenies(line, ignoreDetector string, path string, t *testing.T) {
	assertDeniesDetector(line, ignoreDetector, path, "filename", t)
}

func assertDeniesDetector(line, ignoreDetector string, path string, detectorName string, t *testing.T) {
	assert.True(t, createTalismanRCWithFileIgnores(line, ignoreDetector, []string{}).Deny(testAddition(path), detectorName), "%s is expected to deny a file named %s.", line, path)
}

func assertAccepts(line, ignoreDetector string, path string, t *testing.T, detectorNames ...string) {
	assertAcceptsDetector(line, ignoreDetector, path, "someDetector", t)
}

func assertAcceptsDetector(line, ignoreDetector string, path string, detectorName string, t *testing.T) {
	assert.True(t, createTalismanRCWithFileIgnores(line, ignoreDetector, []string{}).Accept(testAddition(path), detectorName), "%s is expected to accept a file named %s.", line, path)
}

func testAddition(path string) gitrepo.Addition {
	return gitrepo.NewAddition(path, make([]byte, 0))
}

func testAdditionWithData(path string, content []byte) gitrepo.Addition {
	return gitrepo.NewAddition(path, content)
}

func createTalismanRCWithFileIgnores(filename string, detector string, allowedPatterns []string) *TalismanRC {
	fileIgnoreConfig := FileIgnoreConfig{}
	fileIgnoreConfig.FileName = filename
	if detector != "" {
		fileIgnoreConfig.IgnoreDetectors = []string{detector}
	}
	if len(allowedPatterns) != 0 {
		fileIgnoreConfig.AllowedPatterns = allowedPatterns
	}

	return &TalismanRC{FileIgnoreConfig: []FileIgnoreConfig{fileIgnoreConfig}}
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
		fileIgnoreConfigs := []FileIgnoreConfig{
			{
				FileName: "some_filename",
				Checksum: "some_checksum",
			},
		}
		expectedRC := `fileignoreconfig:
- filename: some_filename
  checksum: some_checksum
version: "1.0"
`
		str := SuggestRCFor(fileIgnoreConfigs)
		assert.Equal(t, expectedRC, str)
	})

	t.Run("should ignore invalid configs", func(t *testing.T) {
		fileIgnoreConfigs := []FileIgnoreConfig{
			{
				FileName: "some_filename",
				Checksum: "some_checksum",
			},
		}
		expectedRC := `fileignoreconfig:
- filename: some_filename
  checksum: some_checksum
version: "1.0"
`
		str := SuggestRCFor(fileIgnoreConfigs)
		assert.Equal(t, expectedRC, str)
	})
}
