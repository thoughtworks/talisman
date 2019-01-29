package utility

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
)

//UniqueItems returns the array of strings containing unique items
func UniqueItems(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

//CollectiveSHA256Hash return collective sha256 hash of the passed paths
func CollectiveSHA256Hash(paths []string) string {
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

func hashByte(contentPtr *[]byte) string {
	contents := *contentPtr
	hasher := sha256.New()
	hasher.Write(contents)
	return hex.EncodeToString(hasher.Sum(nil))
}
