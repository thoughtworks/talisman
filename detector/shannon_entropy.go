package detector

import (
	"strings"
	"math"
)

type Entropy struct {

}

func (en *Entropy) GetShannonEntropy(str string, superSet string) float64 {
	if str == "" {
		return 0
	}
	entropy := 0.0
	for _, c := range superSet {
		p := float64(strings.Count(str, string(c))) / float64(len(str))
		if p > 0 {
			entropy -= p * math.Log2(p)
		}
	}
	return entropy
}

func (en *Entropy) GetEntropyCandidatesWithinWord(word string, minCandidateLength int, superSet map[string]bool) []string {
	candidates := []string{}
	count := 0
	subSet := ""
	if len(word) < minCandidateLength {
		return candidates
	}
	for _, c := range word {
		char := string(c)
		if superSet[char] {
			subSet += char
			count++
		} else {
			if count > minCandidateLength {
				candidates = append(candidates, subSet)
			}
			subSet = ""
			count = 0
		}
	}
	if count > minCandidateLength {
		candidates = append(candidates, subSet)
	}
	return candidates
}
