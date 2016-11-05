package whistle

import (
	"bufio"
	"os"
	"testing"
)

func TestTokenRatio(t *testing.T) {
	s1 := "[resmgr] [Thread-2] Marking the host 87987c71-90f7-4d56-acca-c11adae5c2ed as not responding"
	s2 := "[resmgr] [Thread-3] Marking the host d44f268e-2a7c-488e-a6bf-9d76a0fa6905 as not responding"
	ratio, l1, l2, _, _ := TokenRatio(s1, s2)
	t.Log("Token Ratio %d, %d, %d", ratio, l1, l2)
	doTheyMatch, substTokens, canonicalString := DoTheyFuzzyMatch(s1, s2)
	t.Log("Do they match ", doTheyMatch, substTokens, canonicalString)
	parseFileTokenRatio(t, "/Users/roopak/work/pf9-infra/misc/fuzzy-match/token_match_test.txt", s1)

	s1 = "[resmgr1] [Thread-2] Marking the host 87987c71-90f7-4d56-acca-c11adae5c2ed as not responding"
	s2 = "[resmgr2] [Thread-3] Marking the host d44f268e-2a7c-488e-a6bf-9d76a0fa6905 something else that should make it fail more text to see what happens to the mathing algo"
	ratio, l1, l2, _, _ = TokenRatio(s1, s2)
	t.Log("Token Ratio %d", ratio, l1, l2)

	doTheyMatch, substTokens, canonicalString = DoTheyFuzzyMatch(s1, s2)
	t.Log("Do they match ", doTheyMatch, substTokens, canonicalString)

}

func parseFileTokenRatio(t *testing.T, path string, s1 string) {
	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	var data [2]string
	i := 0
	for scanner.Scan() {
		text := scanner.Text()
		data[i%2] = text
		if i%2 != 0 {
			ratio, _, _, _, _ := TokenRatio(data[0], data[1])

			t.Log("Token Ratio %d", ratio)
			doTheyMatch, _, _ := DoTheyFuzzyMatch(data[0], data[1])
			t.Log("Do they match ", doTheyMatch)
		}
		i = i + 1
	}
	ratio, _, _, _, _ := TokenRatio(data[0], s1)
	t.Log("Token ratio %d", ratio)
	t.Log(DoTheyFuzzyMatch(data[0], s1))

}
