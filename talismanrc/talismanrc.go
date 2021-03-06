package talismanrc

import (
	"github.com/spf13/afero"
	"os"
	"talisman/gitrepo"
)

const (
	//DefaultRCFileName represents the name of default file in which all the ignore patterns are configured in new version
	DefaultRCFileName string = ".talismanrc"
)

func SetFs(_fs afero.Fs) {
	fs = _fs
}

func SetRcFilename(rcFileName string) {
	currentRCFileName = rcFileName
}

func readRepoFile() func(string) ([]byte, error) {
	wd, _ := os.Getwd()
	repo := gitrepo.RepoLocatedAt(wd)
	return repo.ReadRepoFileOrNothing
}

func Get() *TalismanRC {
	return ReadConfigFromRCFile(readRepoFile())
}



