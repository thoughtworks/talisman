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
	file, err := afero.TempFile(fs, "", DefaultRCFileName)
	assert.NoError(t, err, "Problem setting up test .talismanrc?")
	talismanRCFile := file.Name()
	SetFs__(fs)

	t.Run("Creates an empty TalismanRC if .talismanrc file doesn't exist", func(t *testing.T) {
		SetRcFilename__("not-a-file")
		emptyRC, err := Load()
		assert.NoError(t, err, "Should not error if there is a problem reading the file")
		assert.Equal(t, &TalismanRC{Version: DefaultRCVersion}, emptyRC)
	})

	t.Run("Loads all valid TalismanRC fields", func(t *testing.T) {
		SetRcFilename__(talismanRCFile)
		err = afero.WriteFile(fs, talismanRCFile, []byte(fullyConfiguredTalismanRC), 0666)
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

	SetRcFilename__(DefaultRCFileName)
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

		rc, err := newPersistedRC(fileContents)
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

		rc, _ := newPersistedRC(fileContent)
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
		persistedTalismanrc, _ := newPersistedRC(talismanRCContents)
		assert.Equal(t, persistedTalismanrc.Threshold, severity.High)
	})

	t.Run("Should read custom severities", func(t *testing.T) {
		talismanRCContents := []byte(`
custom_severities:
- detector: Base64Content
  severity: low
`)
		talismanRC, _ := newPersistedRC(talismanRCContents)
		assert.Equal(t, talismanRC.CustomSeverities, []CustomSeverityConfig{{Detector: "Base64Content", Severity: severity.Low}})
	})
}

func TestShouldIgnoreUnformattedFiles(t *testing.T) {
	for _, s := range []string{"#", "#monkey", "# this monkey likes bananas  "} {
		fileContents := []byte(s)
		talismanRC, err := newPersistedRC(fileContents)
		assert.Nil(t, err, "Should successfully unmarshal commented yaml")
		assert.Equal(t, &TalismanRC{Version: "1.0"}, talismanRC, "Expected commented line '%s' to result in an empty TalismanRC")
	}
}
