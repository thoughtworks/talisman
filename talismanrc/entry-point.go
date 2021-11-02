package talismanrc

import (
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

func ReadConfigFromRCFile(fileReader func(string) ([]byte, error)) *persistedRC {
	fileContents, err := fileReader(currentRCFileName)
	if err != nil {
		panic(err)
	}
	return newPersistedRC(fileContents)
}

func newPersistedRC(fileContents []byte) *persistedRC {
	talismanRCFromFile := persistedRC{}
	err := yaml.Unmarshal(fileContents, &talismanRCFromFile)
	if err != nil {
		logr.Errorf("Unable to parse .talismanrc : %v", err)
		return &persistedRC{}
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

func ConfigFromFile() *persistedRC {
	return ReadConfigFromRCFile(repoFileReader())
}

func MakeWithFileIgnores(fileIgnoreConfigs []FileIgnoreConfig) *persistedRC {
	return &persistedRC{FileIgnoreConfig: fileIgnoreConfigs, Version: DefaultRCVersion}
}

func BuildIgnoreConfig(mode Mode, filepath, checksum string, detectors []string) IgnoreConfig {
	var result IgnoreConfig
	switch mode {
	case HookMode:
		result = &FileIgnoreConfig{FileName: filepath, Checksum: checksum, IgnoreDetectors: detectors}
	case ScanMode:
		result = &FileIgnoreConfig{FileName: filepath, Checksum: checksum, IgnoreDetectors: detectors}
	}
	return result
}
