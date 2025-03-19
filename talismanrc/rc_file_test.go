package talismanrc

import (
	"talisman/detector/severity"
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultRepoFileReader = repoFileReader()

func TestLoadsValidYaml(t *testing.T) {
	var repoFileReader = func(string) ([]byte, error) {
		return []byte(`
---
fileignoreconfig:
- filename: testfile_1.yml
  checksum: file1_checksum

custom_patterns:
- 'pwd_[a-z]{8,20}'`), nil
	}
	t.Run("talismanrc should not fail as long as the yaml structure is correct", func(t *testing.T) {
		setRepoFileReader(repoFileReader)
		rc, _ := For(HookMode)
		assert.Equal(t, 1, len(rc.IgnoreConfigs))
		assert.Equal(t, 1, len(rc.CustomPatterns))
	})
	setRepoFileReader(defaultRepoFileReader)
}

func TestShouldIgnoreUnformattedFiles(t *testing.T) {
	for _, s := range []string{"#", "#monkey", "# this monkey likes bananas  "} {
		setRepoFileReader(func(string) ([]byte, error) {
			return []byte(s), nil
		})

		talismanRC, _ := For(HookMode)
		assert.True(t, talismanRC.AcceptsAll(), "Expected commented line '%s' to result in no ignore patterns", s)
		talismanRC, _ = For(ScanMode)
		assert.True(t, talismanRC.AcceptsAll(), "Expected commented line '%s' to result in no ignore patterns", s)
	}
	setRepoFileReader(defaultRepoFileReader)
}

func TestShouldConvertThresholdToValue(t *testing.T) {
	talismanRCContents := []byte("threshold: high")
	persistedTalismanrc, _ := newPersistedRC(talismanRCContents)
	assert.Equal(t, persistedTalismanrc.Threshold, severity.High)
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
		rc, _ := For(HookMode)
		assert.Equal(t, 3, len(rc.IgnoreConfigs))

		assert.Equal(t, rc.IgnoreConfigs[0].GetFileName(), "testfile_1.yml")
		assert.True(t, rc.IgnoreConfigs[0].ChecksumMatches("file1_checksum"))
		assert.Equal(t, rc.IgnoreConfigs[1].GetFileName(), "testfile_2.yml")
		assert.True(t, rc.IgnoreConfigs[1].ChecksumMatches("file2_checksum"))
		assert.Equal(t, rc.IgnoreConfigs[2].GetFileName(), "testfile_3.yml")
		assert.True(t, rc.IgnoreConfigs[2].ChecksumMatches("file3_checksum"))

		setRepoFileReader(defaultRepoFileReader)
	})

	t.Run("talismanrc.ForScan(ignoreHistory) should populate talismanrc for scan mode with ignore history", func(t *testing.T) {
		setRepoFileReader(repoFileReader)
		rc, _ := ForScan(true)

		assert.Equal(t, 3, len(rc.IgnoreConfigs))
		setRepoFileReader(defaultRepoFileReader)

	})

	t.Run("talismanrc.ForScan(ignoreHistory) should populate talismanrc for scan mode without ignore history", func(t *testing.T) {
		setRepoFileReader(repoFileReader)
		rc, _ := ForScan(false)

		assert.Equal(t, 0, len(rc.IgnoreConfigs))
		setRepoFileReader(defaultRepoFileReader)

	})
}
