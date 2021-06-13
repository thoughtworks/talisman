package talismanrc

import (
	logr "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"reflect"
	"regexp"
	"sort"
	"talisman/detector/severity"

	"talisman/gitrepo"
)

type Mode int

const (
	HookMode = Mode(iota + 1)
	ScanMode
)

type TalismanRC struct {
	IgnoreConfigs    []IgnoreConfig         `yaml:"-"`
	ScopeConfig      []ScopeConfig          `yaml:"-"`
	CustomPatterns   []PatternString        `yaml:"-"`
	CustomSeverities []CustomSeverityConfig `yaml:"-"`
	AllowedPatterns  []*regexp.Regexp       `yaml:"-"`
	Experimental     ExperimentalConfig     `yaml:"-"`
	Threshold        severity.Severity      `yaml:"-"`
	base             *persistedRC
}

type persistedRC struct {
	FileIgnoreConfig []FileIgnoreConfig     `yaml:"fileignoreconfig,omitempty"`
	ScopeConfig      []ScopeConfig          `yaml:"scopeconfig,omitempty"`
	CustomPatterns   []PatternString        `yaml:"custom_patterns,omitempty"`
	CustomSeverities []CustomSeverityConfig `yaml:"custom_severities,omitempty"`
	AllowedPatterns  []string               `yaml:"allowed_patterns,omitempty"`
	Experimental     ExperimentalConfig     `yaml:"experimental,omitempty"`
	Threshold        severity.Severity      `default:"low" yaml:"threshold,omitempty"`
	ScanConfig       struct {
		FileIgnoreConfig []ScanFileIgnoreConfig `yaml:"scanfileignoreconfig,omitempty"`
		ScopeConfig      []ScopeConfig          `yaml:"scopeconfig,omitempty"`
		CustomPatterns   []PatternString        `yaml:"custom_patterns,omitempty"`
		CustomSeverities []CustomSeverityConfig `yaml:"custom_severities,omitempty"`
		AllowedPatterns  []string               `yaml:"allowed_patterns,omitempty"`
		Experimental     ExperimentalConfig     `yaml:"experimental,omitempty"`
		Threshold        severity.Severity      `default:"low" yaml:"threshold,omitempty"`
	} `yaml:"scanconfig,omitempty"`
	Version string `default:"1.0" yaml:"version"`
}

//AcceptsAll returns true if there are no rules specified
func (tRC *TalismanRC) AcceptsAll() bool {
	return len(tRC.effectiveRules("any-detector")) == 0
}

//Accept answers true if the Addition.Path is configured to be checked by the detectors
func (tRC *TalismanRC) Accept(addition gitrepo.Addition, detectorName string) bool {
	return !tRC.Deny(addition, detectorName)
}

//FilterAdditions removes scope files from additions
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
				break
			}
		}
		if !isFilePresentInScope {
			result = append(result, addition)
		}
	}
	return result
}

