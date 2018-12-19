package detector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldReturnEmptyStringWhenCreditCardNumberDoesNotMatchAnyRegex(t *testing.T) {
	assert.Equal(t, "", NewCreditCardDetector().checkCreditCardNumber("1234123412341234"))
}

func TestShouldReturnEmptyStringWhenValidMasterCardPatternIsDetectedButNotValidLuhnNumber(t *testing.T) {
	assert.Equal(t, "", NewCreditCardDetector().checkCreditCardNumber("52222111111111111"))
}

func TestShouldReturnCardNumberWhenAmericanExpressCardNumberIsGiven(t *testing.T) {
	assert.Equal(t, "340000000000009", NewCreditCardDetector().checkCreditCardNumber("340000000000009"))
}

func TestShouldReturnCardNumberWhenDinersClubCardNumberIsGiven(t *testing.T) {
	assert.Equal(t, "30000000000004", NewCreditCardDetector().checkCreditCardNumber("30000000000004"))
}

func TestShouldReturnCardNumberWhenDiscoverCardNumberIsGiven(t *testing.T) {
	assert.Equal(t, "6011000000000004", NewCreditCardDetector().checkCreditCardNumber("6011000000000004"))
}

func TestShouldReturnCardNumberWhenJCBCardNumberIsGiven(t *testing.T) {
	assert.Equal(t, "3530111333300000", NewCreditCardDetector().checkCreditCardNumber("3530111333300000"))
}

func TestShouldReturnCardNumberWhenMaestroCardNumberIsGiven(t *testing.T) {
	assert.Equal(t, "6759649826438453", NewCreditCardDetector().checkCreditCardNumber("6759649826438453"))
}

func TestShouldReturnCardNumberWhenVisaCardNumberIsGiven(t *testing.T) {
	assert.Equal(t, "4111111111111111", NewCreditCardDetector().checkCreditCardNumber("4111111111111111"))
}
