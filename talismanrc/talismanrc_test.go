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

const fullyConfiguredTalismanRC = `
---
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

func init() {
	logr.SetOutput(io.Discard)
}

func TestShouldFilterAllowedPatternsFromAddition(t *testing.T) {
	const hex string = "68656C6C6F20776F726C6421"
	const fileContent string = "Prefix content" + hex
	gitRepoAddition1 := testAdditionWithData("file1", []byte(fileContent))
	talismanrc := &TalismanRC{AllowedPatterns: []*Pattern{{regexp.MustCompile(hex)}}}

	fileContentFiltered := talismanrc.RemoveAllowedPatterns(gitRepoAddition1)

	assert.Equal(t, fileContentFiltered, "Prefix content")
}

func TestShouldFilterAllowedPatternsFromAdditionBasedOnFileConfig(t *testing.T) {
	const hexContent string = "68656C6C6F20776F726C6421"
	const fileContent string = "Prefix content" + hexContent
	gitRepoAddition1 := testAdditionWithData("file1", []byte(fileContent))
	gitRepoAddition2 := testAdditionWithData("file2", []byte(fileContent))
	talismanrc := createTalismanRCWithFileIgnores("file1", "somedetector", []string{hexContent})

	fileContentFiltered1 := talismanrc.RemoveAllowedPatterns(gitRepoAddition1)
	fileContentFiltered2 := talismanrc.RemoveAllowedPatterns(gitRepoAddition2)

	assert.Equal(t, fileContentFiltered1, "Prefix content")
	assert.Equal(t, fileContentFiltered2, fileContent)
}

func TestShouldFilterAllowedPatternsFromAdditionBasedOnFileConfigWithWildcards(t *testing.T) {
	const hexContent string = "68656C6C6F20776F726C6421"
	const fileContent string = "Prefix content" + hexContent
	gitRepoAddition1 := testAdditionWithData("foo/file1.yml", []byte(fileContent))
	gitRepoAddition2 := testAdditionWithData("foo/file2.yml", []byte(fileContent))
	talismanrc := createTalismanRCWithFileIgnores("foo/*.yml", "somedetector", []string{hexContent})

	fileContentFiltered1 := talismanrc.RemoveAllowedPatterns(gitRepoAddition1)
	fileContentFiltered2 := talismanrc.RemoveAllowedPatterns(gitRepoAddition2)

	assert.Equal(t, fileContentFiltered1, "Prefix content")
	assert.Equal(t, fileContentFiltered2, "Prefix content")
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
		},
		"go": {
			testAddition("Gopkg.lock"),
			testAddition("makefile"),
			testAddition("go.mod"),
			testAddition("go.sum"),
			testAddition("submodule/go.mod"),
			testAddition("submodule/go.sum"),
			testAddition("Gopkg.toml"),
			testAddition("Gopkg.lock"),
			testAddition("glide.yaml"),
			testAddition("glide.lock"),
		},
		"images": {
			testAddition("img.jpeg"),
			testAddition("img.jpg"),
			testAddition("img.png"),
			testAddition("img.tiff"),
			testAddition("img.bmp"),
		},
		"bazel": {
			testAddition("bazelfile.bzl"),
		},
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
			filteredAdditions := talismanRCConfig.RemoveScopedFiles(additions)
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
	SetFs__(fs)

	ignoreConfig := FileIgnoreConfig{
		FileName: "Foo",
		Checksum: "SomeCheckSum"}
	t.Run("When .talismanrc is empty", func(t *testing.T) {
		initialRCConfig, _ := Load()
		initialRCConfig.AddIgnores([]FileIgnoreConfig{ignoreConfig})
		newRCConfig, _ := Load()
		assert.Equal(t, 1, len(newRCConfig.FileIgnoreConfig))
		_ = fs.Remove(RCFileName)
	})

	t.Run("When .talismanrc has lots of configurations", func(t *testing.T) {
		err := afero.WriteFile(fs, RCFileName, []byte(fullyConfiguredTalismanRC), 0666)
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
