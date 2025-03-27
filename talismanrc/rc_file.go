package talismanrc

import (
	"fmt"

	logr "github.com/sirupsen/logrus"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

const (
	//DefaultRCFileName represents the name of default file in which all the ignore patterns are configured in new version
	DefaultRCVersion         = "1.0"
	DefaultRCFileName string = ".talismanrc"
)

var (
	fs                = afero.NewOsFs()
	currentRCFileName = DefaultRCFileName
)

// Load creates a TalismanRC struct based on a .talismanrc file, if present
func Load() (*TalismanRC, error) {
	fileContents, err := afero.ReadFile(fs, currentRCFileName)
	if err != nil {
		// File does not exist or is not readable, proceed as if there is no .talismanrc
		fileContents = []byte{}
	}
	return talismanRCFromYaml(fileContents)
}

func talismanRCFromYaml(fileContents []byte) (*TalismanRC, error) {
	talismanRCFromFile := TalismanRC{}
	err := yaml.Unmarshal(fileContents, &talismanRCFromFile)
	if err != nil {
		logr.Errorf("Unable to parse .talismanrc : %v", err)
		fmt.Println(fmt.Errorf("\n\x1b[1m\x1b[31mUnable to parse .talismanrc %s. Please ensure it is following the right YAML structure\x1b[0m\x1b[0m", err))
		return &TalismanRC{}, err
	}
	if talismanRCFromFile.Version == "" {
		talismanRCFromFile.Version = DefaultRCVersion
	}
	return &talismanRCFromFile, nil
}

func SetFs__(_fs afero.Fs) {
	fs = _fs
}

func SetRcFilename__(rcFileName string) {
	currentRCFileName = rcFileName
}
