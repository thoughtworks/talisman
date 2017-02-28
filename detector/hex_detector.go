package detector

const HEX_CHARS = "1234567890abcdefABCDEF"
const HEX_ENTROPY_THRESHOLD = 2.7
const MIN_HEX_SECRET_LENGTH = 20

type HexDetector struct {
	hexMap map[string]bool
	entropy *Entropy
}

func NewHexDetector() *HexDetector {
	bd := HexDetector{}
	bd.initHexMap()
	bd.entropy = &Entropy{}
	return &bd
}

func (hd *HexDetector) initHexMap() {
	hd.hexMap = map[string]bool{}
	for i := 0; i < len(HEX_CHARS); i++ {
		hd.hexMap[string(HEX_CHARS[i])] = true
	}
}

func (hd *HexDetector) checkHexEncoding(word string) string {
	entropyCandidates := hd.entropy.GetEntropyCandidatesWithinWord(word, MIN_HEX_SECRET_LENGTH, hd.hexMap)
	for _, candidate := range entropyCandidates {
		entropy := hd.entropy.GetShannonEntropy(candidate, HEX_CHARS)
		if entropy > HEX_ENTROPY_THRESHOLD {
			return word
		}
	}
	return ""
}