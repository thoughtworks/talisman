package detector

import (
	"talisman/utility"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
