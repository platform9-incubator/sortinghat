package whistle

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
)

var port int
var payloadChannel = make(chan []DataBody, 1024)

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
	router.HandleFunc("/", DefaultPage)
	router.HandleFunc("/data", Ingest)
	router.HandleFunc("/data/recompute", RecomputeBuckets)
	fmt.Println("Listening on  %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
func DefaultPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func Ingest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
        dataBody := ParseSumologic(body)
	payloadChannel <- dataBody
	w.WriteHeader(http.StatusOK)
}

func ingestMessage(chanPayloadBody chan []DataBody) {
	buckets := GetAllBuckets()
	var newBucket *Bucket
	for {
		payloadBody := <-chanPayloadBody
		for _, event := range payloadBody {
			foundMatch, rawLog := ParseMessage(event.Message)
			if foundMatch {
				rawLog.SourceName = event.Host
				rawLog.Category = event.Category
				buckets, rawLog, newBucket = Bucketize(buckets, rawLog)
				rawLog.Id = bson.NewObjectId()
				InsertLog(rawLog)
				if newBucket != nil {
					Info.Println("Upserting bucket")
					UpsertBucket(newBucket)
				}
			} else {
				Error.Println("Can't parse Message, %s \r\n", event)
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
