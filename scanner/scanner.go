package scanner

import (
	"os/exec"
	"strings"
	"talisman/git_repo"
)

func GetAdditions(blobDetails string) []git_repo.Addition {
	var additions []git_repo.Addition
	blobArray := strings.Split(blobDetails, "\n")
	for _, blob := range blobArray {
		objectDetails := strings.Split(blob, " ")
		objectHash := objectDetails[0]
		data := getData(objectHash)
		filePath := objectDetails[1]
		newAddition := git_repo.NewAddition(filePath, data)
		additions = append(additions, newAddition)
	}
	return additions
}

func getData(objectHash string) []byte {
	out, _ := exec.Command("git", "cat-file", "-p", objectHash).CombinedOutput()
	return out
}
