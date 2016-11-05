package whistle

import (
	"bufio"
	"os"
	"regexp"
	"testing"
)

func TestRegExpParsing(t *testing.T) {
	InitLog()
	Init()
	basicRegExTest(t)
	parseFile(t, "/Users/roopak/work/pf9-infra/whistle/src/whistle/reg_exp_parse.txt")
	parseFile(t, "/Users/roopak/work/pf9-infra/misc/fuzzy-match/regex-test.txt")
}

func basicRegExTest(t *testing.T) {
	data := `2015-09-25T17:55:01.610724 ERROR_RABBITMQCTL conductor1`

	var patterns = []string{
		`(?P<DATE_TIME>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})[,.](?P<MILLISECONDS>\d{3}) ([0-9]*) (?P<LOG_LEVEL>INFO|TRACE|ERROR|WARNING) (?P<MESSAGE>.*)`,
		`(?P<DATE_TIME>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})[,.](?P<MILLISECONDS>\d{3}) (?P<LOG_LEVEL>INFO|TRACE|ERROR|WARNING) (?P<MESSAGE>.*)`,
		`(?P<DATE_TIME>\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}).(?P<MICROSECONDS>\d{6})(?P<LOG_LEVEL>) (?P<MESSAGE>.*)`,
		`(?P<DATE_TIME>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})[,.](?P<MILLISECONDS>\d{3}) - (?P<PROG>.*) - (?P<LOG_LEVEL>INFO|TRACE|ERROR|WARNING) (?P<MESSAGE>.*)`,
		`(?P<DATE_TIME>\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}) (?P<MESSAGE>.*)`,                         // 2016/01/04 05:19:53 [error] 19920#0: *423945 connect()
		`(?P<DATE_TIME>\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}).(?P<MICROSECONDS>\d{6}) (?P<MESSAGE>.*)`, // 2015-09-25T17:55:01.610724 ERROR_RABBITMQCTL conductor1

	}

	for idx, pattern := range patterns {
		Info.Println("Compiling %d, %s", idx, pattern)
		regEx := regexp.MustCompile(pattern)
		Info.Println("SubExpr Names: %s", regEx.SubexpNames())
		Info.Println("Match %s", regEx.FindAllStringSubmatch(data, -1))
	}

}

func parseFile(t *testing.T, path string) {
	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		text := scanner.Text()
		foundMatch, _ := ParseMessage(text)
		if !foundMatch {
			t.Log("Failed to Parse: ", text)
			t.Fail()
		}

	}
}
