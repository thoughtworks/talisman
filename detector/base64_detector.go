package detector

const BASE64_CHARS = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
const BASE64_ENTROPY_THRESHOLD = 4.5
const MIN_BASE64_SECRET_LENGTH = 20

type Base64Detector struct {
	base64Map map[string]bool
	aggressiveDetector *Base64AggressiveDetector
	entropy *Entropy
	wordCheck *WordCheck
}

func NewBase64Detector() *Base64Detector {
	bd := Base64Detector{}
	bd.initBase64Map()
	bd.aggressiveDetector = nil
	bd.entropy = &Entropy{}
	return &bd
}

func (bd *Base64Detector) initBase64Map() {
	bd.base64Map = map[string]bool{}
	for i := 0; i < len(BASE64_CHARS); i++ {
		bd.base64Map[string(BASE64_CHARS[i])] = true
	}
}

func (bd *Base64Detector) checkBase64Encoding(word string) string {
	entropyCandidates := bd.entropy.GetEntropyCandidatesWithinWord(word, MIN_BASE64_SECRET_LENGTH, bd.base64Map)
	for _, candidate := range entropyCandidates {
		entropy := bd.entropy.GetShannonEntropy(candidate, BASE64_CHARS)
		if entropy > BASE64_ENTROPY_THRESHOLD && !bd.wordCheck.containsWordsOnly(candidate) {
			return word
		}
	}
	if bd.aggressiveDetector != nil {
		return bd.aggressiveDetector.Test(word)
	}
	return ""
}

