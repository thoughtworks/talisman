package detector

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"talisman/gitrepo"

	log "github.com/Sirupsen/logrus"
)

type fn func(fc *FileContentDetector, word string) string

type FileContentDetector struct {
	base64Detector     *Base64Detector
	hexDetector        *HexDetector
	creditCardDetector *CreditCardDetector
}

func NewFileContentDetector() *FileContentDetector {
	fc := FileContentDetector{}
	fc.base64Detector = NewBase64Detector()
	fc.hexDetector = NewHexDetector()
	fc.creditCardDetector = NewCreditCardDetector()
	return &fc
}

func (fc *FileContentDetector) AggressiveMode() *FileContentDetector {
	fc.base64Detector.aggressiveDetector = &Base64AggressiveDetector{}
	return fc
}

type contentType int

const (
	base64Content contentType = iota
	hexContent
	creditCardContent
)

func (ct contentType) getInfo() string {
	switch ct {
	case base64Content:
		return "Failing file as it contains a base64 encoded text."
	case hexContent:
		return "Failing file as it contains a hex encoded text."
	case creditCardContent:
		return "Failing file as it contains a potential credit card number."
	}
	return ""
}

func (ct contentType) getOutput() string {
	switch ct {
	case base64Content:
		return "Expected file to not to contain base64 encoded texts such as: %s"
	case hexContent:
		return "Expected file to not to contain hex encoded texts such as: %s"
	case creditCardContent:
		return "Expected file to not to contain credit card numbers such as: %s"
	}

	return ""
}

type content struct {
	name        gitrepo.FileName
	path        gitrepo.FilePath
	contentType contentType
	results     []string
}

func (fc *FileContentDetector) Test(additions []gitrepo.Addition, ignoreConfig TalismanRCIgnore, result *DetectionResults) {
	cc := NewChecksumCompare(additions, ignoreConfig)
	re := regexp.MustCompile(`(?i)checksum[ \t]*:[ \t]*[0-9a-fA-F]+`)

	contents := make(chan content, 512)
	ignoredFilePaths := make(chan gitrepo.FilePath, len(additions))

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(len(additions))
	for _, addition := range additions {
		go func(addition gitrepo.Addition) {
			defer waitGroup.Done()
			if ignoreConfig.Deny(addition, "filecontent") || cc.IsScanNotRequired(addition) {
				ignoredFilePaths <- addition.Path
				return
			}

			if string(addition.Name) == DefaultRCFileName {
				content := re.ReplaceAllString(string(addition.Data), "")
				data := []byte(content)
				addition.Data = data
			}
			contents <- content{
				name:        addition.Name,
				path:        addition.Path,
				contentType: base64Content,
				results:     fc.detectFile(addition.Data, checkBase64),
			}
			contents <- content{
				name:        addition.Name,
				path:        addition.Path,
				contentType: hexContent,
				results:     fc.detectFile(addition.Data, checkHex),
			}
			contents <- content{
				name:        addition.Name,
				path:        addition.Path,
				contentType: creditCardContent,
				results:     fc.detectFile(addition.Data, checkCreditCardNumber),
			}
		}(addition)
	}
	go func() {
		waitGroup.Wait()
		close(ignoredFilePaths)
		close(contents)
	}()

	for ignoredChanHasMore, contentChanHasMore := true, true; ignoredChanHasMore || contentChanHasMore; {
		select {
		case ignoredFilePath, hasMore := <-ignoredFilePaths:
			processIgnoredFilepath(ignoredFilePath, result)
			ignoredChanHasMore = hasMore
		case c, hasMore := <-contents:
			processContent(c, result)
			contentChanHasMore = hasMore
		}
	}
}

func processIgnoredFilepath(path gitrepo.FilePath, result *DetectionResults) {
	log.WithFields(log.Fields{
		"filePath": path,
	}).Info("Ignoring addition as it was specified to be ignored.")
	result.Ignore(path, "filecontent")
}

func processContent(c content, result *DetectionResults) {
	for _, res := range c.results {
		if res != "" {
			log.WithFields(log.Fields{
				"filePath": c.path,
			}).Info(c.contentType.getInfo())
			if string(c.name) == DefaultRCFileName {
				result.Warn(c.path, "filecontent", fmt.Sprintf(c.contentType.getOutput(), res), []string{})
			} else {
				result.Fail(c.path, "filecontent", fmt.Sprintf(c.contentType.getOutput(), res), []string{})
			}
		}
	}
}

func (fc *FileContentDetector) detectFile(data []byte, getResult fn) []string {
	content := string(data)
	return fc.checkEachLine(content, getResult)
}

func (fc *FileContentDetector) checkEachLine(content string, getResult fn) []string {
	lines := strings.Split(content, "\n")
	res := []string{}
	for _, line := range lines {
		lineResult := fc.checkEachWord(line, getResult)
		if len(lineResult) > 0 {
			res = append(res, lineResult...)
		}
	}
	return res
}

func (fc *FileContentDetector) checkEachWord(line string, getResult fn) []string {
	words := strings.Fields(line)
	res := []string{}
	for _, word := range words {
		wordResult := getResult(fc, word)
		if wordResult != "" {
			res = append(res, wordResult)
		}
	}
	return res
}

func checkBase64(fc *FileContentDetector, word string) string {
	return fc.base64Detector.checkBase64Encoding(word)
}

func checkCreditCardNumber(fc *FileContentDetector, word string) string {
	return fc.creditCardDetector.checkCreditCardNumber(word)
}

func checkHex(fc *FileContentDetector, word string) string {
	return fc.hexDetector.checkHexEncoding(word)
}
