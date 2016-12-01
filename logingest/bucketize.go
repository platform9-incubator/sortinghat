package logingest
import "gopkg.in/mgo.v2/bson"

func Bucketize(buckets []Bucket, rawLog *RawLog) ([]Bucket, *RawLog, *Bucket) {
	message := rawLog.Message
	var newBucket Bucket
	var pNewBucket *Bucket = nil
	var foundMatch = false
	for _, bucket := range buckets {
		match, substTokens, canonicalString := DoTheyFuzzyMatch(bucket.Message, message)
		if match {
			// associated bucket's id with the message
			// fmt.Println("Found Match")
			foundMatch = true
			rawLog.BucketId = bucket.Id
			rawLog.MessageTokens = substTokens
			if bucket.CanonicalMessage != canonicalString {
				bucket.CanonicalMessage = canonicalString
				pNewBucket = &bucket
			}
			break
		}
	}
	if !foundMatch {
		newBucket.Message = message
		newBucket.Id = bson.NewObjectId()
		rawLog.BucketId = newBucket.Id
		buckets = append(buckets, newBucket)
		pNewBucket = &newBucket
	}
	return buckets, rawLog, pNewBucket

}
