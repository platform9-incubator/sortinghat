package whistle

import (
	"encoding/json"
	"fmt"
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
	Data    []DataBody `json:data`
}


func ParseSumologic(body []byte)  []DataBody {
	var dataBody SumologicBody
	Info.Println("Received request")
	if err := json.Unmarshal(body, &dataBody); err != nil {
		fmt.Println("Error Unmarshaling ", err, body)
	}

	return dataBody.Data;
}

