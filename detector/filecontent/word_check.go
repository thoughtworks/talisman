package filecontent

import (
	"bufio"
	log "github.com/Sirupsen/logrus"
	"os"
	"strings"
)

type WordCheck struct {
}

const AVERAGE_LENGTH_OF_WORDS_IN_ENGLISH = 5 //See http://bit.ly/2qYFzFf for reference

func (en *WordCheck) containsWordsOnly(text string) bool {
	text = strings.ToLower(text)
	file := &os.File{}
	defer file.Close()
	reader := bufio.NewReader(strings.NewReader(DictionaryWordsString))
	if reader == nil {
		return false
	}
	wordCount := howManyWordsExistInText(reader, text)
	if wordCount >= (len(text) / (AVERAGE_LENGTH_OF_WORDS_IN_ENGLISH)) {
		return true
	}
	return false
}

func howManyWordsExistInText(reader *bufio.Reader, text string) int {
	wordCount := 0
	for {
		word, err := reader.ReadString('\n')
		word = strings.Trim(word, "\n")

		if word != "" && len(word) > 2 && strings.Contains(text, word) {
			text = strings.Replace(text, word, "", 1) //already matched
			wordCount++
		}

		if err != nil { //EOF
			log.Debugf("[WordChecker]: Found %d words", wordCount)
			break
		}
	}
	return wordCount
}
