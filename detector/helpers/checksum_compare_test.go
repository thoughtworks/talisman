package helpers

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"talisman/gitrepo"
	mockchecksumcalculator "talisman/internal/mock/checksumcalculator"
	mockutility "talisman/internal/mock/utility"
	"talisman/talismanrc"
	logr "github.com/Sirupsen/logrus"

	"testing"
)

func init() {
	logr.SetOutput(ioutil.Discard)
}
func TestChecksumCompare_IsScanNotRequired(t *testing.T) {

	t.Run("should return false if talismanrc is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSHA256Hasher := mockutility.NewMockSHA256Hasher(ctrl)
		ignoreConfig := talismanrc.NewTalismanRC(nil)
		cc := NewChecksumCompare(nil, mockSHA256Hasher, ignoreConfig)
		addition := gitrepo.Addition{Path: "some.txt"}
		mockSHA256Hasher.EXPECT().CollectiveSHA256Hash([]string{string(addition.Path)}).Return("somesha")

		required := cc.IsScanNotRequired(addition)

		assert.False(t, required)
	})

	t.Run("should loop through talismanrc configs", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSHA256Hasher := mockutility.NewMockSHA256Hasher(ctrl)
		checksumCalculator := mockchecksumcalculator.NewMockChecksumCalculator(ctrl)
		ignoreConfig := talismanrc.TalismanRC{
			FileIgnoreConfig: []talismanrc.FileIgnoreConfig{
				{
					FileName: "some.txt",
					Checksum: "sha1",
				},
			},
		}
		cc := NewChecksumCompare(checksumCalculator, mockSHA256Hasher, &ignoreConfig)
		addition := gitrepo.Addition{Name: "some.txt",}
		mockSHA256Hasher.EXPECT().CollectiveSHA256Hash([]string{string(addition.Path)}).Return("somesha")
		checksumCalculator.EXPECT().CalculateCollectiveChecksumForPattern("some.txt").Return("sha1")

		required := cc.IsScanNotRequired(addition)

		assert.True(t, required)
	})

}
