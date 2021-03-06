package talismanrc

import (
	logr "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"reflect"
	"sort"
	"talisman/detector/severity"

	"talisman/gitrepo"
)

type TalismanRC struct {
	FileIgnoreConfig []FileIgnoreConfig     `yaml:"fileignoreconfig,omitempty"`
	ScopeConfig      []ScopeConfig          `yaml:"scopeconfig,omitempty"`
	CustomPatterns   []PatternString        `yaml:"custom_patterns,omitempty"`
	CustomSeverities []CustomSeverityConfig `yaml:"custom_severities,omitempty"`
	AllowedPatterns  []string               `yaml:"allowed_patterns,omitempty"`
	Experimental     ExperimentalConfig     `yaml:"experimental,omitempty"`
	Threshold        severity.Severity      `default:"1" yaml:"threshold,omitempty"`
	ScanConfig       struct {
		FileIgnoreConfig []ScanFileIgnoreConfig `yaml:"fileignoreconfig,omitempty"`
		ScopeConfig      []ScopeConfig          `yaml:"scopeconfig,omitempty"`
		CustomPatterns   []PatternString        `yaml:"custom_patterns,omitempty"`
		CustomSeverities []CustomSeverityConfig `yaml:"custom_severities,omitempty"`
		AllowedPatterns  []string               `yaml:"allowed_patterns,omitempty"`
		Experimental     ExperimentalConfig     `yaml:"experimental,omitempty"`
		Threshold        severity.Severity      `default:"1" yaml:"threshold,omitempty"`
	} `yaml:"scanconfig,omitempty"`
	Version string `default:"1" yaml:"version,required"`
}

//AcceptsAll returns true if there are no rules specified
func (tRC *TalismanRC) AcceptsAll() bool {
	return len(tRC.effectiveRules("any-detector")) == 0
}

//Accept answers true if the Addition.Path is configured to be checked by the detectors
func (tRC *TalismanRC) Accept(addition gitrepo.Addition, detectorName string) bool {
	return !tRC.Deny(addition, detectorName)
}

func (tRC *TalismanRC) FilterAdditions(additions []gitrepo.Addition) []gitrepo.Addition {
	var applicableScopeFileNames []string
	if tRC.ScopeConfig != nil {
		for _, scope := range tRC.ScopeConfig {
			if len(knownScopes[scope.ScopeName]) > 0 {
				applicableScopeFileNames = append(applicableScopeFileNames, knownScopes[scope.ScopeName]...)
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

func (tRC *TalismanRC) IsEmpty() bool {
	return reflect.DeepEqual(TalismanRC{}, tRC)
}
