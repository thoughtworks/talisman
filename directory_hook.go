package main

import (
	"talisman/gitrepo"
	"talisman/utility"

	log "github.com/Sirupsen/logrus"

	"github.com/bmatcuk/doublestar"
)

type DirectoryHook struct{}

func NewDirectoryHook() *DirectoryHook {
	return &DirectoryHook{}
}

func (p *DirectoryHook) GetFilesFromDirectory(globPattern string) []gitrepo.Addition {
	var result []gitrepo.Addition

	files, _ := doublestar.Glob(globPattern)
	for _, file := range files {
		data, err := ReadFile(file)

		if err != nil {
			continue
		}

		newAddition := gitrepo.NewAddition(file, data)
		result = append(result, newAddition)
	}

	return result
}

func ReadFile(filepath string) ([]byte, error) {
	log.Debugf("reading file %s", filepath)
	return utility.SafeReadFile(filepath)
}
