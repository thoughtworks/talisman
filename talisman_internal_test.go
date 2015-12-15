package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsingShasFromStdIn(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "mockStdin")
	if err != nil {
		panic(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()
	file.WriteString("localRef localSha remoteRef remoteSha")
	file.Seek(0, 0)

	_, oldSha, _, newSha := readRefAndSha(file)
	assert.Equal(t, "localSha", oldSha, "oldSha did not equal 'localSha', got: %s", oldSha)
	assert.Equal(t, "remoteSha", newSha, "newSha did not equal 'remoteSha', got: %s", newSha)
}
