package whistle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestJsonParsing(t *testing.T) {
	testParse(t)
}

func testParse(t *testing.T) {
	file, err := ioutil.ReadFile("/Users/roopak/work/pf9-infra/whistle/src/whistle/json_parse.txt")
	if err != nil {
		fmt.Println("Error reading file")
	}
	var dataBody PayloadBody

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
