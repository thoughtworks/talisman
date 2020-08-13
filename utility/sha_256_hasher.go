package utility

import (
	"crypto/sha256"
	"encoding/hex"
)

type SHA256Hasher interface {
	CollectiveSHA256Hash(paths []string) string
}

type DefaultSHA256Hasher struct {}

//CollectiveSHA256Hash return collective sha256 hash of the passed paths
func (DefaultSHA256Hasher) CollectiveSHA256Hash(paths []string) string {
	var finHash = ""
	for _, path := range paths {
		sbyte := []byte(finHash)
		concatBytes := hashByte(&sbyte)
		nameByte := []byte(path)
		nameHash := hashByte(&nameByte)
		fileBytes, _ := SafeReadFile(path)
		fileHash := hashByte(&fileBytes)
		finHash = concatBytes + fileHash + nameHash
	}
	c := []byte(finHash)
	m := hashByte(&c)
	return m
}

func hashByte(contentPtr *[]byte) string {
	contents := *contentPtr
	hasher := sha256.New()
	hasher.Write(contents)
	return hex.EncodeToString(hasher.Sum(nil))
}
