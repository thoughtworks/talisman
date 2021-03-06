package filecontent

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"talisman/detector/helpers"
	"talisman/detector/severity"
	"talisman/gitrepo"
	"talisman/talismanrc"

	log "github.com/Sirupsen/logrus"
)

type fn func(fc *FileContentDetector, word string) string

type FileContentDetector struct {
	base64Detector         *Base64Detector
	hexDetector            *HexDetector
	creditCardDetector     *CreditCardDetector
	base64EntropyThreshold float64
}

func NewFileContentDetector(tRC *talismanrc.TalismanRC) *FileContentDetector {
	fc := FileContentDetector{}
	fc.base64Detector = NewBase64Detector(tRC)
	fc.hexDetector = NewHexDetector()
	fc.creditCardDetector = NewCreditCardDetector()
	return &fc
}

func (fc *FileContentDetector) AggressiveMode() *FileContentDetector {
	fc.base64Detector.AggressiveDetector = &Base64AggressiveDetector{}
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

func (ct contentType) getMessageFormat() string {
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
	severity    severity.Severity
}

func (fc *FileContentDetector) Test(comparator helpers.ChecksumCompare, currentAdditions []gitrepo.Addition, ignoreConfig *talismanrc.TalismanRC, result *helpers.DetectionResults, additionCompletionCallback func()) {
	contentTypes := []struct {
		contentType
		fn
		severity severity.Severity
	}{
		{
			contentType: base64Content,
			fn:          checkBase64,
			severity:    severity.SeverityConfiguration["Base64Content"],
		},
		{
			contentType: hexContent,
			fn:          checkHex,
			severity:    severity.SeverityConfiguration["HexContent"],
		},
		{
			contentType: creditCardContent,
			fn:          checkCreditCardNumber,
			severity:    severity.SeverityConfiguration["CreditCardContent"],
		},
	}
	re := regexp.MustCompile(`(?i)checksum[ \t]*:[ \t]*[0-9a-fA-F]+`)

	contents := make(chan content, 512)
	ignoredFilePaths := make(chan gitrepo.FilePath, len(currentAdditions))

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(len(currentAdditions))
	for _, addition := range currentAdditions {
		go func(addition gitrepo.Addition) {
			defer waitGroup.Done()
			defer additionCompletionCallback()
			if ignoreConfig.Deny(addition, "filecontent") || comparator.IsScanNotRequired(addition) {
				ignoredFilePaths <- addition.Path
				return
			}

			if string(addition.Name) == talismanrc.DefaultRCFileName {
				content := re.ReplaceAllString(string(addition.Data), "")
				data := []byte(content)
				addition.Data = data
			}
			for _, ct := range contentTypes {
				contents <- content{
					name:        addition.Name,
					path:        addition.Path,
					contentType: ct.contentType,
					results:     fc.detectFile(addition.Data, ct.fn),
					severity:    ct.severity,
				}
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
			if !hasMore {
				ignoredChanHasMore = false
				continue
			}
			processIgnoredFilepath(ignoredFilePath, result)
		case c, hasMore := <-contents:
			if !hasMore {
				contentChanHasMore = false
				continue
			}
			processContent(c, ignoreConfig.Threshold, result)
		}
	}
}

func processIgnoredFilepath(path gitrepo.FilePath, result *helpers.DetectionResults) {
	log.WithFields(log.Fields{
		"filePath": path,
	}).Info("Ignoring addition as it was specified to be ignored.")
	result.Ignore(path, "filecontent")
}

func processContent(c content, threshold severity.Severity, result *helpers.DetectionResults) {
	for _, res := range c.results {
		if res != "" {
			log.WithFields(log.Fields{
				"filePath": c.path,
			}).Info(c.contentType.getInfo())
			if string(c.name) == talismanrc.DefaultRCFileName || !c.severity.ExceedsThreshold(threshold) {
				result.Warn(c.path, "filecontent", fmt.Sprintf(c.contentType.getMessageFormat(), formatForReporting(res)), []string{}, c.severity)
			} else {
				result.Fail(c.path, "filecontent", fmt.Sprintf(c.contentType.getMessageFormat(), formatForReporting(res)), []string{}, c.severity)
			}
		}
	}
}

func formatForReporting(input string) string {
	if len(input) > 50 {
		return input[:47] + "..."
	}
	return input
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
	return fc.base64Detector.CheckBase64Encoding(word)
}

func checkCreditCardNumber(fc *FileContentDetector, word string) string {
	return fc.creditCardDetector.checkCreditCardNumber(word)
}

func checkHex(fc *FileContentDetector, word string) string {
	return fc.hexDetector.CheckHexEncoding(word)
}
