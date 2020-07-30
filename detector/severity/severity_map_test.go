package severity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldReturnSeverityStringForDefinedSeverity(t *testing.T) {
	assert.Equal(t, SeverityValueToString(LowSeverity), "low")
	assert.Equal(t, SeverityValueToString(MediumSeverity), "medium")
	assert.Equal(t, SeverityValueToString(HighSeverity), "high")
}
func TestShouldReturnEmptyForInvalidSeverity(t *testing.T) {
	assert.Equal(t, SeverityValueToString(10), "")
}

func TestShouldReturnSeverityValueForDefinedStrings(t *testing.T) {
	assert.Equal(t, SeverityStringToValue("Low"), LowSeverity)
	assert.Equal(t, SeverityStringToValue("MEDIUM"), MediumSeverity)
	assert.Equal(t, SeverityStringToValue("high"), HighSeverity)
}

func TestShouldReturnSeverityZeroForUnknownStrings(t *testing.T) {
	assert.Equal(t, SeverityStringToValue("FakeSeverity"), SeverityValue(0))
}
