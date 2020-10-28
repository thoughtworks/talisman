package scanner

import (
	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"talisman/git_testing"
	"testing"
)

func init() {
	git_testing.Logger = logrus.WithField("Environment", "Debug")
	git_testing.Logger.Debug("Accetpance test started")
}

func Test_getBlobsFromChannel(t *testing.T) {
	ch := make(chan []string)
	go func() {
		ch <- []string{
			"100644 blob 351324aa7b3c66043e484c2f2c7b7f1842152f35	.gitignore",
			"100644 blob 8715df9907604c8ee8fc5e377821817f84f014fa	.pre-commit-hooks.yaml",
			"commitSha",
		}
	}()
	blobsInCommits := BlobsInCommits{commits: map[string][]string{}}
	getBlobsFromChannel(blobsInCommits, ch)

	commits := blobsInCommits.commits
	assert.Len(t, commits, 2)
	assert.Equal(t, []string{"commitSha"}, commits["351324aa7b3c66043e484c2f2c7b7f1842152f35	.gitignore"])
	assert.Equal(t, []string{"commitSha"}, commits["8715df9907604c8ee8fc5e377821817f84f014fa	.pre-commit-hooks.yaml"])
}
