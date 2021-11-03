package utility

import (
	"crypto/sha256"
	"encoding/hex"
	"talisman/gitrepo"
)

type SHA256Hasher interface {
	CollectiveSHA256Hash(paths []string) string
}

type DefaultSHA256Hasher struct{}

//CollectiveSHA256Hash return collective sha256 hash of the passed paths
func (DefaultSHA256Hasher) CollectiveSHA256Hash(paths []string) string {
	return collectiveSHA256Hash(paths, SafeReadFile)
}

type GitHeadFileSHA256Hasher struct {
	root string
}

type GitFileSHA256Hasher struct {
	root string
}

func NewGitHeadFileSHA256Hasher(root string) GitHeadFileSHA256Hasher {
	return GitHeadFileSHA256Hasher{root}
}

func NewGitFileSHA256Hasher(root string) GitFileSHA256Hasher {
	return GitFileSHA256Hasher{root}
}

func (g GitHeadFileSHA256Hasher) CollectiveSHA256Hash(paths []string) string {
	return collectiveSHA256Hash(paths, gitrepo.NewCommittedRepoFileReader(g.root))
}

func (g GitFileSHA256Hasher) CollectiveSHA256Hash(paths []string) string {
	return collectiveSHA256Hash(paths, gitrepo.NewRepoFileReader(g.root))
}

func hashByte(contentPtr *[]byte) string {
	contents := *contentPtr
	hasher := sha256.New()
	hasher.Write(contents)
	return hex.EncodeToString(hasher.Sum(nil))
}

func collectiveSHA256Hash(paths []string, FileReader func(string) ([]byte, error)) string {
	var finHash = ""
	for _, path := range paths {
		sbyte := []byte(finHash)
		concatBytes := hashByte(&sbyte)
		nameByte := []byte(path)
		nameHash := hashByte(&nameByte)
		fileBytes, _ := FileReader(path)
		fileHash := hashByte(&fileBytes)
		finHash = concatBytes + fileHash + nameHash
	}
	c := []byte(finHash)
	m := hashByte(&c)
	return m
}
