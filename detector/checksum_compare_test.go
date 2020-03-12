package detector

import (
	"talisman/gitrepo"
	"talisman/talismanrc"
	"talisman/utility"
	"testing"

	"github.com/stretchr/testify/assert"
)

var talismanRCWithInCorrectChecksum = `
fileignoreconfig:
- filename : 'some_file.pem'
  checksum : '87139cc4d975333b25b6275f97680604add51b84eb8f4a3b9dcbbc652e6f27ac'
  ignore_detectors : []
`

var talismanRCWithCorrectChecksum = `
fileignoreconfig:
- filename : 'some_file.pem'
  checksum : '87139cc4d975333b25b6275f97680604add51b84eb8f4a3b9dcbbc652e6f27ac'
  ignore_detectors : []
- filename : 'test/some_file.pem'
  checksum : '25bd31a28bf9d4e06327f1c4a5cab2260574ae508803f66adcc393350e994866'
  ignore_detectors : []
`

var talismanRCWithOneCorrectChecksum = `
fileignoreconfig:
- filename : 'some_file.pem'
  checksum : '87139cc4d975333b25b6275f97680604add51b84eb8f4a3b9dcbbc652e6f27ac'
  ignore_detectors : []
- filename : 'test/some1_file.pem'
  checksum : '87139cc4d975333b25b6275f97680604add51b84eb8f4a3b9dcbbc652e6f27ac'
  ignore_detectors : []
`

func TestShouldConsiderBothFilesForDetection(t *testing.T) {
	rc := talismanrc.NewTalismanRC([]byte(talismanRCWithInCorrectChecksum))
	addition1 := gitrepo.NewAddition("some_file.pem", make([]byte, 0))
	addition2 := gitrepo.NewAddition("test/some_file.pem", make([]byte, 0))
	cc := NewChecksumCompare([]gitrepo.Addition{addition1, addition2}, rc)

	filteredRC := cc.FilterIgnoresBasedOnChecksums()

	assert.Len(t, filteredRC.FileIgnoreConfig, 0, "Should return empty ignores and detectors should scan both files")
}

func TestShouldNotConsiderBothFilesForDetection(t *testing.T) {
	addition1 := gitrepo.NewAddition("some_file.pem", make([]byte, 0))
	addition2 := gitrepo.NewAddition("test/some_file.pem", make([]byte, 0))
	rc := talismanrc.NewTalismanRC([]byte(talismanRCWithCorrectChecksum))
	cc := NewChecksumCompare([]gitrepo.Addition{addition1, addition2}, rc)

	filteredRC := cc.FilterIgnoresBasedOnChecksums()

	assert.Len(t, filteredRC.FileIgnoreConfig, 1, "Should return 1 ignore/s (some_file.pem) which detectors should honor")
}

func TestShouldConsiderOneFileForDetection(t *testing.T) {
	addition1 := gitrepo.NewAddition("some_file.pem", make([]byte, 0))
	addition2 := gitrepo.NewAddition("test/some1_file.pem", make([]byte, 0))
	rc := talismanrc.NewTalismanRC([]byte(talismanRCWithOneCorrectChecksum))
	cc := NewChecksumCompare([]gitrepo.Addition{addition1, addition2}, rc)

	filteredRC := cc.FilterIgnoresBasedOnChecksums()

	assert.Len(t, filteredRC.FileIgnoreConfig, 1, "Should return 1 ignore and detectors should scan that file")
}

func TestShouldConsiderBothFilesForDetectionIfTalismanRCIsEmpty(t *testing.T) {
	addition1 := gitrepo.NewAddition("some_file.pem", make([]byte, 0))
	addition2 := gitrepo.NewAddition("test/some_file.pem", make([]byte, 0))
	rc := talismanrc.NewTalismanRC([]byte{})
	cc := NewChecksumCompare([]gitrepo.Addition{addition1, addition2}, rc)

	filteredRC := cc.FilterIgnoresBasedOnChecksums()

	assert.Len(t, filteredRC.FileIgnoreConfig, 0, "Should return empty ignores and detectors should scan both files")
}

func TestShouldReturnCorrectFileHash(t *testing.T) {
	checksumSomeFile := utility.CollectiveSHA256Hash([]string{"some_file.pem"})
	checksumTestSomeFile := utility.CollectiveSHA256Hash([]string{"test/some_file.pem"})
	assert.Equal(t, checksumSomeFile, "87139cc4d975333b25b6275f97680604add51b84eb8f4a3b9dcbbc652e6f27ac", "Should be equal to some_file.pem hash value")
	assert.Equal(t, checksumTestSomeFile, "25bd31a28bf9d4e06327f1c4a5cab2260574ae508803f66adcc393350e994866", "Should be equal to test/some_file.pem hash value")
}

func TestShouldReturnEmptyFileHashWhenNoPathsPassed(t *testing.T) {
	checksum := utility.CollectiveSHA256Hash([]string{})
	assert.Equal(t, checksum, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", "Should be equal to empty hash value when no paths passed")
}
