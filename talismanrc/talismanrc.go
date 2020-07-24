package talismanrc

import (
	logr "github.com/Sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"reflect"
	"regexp"
	"sort"

	"talisman/gitrepo"
)

const (
	//DefaultRCFileName represents the name of default file in which all the ignore patterns are configured in new version
	DefaultRCFileName string = ".talismanrc"
)

var (
	emptyStringPattern = regexp.MustCompile(`^\s*$`)
	fs                 = afero.NewOsFs()
	currentRCFileName  = DefaultRCFileName
)


type FileIgnoreConfig struct {
	FileName        string   `yaml:"filename"`
	Checksum        string   `yaml:"checksum"`
	IgnoreDetectors []string `yaml:"ignore_detectors,omitempty"`
	AllowedPatterns []string `yaml:"allowed_patterns,omitempty"`
}

type ScopeConfig struct {
	ScopeName string `yaml:"scope"`
}

type ExperimentalConfig struct {
	Base64EntropyThreshold float64 `yaml:"base64EntropyThreshold,omitempty"`
}

type PatternString string

type TalismanRC struct {
	FileIgnoreConfig []FileIgnoreConfig `yaml:"fileignoreconfig,omitempty"`
	ScopeConfig      []ScopeConfig      `yaml:"scopeconfig,omitempty"`
	CustomPatterns   []PatternString    `yaml:"custom_patterns,omitempty"`
	AllowedPatterns  []string           `yaml:"allowed_patterns,omitempty"`
	Experimental     ExperimentalConfig `yaml:"experimental,omitempty"`
}

func SetFs(_fs afero.Fs) {
	fs = _fs
}

func SetRcFilename(rcFileName string) {
	currentRCFileName = rcFileName
}

func Get() *TalismanRC {
	return ReadConfigFromRCFile(readRepoFile())
}

func (tRC *TalismanRC) IsEmpty() bool {
	return reflect.DeepEqual(TalismanRC{}, tRC)
}

func ReadConfigFromRCFile(repoFileRead func(string) ([]byte, error)) *TalismanRC {
	fileContents, error := repoFileRead(currentRCFileName)
	if error != nil {
		panic(error)
	}
	return NewTalismanRC(fileContents)
}

func readRepoFile() func(string) ([]byte, error) {
	wd, _ := os.Getwd()
	repo := gitrepo.RepoLocatedAt(wd)
	return repo.ReadRepoFileOrNothing
}

func NewTalismanRC(fileContents []byte) *TalismanRC {
	talismanRC := TalismanRC{}
	err := yaml.Unmarshal(fileContents, &talismanRC)
	if err != nil {
		log.Println("Unable to parse .talismanrc")
		log.Printf("error: %v", err)
		return &talismanRC
	}
	return &talismanRC
}

func (i FileIgnoreConfig) isEffective(detectorName string) bool {
	return !isEmptyString(i.FileName) &&
		contains(i.IgnoreDetectors, detectorName)
}

//AcceptsAll returns true if there are no rules specified
func (tRC *TalismanRC) AcceptsAll() bool {
	return len(tRC.effectiveRules("any-detector")) == 0
}

//Accept answers true if the Addition.Path is configured to be checked by the detectors
func (tRC *TalismanRC) Accept(addition gitrepo.Addition, detectorName string) bool {
	return !tRC.Deny(addition, detectorName)
}

func (tRC *TalismanRC) IgnoreAdditionsByScope(additions []gitrepo.Addition, scopeMap map[string][]string) []gitrepo.Addition {
	var applicableScopeFileNames []string
	if tRC.ScopeConfig != nil {
		for _, scope := range tRC.ScopeConfig {
			if len(scopeMap[scope.ScopeName]) > 0 {
				applicableScopeFileNames = append(applicableScopeFileNames, scopeMap[scope.ScopeName]...)
			}
		}
	}
	var result []gitrepo.Addition
	for _, addition := range additions {
		isFilePresentInScope := false
		for _, fileName := range applicableScopeFileNames {
			if addition.Matches(fileName) {
				isFilePresentInScope = true
			}
		}
		if !isFilePresentInScope {
			result = append(result, addition)
		}
	}
	return result
}

func (tRC *TalismanRC) AddFileIgnores(entriesToAdd []FileIgnoreConfig) {
	if len(entriesToAdd) > 0 {
		logr.Debugf("Adding entries: %v", entriesToAdd)
		talismanRCConfig := Get()
		talismanRCConfig.FileIgnoreConfig = combineFileIgnores(talismanRCConfig.FileIgnoreConfig, entriesToAdd)
		ignoreEntries, _ := yaml.Marshal(&talismanRCConfig)
		file, err := fs.OpenFile(currentRCFileName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("error opening %s: %s", currentRCFileName, err)
		}
		defer func() {
			err := file.Close()
			if err != nil {
				log.Printf("error closing %s: %s", currentRCFileName, err)
			}

		}()
		logr.Debugf("Writing talismanrc: %v", string(ignoreEntries))
		_, err = file.WriteString(string(ignoreEntries))
		if err != nil {
			log.Printf("error writing to %s: %s", currentRCFileName, err)
		}
	}
}

func combineFileIgnores(exsiting, incoming []FileIgnoreConfig) []FileIgnoreConfig {
	existingMap := make(map[string]FileIgnoreConfig)
	for _, fIC := range exsiting {
		existingMap[fIC.FileName] = fIC
	}
	for _, fIC := range incoming {
		existingMap[fIC.FileName] = fIC
	}
	result := make([]FileIgnoreConfig, len(existingMap))
	resultKeys := make([]string, len(existingMap))
	index := 0
	//sort keys in alphabetical order
	for k, _ := range existingMap {
		resultKeys[index] = k
		index++
	}
	sort.Strings(resultKeys)
	//add result entries based on sorted keys
	index = 0
	for _, k := range resultKeys {
		result[index] = existingMap[k]
		index++
	}
	return result
}

//Deny answers true if the Addition.Path is configured to be ignored and not checked by the detectors
func (tRC *TalismanRC) Deny(addition gitrepo.Addition, detectorName string) bool {
	result := false
	for _, pattern := range tRC.effectiveRules(detectorName) {
		result = result || addition.Matches(pattern)
	}
	return result
}

func (tRC *TalismanRC) effectiveRules(detectorName string) []string {
	var result []string
	for _, ignore := range tRC.FileIgnoreConfig {
		if ignore.isEffective(detectorName) {
			result = append(result, ignore.FileName)
		}
	}
	return result
}

func isEmptyString(str string) bool {
	return emptyStringPattern.MatchString(str)
}
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
