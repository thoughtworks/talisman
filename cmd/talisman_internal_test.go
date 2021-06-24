package main

import (
	"github.com/spf13/afero"
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

func Test_validateGitExecutable(t *testing.T) {
	t.Run("given operating systems is windows", func(t *testing.T) {

		operatingSystem := "windows"
		os.Setenv("PATHEXT", ".COM;.EXE;.BAT;.CMD;.VBS;.VBE;.JS;.JSE;.WSF;.WSH;.MSC")

		t.Run("should return error if git executable exists in current directory", func(t *testing.T) {
			fs := afero.NewMemMapFs()
			gitExecutable := "git.exe"
			afero.WriteFile(fs, gitExecutable, []byte("git executable"), 0700)
			err := validateGitExecutable(fs, operatingSystem)
			assert.EqualError(t, err, "not allowed to have git executable located in repository: git.exe")
		})

		t.Run("should return nil if git executable does not exist in current directory", func(t *testing.T) {
			err := validateGitExecutable(afero.NewMemMapFs(), operatingSystem)
			assert.NoError(t, err)
		})

	})

	t.Run("given operating systems is linux", func(t *testing.T) {

		operatingSystem := "linux"

		t.Run("should return nil if git executable exists in current directory", func(t *testing.T) {
			fs := afero.NewMemMapFs()
			gitExecutable := "git.exe"
			afero.WriteFile(fs, gitExecutable, []byte("git executable"), 0700)
			err := validateGitExecutable(fs, operatingSystem)
			assert.NoError(t, err)
		})

		t.Run("should return nil if git executable does not exist in current directory", func(t *testing.T) {
			err := validateGitExecutable(afero.NewMemMapFs(), operatingSystem)
			assert.NoError(t, err)
		})

	})

}
