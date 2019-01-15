package detector

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"talisman/git_repo"
)

type ChecksumCalculator struct{}

//NewChecksumCalculator returns new instance of the CheckSumDetector
func NewChecksumCalculator() *ChecksumCalculator {
	cs := ChecksumCalculator{}
	return &cs
}

//FilterIgnoresBasedOnChecksums filters the file ignores from the TalismanRCIgnore which doesn't have any checksum value or having mismatched checksum value from the .talsimanrc
func (cs *ChecksumCalculator) FilterIgnoresBasedOnChecksums(additions []git_repo.Addition, ignoreConfig TalismanRCIgnore) TalismanRCIgnore {
	fileIgnores := []FileIgnoreConfig{}
	for _, ignore := range ignoreConfig.FileIgnoreConfig {
		var patternpaths []string
		for _, addition := range additions {
			if addition.Matches(ignore.FileName) {
				patternpaths = append(patternpaths, string(addition.Path))
			}
		}
		// Calculate current checksum
		currentChecksum := CalculateCollectiveHash(patternpaths)
		// Compare with previous checksum from FileIgnoreConfig
		if ignore.Checksum == currentChecksum {
			fileIgnores = append(fileIgnores, ignore)
		}
	}
	rc := TalismanRCIgnore{}
	rc.FileIgnoreConfig = fileIgnores
	return rc
}

func hashByte(contentPtr *[]byte) string {
	contents := *contentPtr
	hasher := sha256.New()
	hasher.Write(contents)
	return hex.EncodeToString(hasher.Sum(nil))
}

//CalculateCollectiveHash returns collective hash of the paths passed as argument
func CalculateCollectiveHash(paths []string) string {
	var finHash = ""
	for _, path := range paths {
		sbyte := []byte(finHash)
		concatBytes := hashByte(&sbyte)
		nameByte := []byte(path)
		nameHash := hashByte(&nameByte)
		fileBytes, _ := ioutil.ReadFile(path)
		fileHash := hashByte(&fileBytes)
		finHash = concatBytes + fileHash + nameHash
	}
	c := []byte(finHash)
	m := hashByte(&c)
	return m
}
