package whistle

import (
	"bytes"
	"errors"
	"log"
	"regexp"
	"strconv"
	"time"
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
**/

var messageRegexes []regexp.Regexp
var logger *log.Logger

func Init() {
	var buf bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
	var patterns = []string{
		`(?P<DATE_TIME>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})[,.](?P<MILLISECONDS>\d{3}) ([0-9]*) (?P<LOG_LEVEL>INFO|TRACE|ERROR|WARNING) (?P<MESSAGE>.*)`,
		`(?P<DATE_TIME>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})[,.](?P<MILLISECONDS>\d{3}) (?P<LOG_LEVEL>INFO|TRACE|ERROR|WARNING) (?P<MESSAGE>.*)`,
		`(?P<DATE_TIME>\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}).(?P<MICROSECONDS>\d{6})(?P<LOG_LEVEL>) (?P<MESSAGE>.*)`,
		`(?P<DATE_TIME>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})[,.](?P<MILLISECONDS>\d{3}) - (?P<PROG>.*) - (?P<LOG_LEVEL>INFO|TRACE|ERROR|WARNING) (?P<MESSAGE>.*)`,
		`(?P<DATE_TIME>\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}) (?P<MESSAGE>.*)`,                         // 2016/01/04 05:19:53 [error] 19920#0: *423945 connect()
		`(?P<DATE_TIME>\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}).(?P<MICROSECONDS>\d{6}) (?P<MESSAGE>.*)`, // 2015-09-25T17:55:01.610724 ERROR_RABBITMQCTL conductor1

	}

	messageRegexes = make([]regexp.Regexp, len(patterns))
	for idx, pattern := range patterns {
		Info.Println("Compiling %s", pattern)
		var regEx *regexp.Regexp = regexp.MustCompile(pattern)
		if regEx != nil {
			messageRegexes[idx] = *regEx
		}
	}
}

func matchGroups(matcher regexp.Regexp, data string) (map[string]string, error) {

	result := map[string]string{}

	groupNames := matcher.SubexpNames()
	allMatches := matcher.FindAllStringSubmatch(data, -1)
	if allMatches == nil {
		return result, errors.New("No match found")
	}
	matchResult := allMatches[0]
	if matchResult != nil {
		for i, value := range matchResult {
			result[groupNames[i]] = value
		}
	}
	return result, nil
}

func parseParts(regExp regexp.Regexp, data string) (string, int64, string, string, error) {

	matchResult, err := matchGroups(regExp, data)
	if matchResult != nil {

		logLevel := matchResult["LOG_LEVEL"]
		dateStr := matchResult["DATE_TIME"]
		microseconds := ""
		if millisec_str, ok := matchResult["MILLISECONDS"]; ok {
			milliseconds, _ := strconv.Atoi(millisec_str)
			microseconds = strconv.FormatInt(int64(1000*milliseconds), 10)
		} else {
			microseconds = matchResult["MICROSECONDS"]
		}

		dateTimeLayouts := []string{
			//	Layouts must use the reference time Mon Jan 2 15:04:05 MST 2006 to show the pattern with which to format/parse a given
			`2006-01-02 15:04:05.000000`,
			`2006-01-02T15:04:05.000000`,
			`2006/01/02 15:04:05.000000`,
		}
		var microSec int64 = 0
		timeStr := ""
		dateTimeStr := dateStr + "." + microseconds
		for _, dateTimeLayout := range dateTimeLayouts {
			timeObj, err := time.Parse(dateTimeLayout, dateTimeStr)
			if err != nil {
				continue
			}
			microSec = timeObj.UnixNano() / 1000
			timeStr = timeObj.Format("2006-01-02 15:04:05.000000")
		}
		return logLevel, microSec, timeStr, matchResult["MESSAGE"], err
	} else {
		timeObj := time.Now()
		microSec := timeObj.UnixNano() / 1000
		timeStr := timeObj.Format("2006-01-02 15:04:05.000000")
		return "INFO", microSec, timeStr, "", err
	}
}

func ParseMessage(data string) (bool, *RawLog) {
	var rawLog RawLog
	foundMatch := false
	for _, matcher := range messageRegexes {
		logLevel, microsec, timeStr, message, err := parseParts(matcher, data)

		if err != nil {
			continue
		}
		rawLog.Severity = logLevel
		rawLog.Timestamp = microsec
		rawLog.TimeStr = timeStr
		rawLog.Message = message
		foundMatch = true
		break
	}
	return foundMatch, &rawLog
}
