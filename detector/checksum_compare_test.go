package detector

import (
	"github.com/stretchr/testify/assert"
	"talisman/gitrepo"
	"talisman/talismanrc"
	"testing"
)

func TestChecksumCompare_IsScanNotRequired(t *testing.T) {

	t.Run("should return false if talismanrc is empty", func(t *testing.T) {
		ignoreConfig := talismanrc.NewTalismanRC(nil)
		cc := NewChecksumCompare([]gitrepo.Addition{}, []gitrepo.Addition{}, ignoreConfig)

		required := cc.IsScanNotRequired(gitrepo.Addition{})

		assert.False(t, required)
	})

	t.Run("should return false if talismanrc is empty", func(t *testing.T) {
		ignoreConfig := talismanrc.NewTalismanRC(nil)
		cc := NewChecksumCompare([]gitrepo.Addition{}, []gitrepo.Addition{}, ignoreConfig)

		required := cc.IsScanNotRequired(gitrepo.Addition{})

		assert.False(t, required)
	})

}
