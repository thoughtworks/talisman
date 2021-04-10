package filecontent

import (
	log "github.com/Sirupsen/logrus"
	"talisman/talismanrc"
)

const BASE64_CHARS = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
const BASE64_ENTROPY_THRESHOLD = 4.5
const MIN_BASE64_SECRET_LENGTH = 20

type Base64Detector struct {
	base64Map              map[string]bool
	AggressiveDetector     *Base64AggressiveDetector
	entropy                *Entropy
	wordCheck              *WordCheck
	base64EntropyThreshold float64
}

func NewBase64Detector(tRC *talismanrc.TalismanRC) *Base64Detector {
	bd := Base64Detector{}
	bd.initBase64Map()
	bd.AggressiveDetector = nil

	bd.base64EntropyThreshold = BASE64_ENTROPY_THRESHOLD
	if tRC.GetExperimental().Base64EntropyThreshold > 0.0 {
		bd.base64EntropyThreshold = tRC.Experimental.Base64EntropyThreshold
		log.Debugf("Setting b64 entropy threshold to %f", bd.base64EntropyThreshold)
	}

	bd.entropy = &Entropy{}
	return &bd
}

func (bd *Base64Detector) initBase64Map() {
	bd.base64Map = map[string]bool{}
	for i := 0; i < len(BASE64_CHARS); i++ {
		bd.base64Map[string(BASE64_CHARS[i])] = true
	}
}

func (bd *Base64Detector) CheckBase64Encoding(word string) string {
	entropyCandidates := bd.entropy.GetEntropyCandidatesWithinWord(word, MIN_BASE64_SECRET_LENGTH, bd.base64Map)
	for _, candidate := range entropyCandidates {
		entropy := bd.entropy.GetShannonEntropy(candidate, BASE64_CHARS)
		log.Debugf("Detected entropy for word %s = %f", candidate, entropy)
		if entropy > bd.base64EntropyThreshold && !bd.wordCheck.containsWordsOnly(candidate) {
			return word
		}
	}
	if bd.AggressiveDetector != nil {
		return bd.AggressiveDetector.Test(word)
	}
	return ""
}
