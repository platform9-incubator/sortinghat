package whistle

import (
	"bufio"
	"os"
	"testing"
)

func TestTokenRatio(t *testing.T) {
	data := [] string{
		"cinder.api.middleware.fault [req-d344e75b-8636-4e64-a719-1cfb47c5e3ac 7a365eed6f484dc888bcfda1196dc199 5d1af3fdd63c4eb8a33e9bea20ec35f0 - - -] Caught error: (_mysql_exceptions.OperationalError) (1040, 'Too many connections')",
		"cinder.api.middleware.fault [req-e5c8ef2c-5379-4eab-a821-b2616e162a8e 5a1151856fcb47bea65894be590316a9 fa2524ff80f0476ebd4055a905f06e42 - - -] Caught error: (_mysql_exceptions.OperationalError) (1040, 'Too many connections')",

		"[resmgr] [Thread-2] Marking the host 87987c71-90f7-4d56-acca-c11adae5c2ed as not responding",
		"[resmgr] [Thread-3] Marking the host d44f268e-2a7c-488e-a6bf-9d76a0fa6905 as not responding",

		"[resmgr1] [Thread-2] Marking the host 87987c71-90f7-4d56-acca-c11adae5c2ed as not responding",
		"[resmgr2] [Thread-3] Marking the host d44f268e-2a7c-488e-a6bf-9d76a0fa6905 something else that should make it fail more text to see what happens to the mathing algo",
	}

	for i := 0; i < len(data); i=i+2 {
		s1 := data[i]
		s2 := data[i+1]
		doMatchTest(t, s1, s2)
	}
}

func doMatchTest(t *testing.T, s1 string, s2 string) {
	ratio, l1, l2, _, _ := TokenRatio(s1, s2)
	t.Log("S1:", s1)
	t.Log("S2:", s2)
	t.Log("Token Ratio %d, %d, %d", ratio, l1, l2)
	doTheyMatch, substTokens, canonicalString := DoTheyFuzzyMatch(s1, s2)
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
