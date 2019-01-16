package scanner

import (
	"os/exec"
	"strings"
	"talisman/git_repo"
)

// GetAdditions will get all the additions for entire git history
func GetAdditions(blobDetails string) []git_repo.Addition {
	blobArray := strings.Split(blobDetails, "commit")
	commitsInBlob := make(map[string][]string)
	for _, commit := range blobArray {
		if commit != "" {
			objects := strings.Split(commit, " ")
			commitHash := objects[1]
			blobs := objects[2]
			for _, blob := range strings.Split(blobs, "\n") {
				commitsInBlob[blob] = append(commitsInBlob[blob], commitHash)
			}
		}
	}

	var additions []git_repo.Addition
	for blob := range commitsInBlob {
		objectDetails := strings.Split(blob, "\t")
		objectHash := objectDetails[0]
		data := getData(objectHash)
		filePath := objectDetails[1]
		newAddition := git_repo.NewScannerAddition(filePath, commitsInBlob[blob], data)
		additions = append(additions, newAddition)
	}
	return additions
}

func getData(objectHash string) []byte {
	out, _ := exec.Command("git", "cat-file", "-p", objectHash).CombinedOutput()
	return out
}
