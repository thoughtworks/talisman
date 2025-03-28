package talismanrc

import (
	"fmt"

	logr "github.com/sirupsen/logrus"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

const (
	// RCFileName represents the name of default file in which all the ignore patterns are configured in new version
	RCFileName       = ".talismanrc"
	DefaultRCVersion = "1.0"
)

var (
	fs = afero.NewOsFs()
)

// Load creates a TalismanRC struct based on a .talismanrc file, if present
func Load() (*TalismanRC, error) {
	fileContents, err := afero.ReadFile(fs, RCFileName)
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

func (tRC *TalismanRC) saveToFile() {
	ignoreEntries, _ := yaml.Marshal(&tRC)
	err := afero.WriteFile(fs, RCFileName, ignoreEntries, 0644)
	if err != nil {
		logr.Errorf("error writing to %s: %s", RCFileName, err)
	}
}

func SetFs__(_fs afero.Fs) {
	fs = _fs
}
