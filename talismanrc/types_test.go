package talismanrc

import (
	"io"
	"regexp"
	"strings"
	"testing"

	logr "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func init() {
	logr.SetOutput(io.Discard)
}

func TestCustomMarshalling(t *testing.T) {
	t.Run("Can unmarshal yaml into a Pattern struct", func(t *testing.T) {
		savedPattern := []byte("text-pattern")
		fromText := Pattern{}
		err := yaml.Unmarshal(savedPattern, &fromText)
		assert.Nil(t, err, "Should have unmarshalled %s into a Pattern", savedPattern)
		assert.Equal(t, Pattern{regexp.MustCompile(string(savedPattern))}, fromText)
	})

	t.Run("Can marshal a Pattern struct into yaml", func(t *testing.T) {
		pattern := Pattern{regexp.MustCompile("pattern")}
		str, err := yaml.Marshal(pattern)
		assert.Nil(t, err, "Should have marshalled %v into a string of yaml", pattern)
		assert.Equal(t, pattern.String(), strings.TrimSpace(string(str)))
	})
}