func (tRC *persistedRC) AddIgnores(mode Mode, entriesToAdd []IgnoreConfig) {
	if len(entriesToAdd) > 0 {
		logr.Debugf("Adding entries: %v", entriesToAdd)
		talismanRCConfig := ConfigFromFile()
		if mode == HookMode {
			fileIgnoreEntries := make([]FileIgnoreConfig, len(entriesToAdd))
			for idx, entry := range entriesToAdd {
				newVal, _ := entry.(*FileIgnoreConfig)
				fileIgnoreEntries[idx] = *newVal
			}
			talismanRCConfig.FileIgnoreConfig = combineFileIgnores(talismanRCConfig.FileIgnoreConfig, fileIgnoreEntries)
		} else {
			scanFileIgnoreEntries := make([]ScanFileIgnoreConfig, len(entriesToAdd))
			for idx, entry := range entriesToAdd {
				newVal, _ := entry.(*ScanFileIgnoreConfig)
				scanFileIgnoreEntries[idx] = *newVal
			}
			talismanRCConfig.ScanConfig.FileIgnoreConfig = combineScanFileIgnores(talismanRCConfig.ScanConfig.FileIgnoreConfig, scanFileIgnoreEntries)
		}
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

func combineScanFileIgnores(existing, incoming []ScanFileIgnoreConfig) []ScanFileIgnoreConfig {
	existingMap := make(map[string]ScanFileIgnoreConfig)
	for _, fIC := range existing {
		existingMap[fIC.FileName] = fIC
	}
	for _, fIC := range incoming {
		if efIC, ok := existingMap[fIC.FileName]; ok {
			efIC.AllowedPatterns = combine(efIC.AllowedPatterns, fIC.AllowedPatterns)
			efIC.Checksums = combine(efIC.Checksums, fIC.Checksums)
			efIC.IgnoreDetectors = combine(efIC.IgnoreDetectors, fIC.IgnoreDetectors)
		} else {
			existingMap[fIC.FileName] = fIC
		}
	}
	result := make([]ScanFileIgnoreConfig, len(existingMap))
	resultKeys := make([]string, len(existingMap))
	index := 0
	//sort keys in alphabetical order
	for k := range existingMap {
		resultKeys[index] = k
		index++
	}
	sort.Strings(resultKeys)
	//add result entries based on sorted keys
	for idx, k := range resultKeys {
		result[idx] = existingMap[k]
	}
	return result
}

func combine(existing []string, incoming []string) []string {
	combinedMap := make(map[string]bool)
	result := make([]string, 0)
	for _, src := range [][]string{existing, incoming} {
		for _, v := range src {
			if _, ok := combinedMap[v]; !ok {
				combinedMap[v] = true
				result = append(result, v)
			}
		}
	}
	return result
}

//Deny answers true if the Addition.Path is configured to be ignored and not checked by the detectors
func (tRC *TalismanRC) Deny(addition gitrepo.Addition, detectorName string) bool {
	for _, pattern := range tRC.effectiveRules(detectorName) {
		if addition.Matches(pattern) {
			return true
		}
	}
	return false
}

func (tRC *TalismanRC) effectiveRules(detectorName string) []string {
	var result []string
	for _, ignore := range tRC.IgnoreConfigs {
		if ignore.isEffective(detectorName) {
			result = append(result, ignore.GetFileName())
		}
	}

	return result
}

func (tRC *persistedRC) IsEmpty() bool {
	return reflect.DeepEqual(&persistedRC{}, tRC)
}

func fromPersistedRC(configFromTalismanRCFile *persistedRC, mode Mode) *TalismanRC {
	tRC := TalismanRC{}
	if mode == HookMode {
		tRC.Threshold = configFromTalismanRCFile.Threshold
		tRC.ScopeConfig = configFromTalismanRCFile.ScopeConfig
		tRC.Threshold = configFromTalismanRCFile.Threshold
		tRC.Experimental = configFromTalismanRCFile.Experimental
		tRC.CustomPatterns = configFromTalismanRCFile.CustomPatterns
		tRC.Experimental = configFromTalismanRCFile.Experimental
		tRC.AllowedPatterns = make([]*regexp.Regexp, len(configFromTalismanRCFile.AllowedPatterns))
		for i, p := range configFromTalismanRCFile.AllowedPatterns {
			tRC.AllowedPatterns[i] = regexp.MustCompile(p)
		}
		tRC.IgnoreConfigs = make([]IgnoreConfig, len(configFromTalismanRCFile.FileIgnoreConfig))
		for i, v := range configFromTalismanRCFile.FileIgnoreConfig {
			tRC.IgnoreConfigs[i] = IgnoreConfig(&v)
		}
	}

	if mode == ScanMode {
		scanconfigFromTalismanRCFile := configFromTalismanRCFile.ScanConfig
		tRC.Threshold = scanconfigFromTalismanRCFile.Threshold
		tRC.ScopeConfig = scanconfigFromTalismanRCFile.ScopeConfig
		tRC.Threshold = scanconfigFromTalismanRCFile.Threshold
		tRC.Experimental = scanconfigFromTalismanRCFile.Experimental
		tRC.CustomPatterns = scanconfigFromTalismanRCFile.CustomPatterns
		tRC.Experimental = scanconfigFromTalismanRCFile.Experimental
		tRC.AllowedPatterns = make([]*regexp.Regexp, len(scanconfigFromTalismanRCFile.AllowedPatterns))
		for i, p := range scanconfigFromTalismanRCFile.AllowedPatterns {
			tRC.AllowedPatterns[i] = regexp.MustCompile(p)
		}
		tRC.IgnoreConfigs = make([]IgnoreConfig, len(scanconfigFromTalismanRCFile.FileIgnoreConfig))
		for i, v := range scanconfigFromTalismanRCFile.FileIgnoreConfig {
			tRC.IgnoreConfigs[i] = IgnoreConfig(&v)
		}
	}
	tRC.base = configFromTalismanRCFile
	return &tRC
}

func For(mode Mode) *TalismanRC {
	configFromTalismanRCFile := ConfigFromFile()
	return fromPersistedRC(configFromTalismanRCFile, mode)
}
