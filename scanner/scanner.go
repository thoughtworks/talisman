package scanner

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
	"talisman/gitrepo"
)

// BlobsInCommits is a map of blob and list of the commits the blobs is present in.
type BlobsInCommits struct {
	commits map[string][]string
}

// GetAdditionsInCommitRange will get all the additions from "afterCommitNumber"th commit to "afterCommitNumber+numberOfCommits"th commit
func GetAdditionsInCommitRange(afterCommitNumber uint64, numberOfCommits uint64) []gitrepo.Addition {
	blobsInCommits := getBlobsInCommitRange(afterCommitNumber, numberOfCommits)
	var additions []gitrepo.Addition
	for blob := range blobsInCommits.commits {
		objectDetails := strings.Split(blob, "\t")
		objectHash := objectDetails[0]
		data := getData(objectHash)
		filePath := objectDetails[1]
		newAddition := gitrepo.NewScannerAddition(filePath, blobsInCommits.commits[blob], data)
		additions = append(additions, newAddition)
	}
	return additions
}

func getBlobsInCommitRange(afterCommitNumber uint64, numberOfCommits uint64) BlobsInCommits {
	commits := getAllCommitsInRange(afterCommitNumber, numberOfCommits)
	blobsInCommits := newBlobsInCommit()
	result := make(chan []string, len(commits))
	for _, commit := range commits {
		go putBlobsInChannel(commit, result)
	}
	for i := 1; i < len(commits); i++ {
		getBlobsFromChannel(blobsInCommits, result)
	}
	return blobsInCommits
}

func putBlobsInChannel(commit string, result chan []string) {
	if commit != "" {
		blobDetailsBytes, _ := exec.Command("git", "ls-tree", "-r", commit).CombinedOutput()
		blobDetailsList := strings.Split(string(blobDetailsBytes), "\n")
		blobDetailsList = append(blobDetailsList, commit)
		result <- blobDetailsList
	}
}

func getBlobsFromChannel(blobsInCommits BlobsInCommits, result chan []string) {
	blobs := <-result
	commit := blobs[len(blobs)-1]
	for _, blob := range blobs[:len(blobs)] {
		if blob != "" && blob != commit {
			blobDetailsString := strings.Split(blob, " ")
			blobDetails := strings.Split(blobDetailsString[2], "	")
			blobHash := blobDetails[0] + "\t" + blobDetails[1]
			blobsInCommits.commits[blobHash] = append(blobsInCommits.commits[blobHash], commit)
		}
	}
}

func getAllCommitsInRange(afterCommitNumber uint64, numberOfCommits uint64) []string {
	n := strconv.FormatUint(numberOfCommits, 10)
	skip := strconv.FormatUint(afterCommitNumber, 10)
	out, err := exec.Command("git", "log", "--all", "-"+n, "--skip="+skip, "--pretty=%H").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(out), "\n")
}

func getData(objectHash string) []byte {
	out, _ := exec.Command("git", "cat-file", "-p", objectHash).CombinedOutput()
	return out
}

func newBlobsInCommit() BlobsInCommits {
	commits := make(map[string][]string)
	return BlobsInCommits{commits: commits}
}
