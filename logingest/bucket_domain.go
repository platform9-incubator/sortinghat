package logingest

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"bitbucket.org/platform9/sortinghat/vendor/gopkg.in/mgo.v2"
)

type Bucket struct {
	Id               bson.ObjectId `bson:"_id"`
	Message          string        `bson:"message"`
	CanonicalMessage string        `bson:"canonical_message"`
}

type RawLog struct {
	Id            bson.ObjectId `bson:"_id"`
	Severity      string        `bson:"severity"`
	Message       string        `bson:"message"`
	BucketId      bson.ObjectId `bson:"bucket_id"`
	Timestamp     int64         `bson:"timestamp"`
	TimeStr       string        `bson:"time_str"`
	SourceName    string        `bson:"source_name"`
	Category      string        `bson:"category"`
	MessageTokens []string      `bson:"message_tokens"`
}

type BucketDomain struct {
	buckets *mgo.Collection
	rawLogs *mgo.Collection
}

var bucketManager BucketDomain

func InitDomain() {
	session, _ := mgo.Dial("mongo:27017")
	bucketManager.buckets = session.DB("pf9_logs").C("pf9_bucket_logs")
	bucketManager.rawLogs = session.DB("pf9_logs").C("pf9_raw_logs")
}

func GetAllBuckets() []Bucket {
	var buckets []Bucket
	bucketManager.buckets.Find(nil).All(&buckets)
	return buckets
}

func UpsertBucket(bucket *Bucket) {
	bucketManager.buckets.Upsert(bson.M{"_id": bucket.Id}, bucket)
}

func InsertLog(rawLog *RawLog) {
	err := bucketManager.rawLogs.Insert(rawLog)
	if err != nil {
		Error.Fatal("Error saving rawLog", err)
	}
}

func UpdateLog(rawLog *RawLog) {
	bucketManager.rawLogs.Update(bson.M{"_id": rawLog.Id}, rawLog)
}

func GetAllUnprocessedLogs() []RawLog {
	var rawLogs []RawLog
	bucketManager.rawLogs.Find(bson.M{"bucket_id": bson.M{"$exits": "0"}}).All(&rawLogs)
	return rawLogs
}

func GetAllLogs() []RawLog {
	var rawLogs []RawLog
	bucketManager.rawLogs.Find(nil).All(&rawLogs)
	return rawLogs
}
