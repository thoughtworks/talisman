package helpers

import (
	"io"
	"talisman/gitrepo"
	mockchecksumcalculator "talisman/internal/mock/checksumcalculator"
	"talisman/talismanrc"

	"github.com/golang/mock/gomock"
	logr "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"testing"
)

func init() {
	logr.SetOutput(io.Discard)
}
func TestChecksumCompare_IsScanNotRequired(t *testing.T) {

	t.Run("should return false if talismanrc is empty", func(t *testing.T) {
		ignoreConfig := &talismanrc.TalismanRC{
			IgnoreConfigs: []talismanrc.IgnoreConfig{},
		}
		cc := NewChecksumCompare(nil, ignoreConfig)
		addition := gitrepo.Addition{Path: "some.txt"}

		required := cc.IsScanNotRequired(addition)

		assert.False(t, required)
	})

	t.Run("should loop through talismanrc configs", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		checksumCalculator := mockchecksumcalculator.NewMockChecksumCalculator(ctrl)
		ignoreConfig := talismanrc.TalismanRC{
			IgnoreConfigs: []talismanrc.IgnoreConfig{
				&talismanrc.FileIgnoreConfig{
					FileName: "some.txt",
					Checksum: "sha1",
				},
			},
		}
		cc := NewChecksumCompare(checksumCalculator, &ignoreConfig)
		addition := gitrepo.Addition{Name: "some.txt"}
		checksumCalculator.EXPECT().CalculateCollectiveChecksumForPattern("some.txt").Return("sha1")

		required := cc.IsScanNotRequired(addition)

		assert.True(t, required)
	})

}
