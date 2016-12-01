package logingest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
//	"strconv"
//	"strings"
//	"strconv"
//	"strings"
)

func TestJsonParsing(t *testing.T) {
	//testParse1(t)
	testParse2(t)
}

func testParse1(t *testing.T) {
	file, err := ioutil.ReadFile("/Users/roopak/work/pf9-infra/whistle/src/whistle/papertrail_parse.txt")
	if err != nil {
		fmt.Println("Error reading file")
	}
	var dataBody PaperTrailPayloadBody

	if err = json.Unmarshal(file, &dataBody); err != nil {
		fmt.Println("Marshalled failed", err)
	} else {
		fmt.Println("Marshalled  ", dataBody)
	}

	fmt.Println("Event Legnth", len(dataBody.Payload.Events))
	fmt.Println("Metadata Length", len(dataBody.Meta))
	fmt.Println(dataBody.Meta[0].Category)
	if len(dataBody.Payload.Events) != len(dataBody.Meta) {
		t.Fail()
	}

}
/*
var parse_state int = 0
func parseNext(r rune) rune {
	switch {
		case parse_state == 0  && r == '\\':
			parse_state = 1
			return -1
		case parse_state == 1 && r == '"':
			parse_state = 0
		        return '"'
		case parse_state == 1 && r == '\\':
			parse_state = 0
		        return '\\'
		case parse_state == 1 && r == '/':
			parse_state = 0
		        return '/'
		case parse_state == 1 :
			parse_state = 0
			return "\\"+r
		}
	return r
}
*/
func testParse2(t *testing.T) {
	file, err := ioutil.ReadFile("/Users/roopak/work/src/bitbucket.org/platform9/whistle-log/whistle/sumologic_parse4.txt")
	if err != nil {
		fmt.Println("Error reading file")
		t.Errorf("Error reading file")
	}
	dataBody := ParseSumologic(file)

	fmt.Println("Event Length", len(dataBody))
}
