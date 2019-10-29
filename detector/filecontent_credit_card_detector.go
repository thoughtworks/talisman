package detector

import "regexp"

type CreditCardDetector struct {
	creditCardRegex []*regexp.Regexp
}

var (
	creditCardPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?:3[47][0-9]{13})`),
		regexp.MustCompile(`(?:3(?:0[0-5]|[68][0-9])[0-9]{11})`),
		regexp.MustCompile(`^65[4-9][0-9]{13}|64[4-9][0-9]{13}|6011[0-9]{12}|(622(?:12[6-9]|1[3-9][0-9]|[2-8][0-9][0-9]|9[01][0-9]|92[0-5])[0-9]{10})$`),
		regexp.MustCompile(`^(?:2131|1800|35\d{3})\d{11}$`),
		regexp.MustCompile(`^(5018|5020|5038|6304|6759|6761|6763)[0-9]{8,15}$`),
		regexp.MustCompile(`(?:(?:5[1-5][0-9]{2}|222[1-9]|22[3-9][0-9]|2[3-6][0-9]{2}|27[01][0-9]|2720)[0-9]{12})`),
		regexp.MustCompile(`((?:4[0-9]{12})(?:[0-9]{3})?)`),

	}
)

func (detector CreditCardDetector) checkCreditCardNumber(content string) string {
	if !isLuhnNumber(content) {
		return ""
	}
	for _, regex := range detector.creditCardRegex {
		if regex.MatchString(content) {
			return content
		}
	}
	return ""
}

func NewCreditCardDetector() *CreditCardDetector {
	return &CreditCardDetector{creditCardPatterns}
}

func isLuhnNumber(content string) bool {
	var isAlternate bool
	var checksum int

	for position := len(content) - 1; position > -1; position-- {
		const ASCII_INDEX = 48
		digit := int(content[position] - ASCII_INDEX)
		if isAlternate {
			digit = digit * 2
			if digit > 9 {
				digit = (digit % 10) + 1
			}
		}
		isAlternate = !isAlternate
		checksum += digit
	}
	if checksum%10 == 0 {
		return true
	}
	return false
}
