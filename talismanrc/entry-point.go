package talismanrc

import (
	logr "github.com/Sirupsen/logrus"
	"os"
	"regexp"
	"talisman/gitrepo"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

var (
	emptyStringPattern = regexp.MustCompile(`^\s*$`)
	fs                 = afero.NewOsFs()
	currentRCFileName  = DefaultRCFileName
)


func ReadConfigFromRCFile(repoFileRead func(string) ([]byte, error)) *TalismanRC {
	fileContents, error := repoFileRead(currentRCFileName)
	if error != nil {
		panic(error)
	}
	return NewTalismanRC(fileContents)
}

func NewTalismanRC(fileContents []byte) *TalismanRC {
	talismanRCFromFile := TalismanRC{}
	err := yaml.Unmarshal(fileContents, &talismanRCFromFile)
	if err != nil {
		logr.Errorf("Unable to parse .talismanrc : %v",err)
		return &TalismanRC{}
	}
	if talismanRCFromFile.Version == "" {
		talismanRCFromFile.Version = DefaultRCVersion
	}
	return &talismanRCFromFile
}

const (
	//DefaultRCFileName represents the name of default file in which all the ignore patterns are configured in new version
	DefaultRCVersion         = "1.0"
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

func MakeWithFileIgnores(fileIgnoreConfigs []FileIgnoreConfig) *TalismanRC {
	return &TalismanRC{FileIgnoreConfig: fileIgnoreConfigs, Version: DefaultRCVersion}
}
