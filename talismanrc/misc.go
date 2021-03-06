package talismanrc

import (
	"log"
	yaml "gopkg.in/yaml.v2"
)

func ReadConfigFromRCFile(repoFileRead func(string) ([]byte, error)) *TalismanRC {
	fileContents, error := repoFileRead(currentRCFileName)
	if error != nil {
		panic(error)
	}
	return NewTalismanRC(fileContents)
}

func NewTalismanRC(fileContents []byte) *TalismanRC {
	talismanRCFromFile := TalismanRC{}
	err := yaml.Unmarshal(fileContents, &talismanRCFromFile)
	if err != nil {
		log.Println("Unable to parse .talismanrc")
		log.Printf("error: %v", err)
		return &TalismanRC{}
	}
	return &talismanRCFromFile
}
