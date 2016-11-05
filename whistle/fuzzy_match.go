package whistle

import (
	"strings"
)

func max(l1 int, l2 int) int {
	if l1 > l2 {
		return l1
	} else {
		return l2
	}
}

func min(l1 int, l2 int, l3 int) int {
	minVal := l1
	if l2 < minVal {
		minVal = l2
	} else if l3 < minVal {
		minVal = l3
	}
	return minVal
}

func TokenRatio(s1 string, s2 string) (int, int, int, []string, []string) {

	tokens1 := strings.Split(s1, " ")
	tokens2 := strings.Split(s2, " ")

	largestLen := max(len(tokens1), len(tokens2))

	distance, substTokens, canonicalString := levenshteinDistance(tokens1, tokens2)
	return (100 * (largestLen - distance) / largestLen), len(tokens1), len(tokens2), substTokens, canonicalString
}

func levenshteinDistance(tokens1 []string, tokens2 []string) (int, []string, []string) {

	// Cost of insert and delete is higher as compared to substitution cost
	const insertCost int = 2
	const delCost int = 2
	const maxCost = 100000

	substTokens := []string{""}
	canonicalString := []string{""}

	lastRow := make([]int, len(tokens2)+1)

	for i := 0; i < len(lastRow); i++ {
		lastRow[i] = i
	}

	lastSubstCost := 0
	lastCost := 0

	totalInsertionCost := 0
	totalDeletionCost := 0
	totalSubstitutionCost := 0

	for idx1 := 1; idx1 < len(tokens1); idx1++ {
		token1 := tokens1[idx1]
		newRow := make([]int, len(tokens2)+1)
		newRow[0] = idx1
		// set the substitution cost = 1
		substitutionCost := 1

		minSubstCost := maxCost
		minCost := maxCost
		minSubstToken := ""
		for idx2 := 1; idx2 < len(tokens2); idx2++ {
			token2 := tokens2[idx2]
			if token1 == token2 {
				substitutionCost = 0
			}
			totalInsertionCost = newRow[idx2-1] + insertCost
			totalDeletionCost = lastRow[idx2] + delCost
			totalSubstitutionCost = lastRow[idx2-1] + substitutionCost
			minCurrentTotalCost := min(totalSubstitutionCost, totalDeletionCost, totalInsertionCost)
			newRow[idx2] = minCurrentTotalCost

			minCost = min(minCurrentTotalCost, minCost, 100000)
			// if this was a true substitution and NOT a match
			// and if this is the lowest we have encountered this ROW
			if minCurrentTotalCost == totalSubstitutionCost && substitutionCost > 0 && minSubstCost > minCurrentTotalCost {
				minSubstCost = minCurrentTotalCost
				minSubstToken = token2
			}

		}

		if minCost == lastCost {
			canonicalString = append(canonicalString, token1)
		} else {
			canonicalString = append(canonicalString, "%s")
			lastCost = minCost
		}

		if minCost == minSubstCost && len(minSubstToken) > 0 && minSubstCost > lastSubstCost {
			substTokens = append(substTokens, minSubstToken)
			lastSubstCost = minSubstCost
		}

		lastRow = newRow
	}

	return lastRow[len(lastRow)-2], substTokens, canonicalString
}

const MIN_FUZZ_RATIO int = 65.0
const MAX_FUZZ_RATIO int = 80.0

func DoTheyFuzzyMatch(key string, message string) (bool, []string, string) {
	ratio, l1, l2, substTokens, canonicalStringParts := TokenRatio(key, message)
	max_len := max(l1, l2)
	fuzz_ratio := MIN_FUZZ_RATIO
	if max_len > 15 {
		fuzz_ratio = MAX_FUZZ_RATIO
	} else if 15 > max_len && max_len > 4 {

		fuzz_ratio = int(float32(max_len)*1.7) + MIN_FUZZ_RATIO
	}

	trimmedSubstTokens := make([]string, 0, len(substTokens))

	for _, tok := range substTokens {
		trimmedSubstTokens = append(trimmedSubstTokens, tok)
	}
	return fuzz_ratio <= ratio, trimmedSubstTokens, strings.Join(canonicalStringParts, " ")
}
