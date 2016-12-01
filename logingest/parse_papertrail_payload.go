package logingest

import (
	"encoding/json"
	"fmt"
)


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

type PaperTrailPayloadBody struct {
	Payload PayloadStruct  `json:"payload"`
	Meta    []CategoryBody `json:meta`
}

func ParsePapertrailMessage(body []byte) []DataBody {
	var dataBody PaperTrailPayloadBody
	Info.Println("Received request")
	if err := json.Unmarshal(body, &dataBody); err != nil {
		fmt.Println("Error Unmarshaling ", err, body)
	}

	// Post Preprocessing
	out := make([]DataBody, len(dataBody.Payload.Events))
	for idx, event := range dataBody.Payload.Events {

		out[idx].Message = event.Message
		if idx < len(dataBody.Meta) {

			out[idx].Category = dataBody.Meta[idx].Category
		} else {

			out[idx].Category =  ""
		}
	}
	return out
}
