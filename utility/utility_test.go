package utility

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestShouldReadNormalFileCorrectly(t *testing.T) {
	tempFile, _ := ioutil.TempFile(os.TempDir(), "somefile")
	dataToBeWrittenInFile := []byte{0, 1, 2, 3}
	tempFile.Write(dataToBeWrittenInFile)
	tempFile.Close()

	readDataFromFileUsingIoutilDotReadFile, _ := ioutil.ReadFile(tempFile.Name())
	readDataFromFileUsingSafeFileRead, _ := SafeReadFile(tempFile.Name())
	os.Remove(tempFile.Name())

	assert.Equal(t, readDataFromFileUsingIoutilDotReadFile, dataToBeWrittenInFile)
	assert.Equal(t, readDataFromFileUsingSafeFileRead, dataToBeWrittenInFile)
}

func TestShouldNotReadSymbolicLinkTargetFile(t *testing.T) {
	tempFile, _ := ioutil.TempFile(os.TempDir(), "somefile")
	dataToBeWrittenInFile := []byte{0, 1, 2, 3}
	tempFile.Write(dataToBeWrittenInFile)
	tempFile.Close()
	symlinkFileName := tempFile.Name() + "symlink"
	os.Symlink(tempFile.Name(), symlinkFileName)

	readDataFromSymlinkUsingIoutilDotReadFile, _ := ioutil.ReadFile(symlinkFileName)
	readDataFromSymlinkUsingSafeFileRead, _ := SafeReadFile(symlinkFileName)
	os.Remove(symlinkFileName)
	os.Remove(tempFile.Name())

	assert.Equal(t, readDataFromSymlinkUsingIoutilDotReadFile, dataToBeWrittenInFile)
	assert.Equal(t, readDataFromSymlinkUsingSafeFileRead, []byte{})
}
