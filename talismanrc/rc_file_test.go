package talismanrc

import (
	"regexp"
	"talisman/detector/severity"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestLoadingFromFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	SetFs__(fs)

	t.Run("Creates an empty TalismanRC if .talismanrc file doesn't exist", func(t *testing.T) {
		talismanRCFileExists, _ := afero.Exists(fs, RCFileName)
		assert.False(t, talismanRCFileExists, ".talismanrc file should NOT exist for this test!")
		emptyRC, err := Load()
		assert.NoError(t, err, "Should not error if there is a problem reading the file")
		assert.Equal(t, &TalismanRC{Version: DefaultRCVersion}, emptyRC)
	})

	t.Run("Loads all valid TalismanRC fields", func(t *testing.T) {
		err := afero.WriteFile(fs, RCFileName, []byte(fullyConfiguredTalismanRC), 0666)
		assert.NoError(t, err, "Problem setting up test .talismanrc?")

		talismanRCFromFile, _ := Load()
		expectedTalismanRC := &TalismanRC{
			FileIgnoreConfig: []FileIgnoreConfig{
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
		assert.Equal(t, expectedTalismanRC, talismanRCFromFile)
	})
}

func TestWritingToFile(t *testing.T) {
	tRC := &TalismanRC{Version: DefaultRCVersion}
	fs := afero.NewMemMapFs()
	SetFs__(fs)

	t.Run("When there is no .talismanrc file", func(t *testing.T) {
		talismanRCFileExists, _ := afero.Exists(fs, RCFileName)
		assert.False(t, talismanRCFileExists, "Problem setting up tests")
		tRC.saveToFile()
		talismanRCFileExists, _ = afero.Exists(fs, RCFileName)
		assert.True(t, talismanRCFileExists, "Should have created a new .talismanrc file")
		fileContents, _ := afero.ReadFile(fs, RCFileName)
		assert.Equal(t, "version: \"1.0\"\n", string(fileContents))
	})

	t.Run("When there already is a .talismanrc file", func(t *testing.T) {
		err := afero.WriteFile(fs, RCFileName, []byte("Some existing content to overwrite"), 0666)
		assert.NoError(t, err, "Problem setting up tests")
		tRC.saveToFile()
		talismanRCFileExists, _ := afero.Exists(fs, RCFileName)
		assert.True(t, talismanRCFileExists, "Should have created a new .talismanrc file")
		fileContents, _ := afero.ReadFile(fs, RCFileName)
		assert.Equal(t, "version: \"1.0\"\n", string(fileContents))
	})
}

func TestUnmarshalsValidYaml(t *testing.T) {
	t.Run("Should not fail as long as the yaml structure is correct", func(t *testing.T) {
		fileContents := []byte(`
---
fileignoreconfig:
- filename: testfile_1.yml
  checksum: file1_checksum
custom_patterns:
- 'pwd_[a-z]{8,20}'`)

		rc, err := talismanRCFromYaml(fileContents)
		assert.Nil(t, err, "Should successfully unmarshal valid yaml")
		assert.Equal(t, 1, len(rc.FileIgnoreConfig))
		assert.Equal(t, 1, len(rc.CustomPatterns))
	})

	t.Run("Should read multiple entries in rc file correctly", func(t *testing.T) {
		fileContent := []byte(`
fileignoreconfig:
- filename: testfile_1.yml
  checksum: file1_checksum
- filename: testfile_2.yml
  checksum: file2_checksum
- filename: testfile_3.yml
  checksum: file3_checksum`)

		rc, _ := talismanRCFromYaml(fileContent)
		assert.Equal(t, 3, len(rc.FileIgnoreConfig))

		assert.Equal(t, rc.FileIgnoreConfig[0].GetFileName(), "testfile_1.yml")
		assert.True(t, rc.FileIgnoreConfig[0].ChecksumMatches("file1_checksum"))
		assert.Equal(t, rc.FileIgnoreConfig[1].GetFileName(), "testfile_2.yml")
		assert.True(t, rc.FileIgnoreConfig[1].ChecksumMatches("file2_checksum"))
		assert.Equal(t, rc.FileIgnoreConfig[2].GetFileName(), "testfile_3.yml")
		assert.True(t, rc.FileIgnoreConfig[2].ChecksumMatches("file3_checksum"))
	})

	t.Run("Should read severity level", func(t *testing.T) {
		talismanRCContents := []byte("threshold: high")
		persistedTalismanrc, _ := talismanRCFromYaml(talismanRCContents)
		assert.Equal(t, persistedTalismanrc.Threshold, severity.High)
	})

	t.Run("Should read custom severities", func(t *testing.T) {
		talismanRCContents := []byte(`
custom_severities:
- detector: Base64Content
  severity: low
`)
		talismanRC, _ := talismanRCFromYaml(talismanRCContents)
		assert.Equal(t, talismanRC.CustomSeverities, []CustomSeverityConfig{{Detector: "Base64Content", Severity: severity.Low}})
	})
}

func TestShouldIgnoreUnformattedFiles(t *testing.T) {
	for _, s := range []string{"#", "#monkey", "# this monkey likes bananas  "} {
		fileContents := []byte(s)
		talismanRC, err := talismanRCFromYaml(fileContents)
		assert.Nil(t, err, "Should successfully unmarshal commented yaml")
		assert.Equal(t, &TalismanRC{Version: "1.0"}, talismanRC, "Expected commented line '%s' to result in an empty TalismanRC")
	}
}
