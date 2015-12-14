package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsingShasFromStdIn(t *testing.T) {
	file := tempFileWithContents("localRef localSha remoteRef remoteSha")
	defer os.Remove(file.Name())

	_, oldSha, _, newSha := readRefAndSha(file)

	assert.Equal(t, "localSha", oldSha, "oldSha did not equal 'localSha', got: %s", oldSha)
	assert.Equal(t, "remoteSha", newSha, "newSha did not equal 'remoteSha', got: %s", newSha)
}

// caller is responsible for closing the file
func tempFileWithContents(contents string) *os.File {
	file, err := ioutil.TempFile(os.TempDir(), "stdin")
	check(err)
	file.WriteString(contents)
	file.Seek(0, 0)
	return file
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
