package scanner

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"log"
	"os/exec"
	"strings"
	"talisman/gitrepo"
)

// BlobsInCommits is a map of blob and list of the commits the blobs is present in.
type BlobsInCommits struct {
	commits map[string][]string
}

// GetAdditions will get all the additions for entire git history
func GetAdditions() []gitrepo.Addition {
	blobsInCommits := getBlobsInCommit()
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

func getBlobsInCommit() BlobsInCommits {
	commits := getAllCommits()
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

func getAllCommits() []string {
	var commitsHash []string
	r, err := git.PlainOpen("")
	if err != nil {
		log.Fatal(err)
	}
	cIter, _ := r.Log(&git.LogOptions{All: true})
	cIter.ForEach(func(c *object.Commit) error {
		commitsHash = append(commitsHash, c.Hash.String())
		return nil
	})
	return commitsHash
}

func getData(objectHash string) []byte {
	out, _ := exec.Command("git", "cat-file", "-p", objectHash).CombinedOutput()
	return out
}

func newBlobsInCommit() BlobsInCommits {
	commits := make(map[string][]string)
	return BlobsInCommits{commits: commits}
}
