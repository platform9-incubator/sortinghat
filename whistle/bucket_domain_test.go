package whistle

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestBucketDomain(t *testing.T) {
	InitDomain()
	buckets := GetAllBuckets()
	for _, bucket := range buckets {
		t.Logf("Bucket %s \r\n", bucket)
	}

	unprocessLogs := GetAllUnprocessedLogs()
	for _, log := range unprocessLogs {
		t.Logf("Log %s \r\n", log)
	}
	testMarshal(t)
}

func testMarshal(t *testing.T) {
	var rawLog RawLog
	rawLog.SourceName = "sourceName"
	rawLog.Category = "prod"
	rawLog.Message = " Host X not responding"
	rawLog.TimeStr = "2015-12-05 18:16:32.203"

	out, err := bson.Marshal(&rawLog)
	if err != nil {
		fmt.Println("Error Marshaling: %s", err.Error())
		t.Fail()
	}
	fmt.Println("Test Marshal: ", string(out))
}
