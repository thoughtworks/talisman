package detector

import (
	"gopkg.in/yaml.v2"
	"log"
	"reflect"
	"regexp"
	"strings"

	"talisman/gitrepo"
)

const (
	//LinePattern represents a line in the ignorefile with an optional comment
	LinePattern string = "^([^#]+)?\\s*(#(.*))?$"

	//IgnoreDetectorCommentPattern represents a special comment that ignores only certain detectors
	IgnoreDetectorCommentPattern string = "^ignore:([^\\s]+).*$"

	//DefaultRCFileName represents the name of default file in which all the ignore patterns are configured in new version
	DefaultRCFileName string = ".talismanrc"
)

//Ignores represents a set of patterns that have been configured to be ignored by the Detectors.
//Detectors are expected to honor these ignores.
type Ignores struct {
	patterns []Ignore
}
//Ignore represents a single pattern and its comment
type Ignore struct {
	pattern string
	comment string
	ignoredDetectors []string
}

type FileIgnoreConfig struct {
	FileName        string `yaml:"filename"`
	Checksum        string `yaml:"checksum"`
	IgnoreDetectors []string `yaml:"ignore_detectors"`
}

type ScopeConfig struct {
	ScopeName string `yaml:"scope"`
}

type TalismanRCIgnore struct {
	FileIgnoreConfig []FileIgnoreConfig `yaml:"fileignoreconfig"`
	ScopeConfig      []ScopeConfig      `yaml:"scopeconfig"`
}

func (ignore TalismanRCIgnore) IsEmpty() bool {
	return reflect.DeepEqual(TalismanRCIgnore{}, ignore)
}

func ReadConfigFromRCFile(repoFileRead func(string) ([]byte, error)) TalismanRCIgnore {
	fileContents, error := repoFileRead(DefaultRCFileName)
	if error != nil {
		panic(error)
	}
	return NewTalismanRCIgnore(fileContents)
}


func NewTalismanRCIgnore(fileContents []byte) (TalismanRCIgnore) {
	talismanRCIgnore := TalismanRCIgnore{}
	err := yaml.Unmarshal([]byte(fileContents), &talismanRCIgnore)
	if err != nil {
		log.Println("Unable to parse .talismanrc")
		log.Printf("error: %v", err)
		return talismanRCIgnore
	}
	return talismanRCIgnore
}

func NewIgnore(pattern string, comment string) Ignore {
	var ignoredDetectors []string
	ignorePattern := regexp.MustCompile(IgnoreDetectorCommentPattern)
	match := ignorePattern.FindStringSubmatch(comment)
	if match != nil {
		ignoredDetectors = strings.Split(match[1], ",")
	}

	return Ignore{
		pattern: pattern,
		comment: comment,
		ignoredDetectors: ignoredDetectors,
	}
}

func (i FileIgnoreConfig) isEffective(detectorName string) bool {
	return !isEmptyString(i.FileName) &&
		contains(i.IgnoreDetectors, detectorName)
}


//NewIgnores builds a new Ignores with the patterns specified in the ignoreSpecs
//Empty lines and comments are ignored.
func NewIgnores(lines ...string) Ignores {
	var ignores []Ignore
	for _, line := range lines {
		var commentPattern = regexp.MustCompile(LinePattern)
		groups := commentPattern.FindStringSubmatch(line)
		if len(groups) == 4 {
			ignores = append(ignores, NewIgnore(strings.TrimSpace(groups[1]), strings.TrimSpace(groups[3])))
		}
	}
	return Ignores{ignores}
}

//AcceptsAll returns true if there are no rules specified
func (i TalismanRCIgnore) AcceptsAll() bool {
	return len(i.effectiveRules("any-detector")) == 0
}

//Accept answers true if the Addition.Path is configured to be checked by the detectors
func (i TalismanRCIgnore) Accept(addition gitrepo.Addition, detectorName string) bool {
	return !i.Deny(addition, detectorName)
}

func IgnoreAdditionsByScope(additions []gitrepo.Addition, rcConfigIgnores TalismanRCIgnore, scopeMap map[string][]string) []gitrepo.Addition {
	var applicableScopeFileNames []string
	if rcConfigIgnores.ScopeConfig != nil {
		for _, scope := range rcConfigIgnores.ScopeConfig {
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
//Deny answers true if the Addition.Path is configured to be ignored and not checked by the detectors
func (i TalismanRCIgnore) Deny(addition gitrepo.Addition, detectorName string) bool {
	result := false
	for _, pattern := range i.effectiveRules(detectorName) {
		result = result || addition.Matches(pattern)
	}
	return result
}

func (i TalismanRCIgnore) effectiveRules(detectorName string) []string {
	var result []string
	for _, ignore := range i.FileIgnoreConfig {
		if ignore.isEffective(detectorName) {
			result = append(result, ignore.FileName)
		}
	}
	return result
}

func isEmptyString(str string) bool {
	var emptyStringPattern = regexp.MustCompile("^\\s*$")
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
