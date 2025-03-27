package talismanrc

import (
	"talisman/detector/severity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalsValidYaml(t *testing.T) {
	t.Run("talismanrc should not fail as long as the yaml structure is correct", func(t *testing.T) {
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

	t.Run("talismanrc.For(mode) should read multiple entries in rc file correctly", func(t *testing.T) {
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
}

func TestShouldIgnoreUnformattedFiles(t *testing.T) {
	for _, s := range []string{"#", "#monkey", "# this monkey likes bananas  "} {
		fileContents := []byte(s)
		talismanRC, err := newPersistedRC(fileContents)
		assert.Nil(t, err, "Should successfully unmarshal commented yaml")
		assert.Equal(t, &TalismanRC{Version: "1.0"}, talismanRC, "Expected commented line '%s' to result in an empty TalismanRC")
	}
}

func TestShouldConvertThresholdToValue(t *testing.T) {
	talismanRCContents := []byte("threshold: high")
	persistedTalismanrc, _ := newPersistedRC(talismanRCContents)
	assert.Equal(t, persistedTalismanrc.Threshold, severity.High)
}
