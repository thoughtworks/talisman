package detector

import (
	"bufio"
	"fmt"
	"strings"
	"os"
)

type WordCheck struct {
}

const AVERAGE_LENGTH_OF_WORDS_IN_ENGLISH = 5 //See http://bit.ly/2qYFzFf for reference
const UNIX_WORDS_PATH = "/usr/share/dict/words" //See https://en.wikipedia.org/wiki/Words_(Unix) for reference
const UNIX_WORDS_ALTERNATIVE_PATH = "/usr/dict/words" //See https://en.wikipedia.org/wiki/Words_(Unix) for reference

func (en *WordCheck) containsWordsOnly(text string) bool {
	text = strings.ToLower(text)
	file := &os.File{}
	defer file.Close()
	reader := getWordsFileReader(file, UNIX_WORDS_PATH, UNIX_WORDS_ALTERNATIVE_PATH)
	if reader == nil {
		return false
	}
	wordCount := howManyWordsExistInText(reader, text)
	if wordCount >= (len(text) / (AVERAGE_LENGTH_OF_WORDS_IN_ENGLISH)) {
		return true
	}
	return false
}

func getWordsFileReader(file *os.File, filePaths... string) *bufio.Reader {
	for _, filePath := range filePaths {
		var err error = nil
		file, err = os.Open(filePath)
		if err != nil {
			continue
		}
		return bufio.NewReader(file)
	}
	return nil
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
			fmt.Println("wordCount, ", wordCount)
			break
		}
	}
	return wordCount
}
