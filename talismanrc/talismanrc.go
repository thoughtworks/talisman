package talismanrc

import (
	"os"
	"regexp"
	"sort"

	logr "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"talisman/detector/severity"

	"talisman/gitrepo"
)

type Mode int
type CommitID string

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
	FileIgnoreConfig []FileIgnoreConfig              `yaml:"fileignoreconfig,omitempty"`
	ScopeConfig      []ScopeConfig                   `yaml:"scopeconfig,omitempty"`
	CustomPatterns   []PatternString                 `yaml:"custom_patterns,omitempty"`
	CustomSeverities []CustomSeverityConfig          `yaml:"custom_severities,omitempty"`
	AllowedPatterns  []string                        `yaml:"allowed_patterns,omitempty"`
	Experimental     ExperimentalConfig              `yaml:"experimental,omitempty"`
	Threshold        severity.Severity               `default:"low" yaml:"threshold,omitempty"`
	ScanConfig       map[CommitID][]FileIgnoreConfig `yaml:"scanconfig,omitempty"`
	Version          string                          `default:"2.0" yaml:"version"`
}

//SuggestRCFor returns the talismanRC file content corresponding to input ignore configs
func SuggestRCFor(configs []IgnoreConfig) string {
	fileIgnoreConfigs := []FileIgnoreConfig{}
	for _, config := range configs {
		fIC, ok := config.(*FileIgnoreConfig)
		if ok {
			fileIgnoreConfigs = append(fileIgnoreConfigs, *fIC)
		} else {
			logr.Debugf("Ignoring unknown IgnoreConfig : %#v", config)
		}
	}
	pRC := persistedRC{FileIgnoreConfig: fileIgnoreConfigs}
	result, _ := yaml.Marshal(pRC)

	return string(result)
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
		}
		ignoreEntries, _ := yaml.Marshal(&talismanRCConfig)
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

func fromPersistedRC(configFromTalismanRCFile *persistedRC, mode Mode) *TalismanRC {
	tRC := TalismanRC{}

	tRC.Threshold = configFromTalismanRCFile.Threshold
	tRC.ScopeConfig = configFromTalismanRCFile.ScopeConfig
	tRC.Experimental = configFromTalismanRCFile.Experimental
	tRC.CustomPatterns = configFromTalismanRCFile.CustomPatterns
	tRC.Experimental = configFromTalismanRCFile.Experimental
	tRC.AllowedPatterns = make([]*regexp.Regexp, len(configFromTalismanRCFile.AllowedPatterns))
	for i, p := range configFromTalismanRCFile.AllowedPatterns {
		tRC.AllowedPatterns[i] = regexp.MustCompile(p)
	}

	if mode == HookMode {
		tRC.IgnoreConfigs = make(
			[]IgnoreConfig,
			len(configFromTalismanRCFile.FileIgnoreConfig),
		)

		for i := range configFromTalismanRCFile.FileIgnoreConfig {
			tRC.IgnoreConfigs[i] = &configFromTalismanRCFile.FileIgnoreConfig[i]
		}
	}
	tRC.base = configFromTalismanRCFile

	return &tRC
}

func For(mode Mode) *TalismanRC {
	configFromTalismanRCFile := ConfigFromFile()
	talismanRC := fromPersistedRC(configFromTalismanRCFile, mode)
	return talismanRC
}
