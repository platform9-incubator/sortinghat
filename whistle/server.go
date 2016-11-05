package whistle

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
)

var port int
var payloadChannel = make(chan PayloadBody, 1024)

func init() {
	flag.IntVar(&port, "port", 8080, "Listen port for Whistle")
	flag.Parse()
}

func StartServer() {
	InitLog()
	Init()
	InitDomain()
	go ingestMessage(payloadChannel)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/data", Ingest)
	router.HandleFunc("/data/recompute", RecomputeBuckets)
	fmt.Println("Listening on  %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}

type CategoryBody struct {
	Category string `json:"category"`
}

type PayloadStruct struct {
	Events []EventBody `json:"events"`
}

type EventBody struct {
	ReceivedAt        string `json:"received_at"`
	Program           string `json:"program"`
	SourceName        string `json:"source_name"`
	Severity          string `json:"severity"`
	Hostname          string `json:"hostname"`
	SourceIp          string `json:"source_ip"`
	DisplayReceivedAt string `json:"display_received_At"`
	Message           string `json:"message"`
	Id                int64  `json:"id"`
	Category          string
}

type PayloadBody struct {
	Payload PayloadStruct  `json:"payload"`
	Meta    []CategoryBody `json:meta`
}

/**
Following JSON is an example of the data received from the website
{
"payload":
{
 "max_id": "559847269952958468",
 "min_id": "559832098555604997",
 "events": [
  {
   "received_at": "2015-07-20T09:59:55-07:00",
   "program": "janitor-daemon.log",
   "source_name": "catfood",
   "severity": "Notice",
   "hostname": "catfood_cat1",
   "facility": "User",
   "source_ip": "52.8.193.220",
   "source_id": 93129174,
   "display_received_at": "Jul 20 09:59:55",
   "message": "2015-07-20 16:59:55,927 - janitor-daemon  - Connection error: ('Connection aborted.', error(101, 'Network is unreachable')), will retry in a bit",
   "id": 559832098555604997
  }
 ],
 "meta" : [{"category": "cat1"}, {"category": "cat2"}]
}
}
*/
func Ingest(w http.ResponseWriter, r *http.Request) {
	var dataBody PayloadBody
	Info.Println("Received request")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &dataBody); err != nil {
		fmt.Println("Error Unmarshaling ", err, body)
	}
	// Preprocessing
	for idx, _ := range dataBody.Payload.Events {
		if idx < len(dataBody.Meta) {
			dataBody.Payload.Events[idx].Category = dataBody.Meta[idx].Category
		} else {
			dataBody.Payload.Events[idx].Category = ""
		}
	}

	payloadChannel <- dataBody
	w.WriteHeader(http.StatusOK)
}

func ingestMessage(chanPayloadBody chan PayloadBody) {
	buckets := GetAllBuckets()
	var newBucket *Bucket
	for {
		payloadBody := <-chanPayloadBody
		for _, event := range payloadBody.Payload.Events {
			foundMatch, rawLog := ParseMessage(event.Message)
			if foundMatch {
				rawLog.SourceName = event.SourceName
				rawLog.Category = event.Category
				buckets, rawLog, newBucket = Bucketize(buckets, rawLog)
				rawLog.Id = bson.NewObjectId()
				InsertLog(rawLog)
				if newBucket != nil {
					Info.Println("Upserting bucket")
					UpsertBucket(newBucket)
				}
			} else {
				Error.Println("Can't parse Message, %s \r\n", event.Message)
			}
		}
	}
}

func RecomputeBuckets(w http.ResponseWriter, r *http.Request) {
	Info.Println("Recompute request received")
	buckets := GetAllBuckets()
	Info.Println(" Number of existing buckets ", len(buckets))
	logs := GetAllLogs()
	Info.Println(" Number of Logs ", len(logs))
	var newBucket *Bucket
	for _, log := range logs {
		buckets, _, newBucket = Bucketize(buckets, &log)
		UpdateLog(&log)
		if newBucket != nil {
			Info.Println("-- Upserting new bucket -- ", newBucket.Message)

			UpsertBucket(newBucket)
		}
	}
	Info.Println("Done Processing  ")
	w.WriteHeader(http.StatusOK)
}
