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

func ReadConfigFromRCFile(repoFileRead func(string) ([]byte, error)) *persistedRC {
	fileContents, err := repoFileRead(currentRCFileName)
	if err != nil {
		panic(err)
	}
	return NewTalismanRC(fileContents)
}

func NewTalismanRC(fileContents []byte) *persistedRC {
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

func ConfigFromFile() *persistedRC {
	return ReadConfigFromRCFile(readRepoFile())
}

func MakeWithFileIgnores(fileIgnoreConfigs []FileIgnoreConfig) *persistedRC {
	return &persistedRC{FileIgnoreConfig: fileIgnoreConfigs, Version: DefaultRCVersion}
}

func BuildIgnoreConfig(mode Mode, filepath, checksum string, detectors []string) IgnoreConfig {
	switch mode {
	case HookMode:
		return &FileIgnoreConfig{FileName: filepath, Checksum: checksum, IgnoreDetectors: detectors}
	case ScanMode:
		return &ScanFileIgnoreConfig{FileName: filepath, Checksums: []string{checksum}, IgnoreDetectors: detectors}
	default:
		return &FileIgnoreConfig{FileName: filepath, Checksum: checksum, IgnoreDetectors: detectors}
	}
}
