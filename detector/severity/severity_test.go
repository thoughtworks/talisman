package severity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldReturnSeverityStringForDefinedSeverity(t *testing.T) {
	assert.Equal(t, String(Low), "low")
	assert.Equal(t, String(Medium), "medium")
	assert.Equal(t, String(High), "high")
}
func TestShouldReturnEmptyForInvalidSeverity(t *testing.T) {
	assert.Equal(t, String(10), "")
}

func TestShouldReturnSeverityForDefinedStrings(t *testing.T) {
	severityValue, _ := FromString("low")
	assert.Equal(t, severityValue, Low)
	severityValue, _ = FromString("mEDIum")
	assert.Equal(t, severityValue, Medium)
	severityValue, _ = FromString("HIGH")
	assert.Equal(t, severityValue, High)
}

func TestShouldReturnSeverityZeroWithErrorForUnknownStrings(t *testing.T) {
	severityValue, err := FromString("FakeSeverity")
	assert.Equal(t, Severity(0), severityValue)
	assert.Error(t, err)
	assert.Equal(t, "unknown severity FakeSeverity", err.Error())
}
