package talismanrc

import (
	"os"
	"sort"

	logr "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"talisman/detector/severity"

	"talisman/gitrepo"
)

type TalismanRC struct {
	FileIgnoreConfig []FileIgnoreConfig     `yaml:"fileignoreconfig,omitempty"`
	ScopeConfig      []ScopeConfig          `yaml:"scopeconfig,omitempty"`
	CustomPatterns   []PatternString        `yaml:"custom_patterns,omitempty"`
	CustomSeverities []CustomSeverityConfig `yaml:"custom_severities,omitempty"`
	AllowedPatterns  []*Pattern             `yaml:"allowed_patterns,omitempty"`
	Experimental     ExperimentalConfig     `yaml:"experimental,omitempty"`
	Threshold        severity.Severity      `yaml:"threshold,omitempty"`
	Version          string                 `yaml:"version"`
}

// SuggestRCFor returns a string representation of a .talismanrc for the specified FileIgnoreConfigs
func SuggestRCFor(configs []FileIgnoreConfig) string {
	tRC := TalismanRC{FileIgnoreConfig: configs, Version: DefaultRCVersion}
	result, _ := yaml.Marshal(tRC)

	return string(result)
}

// RemoveScopedFiles removes scope files from additions
func (tRC *TalismanRC) RemoveScopedFiles(additions []gitrepo.Addition) []gitrepo.Addition {
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
				break
			}
		}
		if !isFilePresentInScope {
			result = append(result, addition)
		}
	}
	return result
}

// AddIgnores inserts the specified FileIgnoreConfigs to an existing .talismanrc file, or creates one if it doesn't exist.
func (tRC *TalismanRC) AddIgnores(entriesToAdd []FileIgnoreConfig) {
	if len(entriesToAdd) > 0 {
		logr.Debugf("Adding entries: %v", entriesToAdd)
		tRC.FileIgnoreConfig = combineFileIgnores(tRC.FileIgnoreConfig, entriesToAdd)

		ignoreEntries, _ := yaml.Marshal(&tRC)
		file, err := fs.OpenFile(currentRCFileName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logr.Errorf("error opening %s: %s", currentRCFileName, err)
		}
		defer func() {
			err := file.Close()
			if err != nil {
				logr.Errorf("error closing %s: %s", currentRCFileName, err)
			}
		}()
		logr.Debugf("Writing talismanrc: %v", string(ignoreEntries))
		_, err = file.WriteString(string(ignoreEntries))
		if err != nil {
			logr.Errorf("error writing to %s: %s", currentRCFileName, err)
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
	for k := range existingMap {
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

// RemoveAllowedPatterns removes globally- and per-file allowed patterns from an Addition
func (tRC *TalismanRC) RemoveAllowedPatterns(addition gitrepo.Addition) string {
	additionPathAsString := string(addition.Path)
	// Processing global allowed patterns
	for _, pattern := range tRC.AllowedPatterns {
		addition.Data = pattern.ReplaceAll(addition.Data, []byte(""))
	}

	// Processing allowed patterns based on file path
	for _, ignoreConfig := range tRC.FileIgnoreConfig {
		if ignoreConfig.GetFileName() == additionPathAsString {
			for _, pattern := range ignoreConfig.GetAllowedPatterns() {
				addition.Data = pattern.ReplaceAll(addition.Data, []byte(""))
			}
		}
	}
	return string(addition.Data)
}

// Deny answers true if the Addition should NOT be checked by the specified detector
func (tRC *TalismanRC) Deny(addition gitrepo.Addition, detectorName string) bool {
	for _, pattern := range tRC.effectiveRules(detectorName) {
		if addition.Matches(pattern) {
			return true
		}
	}
	return false
}

// Accept answers true if the Addition should be checked by the specified detector
func (tRC *TalismanRC) Accept(addition gitrepo.Addition, detectorName string) bool {
	return !tRC.Deny(addition, detectorName)
}

func (tRC *TalismanRC) effectiveRules(detectorName string) []string {
	var result []string
	for _, ignore := range tRC.FileIgnoreConfig {
		if ignore.isEffective(detectorName) {
			result = append(result, ignore.GetFileName())
		}
	}
	return result
}
