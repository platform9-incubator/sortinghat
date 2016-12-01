package logingest
import "testing"

func TestBucketize(t *testing.T) {
	InitDomain()

	buckets := GetAllBuckets()
	allLogs := GetAllLogs()
	t.Logf("Number of Logs %d", len(allLogs))
	for _, log := range allLogs {
		buckets, _, _ = Bucketize(buckets, &log)
	}

	t.Logf("Buckets --%d--- \r\n", len(buckets))

	for _, bucket := range buckets {
		t.Logf("Bucket %s \r\n", bucket)
	}

}
