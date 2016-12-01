package whistle

import (
	"encoding/json"
	"fmt"
	"strings"
)


/**
Following JSON is an example of the data received from the website
{
  "data": [
    {
      "Message": "2016-11-23 23:26:25.197 16893 ERROR oslo.messaging._drivers.impl_rabbit [req-f1b416e1-f788-48e4-bd6a-7e6334cc727c - - - - -] AMQP server 127.0.0.1:5672 closed the connection. Check login credentials: Socket closed",
      "Time": 1479943585197,
      "Host": "pf9-sm-lab-vmwest1.platform9.net",
      "Category": "prod/openstack/nova",
      "Name": "/var/log/nova/nova-consoleauth.log",
      "Collector": "pf9-sm-lab-vmwest1.platform9.net"
    }
  ]
}
*/


type SumologicBody struct {
	Data []DataBody `json:data`
}

func coalesceStackTraces(sumologicBody SumologicBody) []DataBody {
	// host2log is a map of
	// hostname : { log_name: [timestamp, index-of-the-last-Data object in databody]}
	type TimeIdx struct {
		time int
		idx  int
	}
	var host2log map[string]map[string]TimeIdx
	host2log = make(map[string]map[string]TimeIdx)

	var ret = make([]DataBody, len(sumologicBody.Data))
	fmt.Println("Number of incoming event", len(sumologicBody.Data))
	index := 0
	for _, data := range sumologicBody.Data {
		if _, hostExists := host2log[data.Host] ; hostExists {
			if tidx, logExists := host2log[data.Host][data.Name] ; logExists {
				if (tidx.time == data.Time) {
					ret[tidx.idx].Message = fmt.Sprintf("%s\n%s", ret[tidx.idx].Message, data.Message)
					continue
				}
			}
		} else {
			host2log[data.Host] = map[string]TimeIdx {}
		}
		ret[index] = data
		var tidx TimeIdx
		tidx.idx = index
		tidx.time = data.Time
		host2log[data.Host][data.Name] = tidx
		index = index + 1
	}
	fmt.Println("Number of coalesced events", index)
	return ret[0:index]
}

func ParseSumologic(body []byte) []DataBody {

	var dataBody SumologicBody
	//Info.Println("Received request")
	var data string
	data = strings.Replace(string(body), "\\\"", "\"", -1)
	data = strings.Replace(data, "\\/", "/", -1)
	data = strings.Replace(data, "\\\\", "\\", -1)
	if err := json.Unmarshal([]byte(data), &dataBody); err != nil {
		fmt.Println("Error Unmarshaling ", err, body)
	}

	return coalesceStackTraces(dataBody);
}

