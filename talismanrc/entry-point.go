package talismanrc

import (
	"fmt"
	"regexp"
	"talisman/utility"

	logr "github.com/sirupsen/logrus"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

var (
	emptyStringPattern = regexp.MustCompile(`^\s*$`)
	fs                 = afero.NewOsFs()
	currentRCFileName  = DefaultRCFileName
)

func ReadConfigFromRCFile(fileReader func(string) ([]byte, error)) (*persistedRC, error) {
	fileContents, err := fileReader(currentRCFileName)
	if err != nil {
		panic(err)
	}
	return newPersistedRC(fileContents)
}

func newPersistedRC(fileContents []byte) (*persistedRC, error) {
	talismanRCFromFile := persistedRC{}
	err := yaml.Unmarshal(fileContents, &talismanRCFromFile)
	if err != nil {
		logr.Errorf("Unable to parse .talismanrc : %v", err)
		fmt.Println(fmt.Errorf("\n\x1b[1m\x1b[31mUnable to parse .talismanrc %s. Please ensure it is following the right YAML structure\x1b[0m\x1b[0m", err))
		return &persistedRC{}, err
	}
	if talismanRCFromFile.Version == "" {
		talismanRCFromFile.Version = DefaultRCVersion
	}
	return &talismanRCFromFile, nil
}

const (
	//DefaultRCFileName represents the name of default file in which all the ignore patterns are configured in new version
	DefaultRCVersion         = "1.0"
	DefaultRCFileName string = ".talismanrc"
)

func SetFs__(_fs afero.Fs) {
	fs = _fs
}

func SetRcFilename__(rcFileName string) {
	currentRCFileName = rcFileName
}

type RepoFileReader func(string) ([]byte, error)

var repoFileReader = func() RepoFileReader {

	return func(path string) ([]byte, error) {
		data, err := utility.SafeReadFile(path)
		if err != nil {
			return []byte{}, nil
		}
		return data, nil
	}
}

func setRepoFileReader(rfr RepoFileReader) {
	repoFileReader = func() RepoFileReader { return rfr }
}

func ConfigFromFile() (*persistedRC, error) {
	return ReadConfigFromRCFile(repoFileReader())
}

func MakeWithFileIgnores(fileIgnoreConfigs []FileIgnoreConfig) *persistedRC {
	return &persistedRC{FileIgnoreConfig: fileIgnoreConfigs, Version: DefaultRCVersion}
}

func BuildIgnoreConfig(filepath, checksum string, detectors []string) IgnoreConfig {
	return &FileIgnoreConfig{FileName: filepath, Checksum: checksum, IgnoreDetectors: detectors}
}
