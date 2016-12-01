# Copyright (c) 2015 Platform9 Systems Inc.
# This provides parsing of the incoming messages and
# extracting interesting information like date/time,
__author__ = 'roopak'

from Queue import Queue
from pymongo import MongoClient
import re
from datetime import datetime, date
from log_handler import logging

from bson.json_util import dumps,loads

"""
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
"""
logger = logging.getLogger(__name__)
_queue = Queue()

patterns = ['(?P<DATE_TIME>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})[,.](?P<MILLISECONDS>\d{3}) (?P<LOG_LEVEL>INFO|TRACE|ERROR|WARNING)(?P<MESSAGE>.*)',
            '(?P<DATE_TIME>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})[,.](?P<MILLISECONDS>\d{3}) ([0-9]*) (?P<LOG_LEVEL>INFO|TRACE|ERROR|WARNING)(?P<MESSAGE>.*)',
            '(?P<DATE_TIME>\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}).(?P<MICROSECONDS>\d{6})(?P<LOG_LEVEL>)(?P<MESSAGE>.*)',
            #2015-12-15 04:47:05,211 - janitor-daemon - ERROR]
            '(?P<DATE_TIME>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})[,.](?P<MILLISECONDS>\d{3}) - (?P<PROG>.*) - (?P<LOG_LEVEL>INFO|TRACE|ERROR|WARNING)(?P<MESSAGE>.*)']

date_time_fmts = ['%Y-%m-%d %H:%M:%S.%f', '%Y-%m-%d %H:%M:%S.%f', '%Y-%m-%dT%H:%M:%S.%f', '%Y-%m-%d %H:%M:%S.%f']
regexps = map(lambda x: re.compile(x), patterns)

# Received TS format
# "received_at": "2015-07-20T09:59:55-07:00"
rcvd_time_regex = re.compile('(?P<DATE_TIME>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})(?P<TZ_HR>[+-]\d{2}):(?P<TZ_MIN>\d{2})')
rcvd_date_time_fmt = '%Y-%m-%dT%H:%M%S%z'

def _get_parts(match, fmt):
    groups = match.groupdict()
    log_level = groups["LOG_LEVEL"]
    date_str = groups["DATE_TIME"]
    microseconds = ""
    if (groups.has_key("MILLISECONDS")):
        microseconds = str(int(groups["MILLISECONDS"]) * 1000)
    else:
        microseconds = groups["MICROSECONDS"]

    date_time_str = date_str+"."+microseconds
    date_obj = datetime.strptime(date_time_str, fmt)
    return log_level, date_obj, groups["MESSAGE"]

def _parse_rcvd_time(ts):
    match = rcvd_time_regex.match(ts)
    if match:
        date_str = "%s%s%s" % (match.group('DATE_TIME'), match.group('TZ_HR'), match.group('TZ_MIN'))
        return datetime.strptime(date_str, rcvd_date_time_fmt)
    # TODO: print error
    return datetime.utcnow()

def _convert_timestamp(ms_str):
    date_obj = datetime.utcfromtimestamp(float(ms_str))
    return date_obj.strftime('%Y-%m-%d %H:%M:%S.%f')


def _to_timestamp(dt, epoch=datetime(1970, 1, 1)):
    td = dt - epoch
    return td.total_seconds()


def _extract_event_data(data):
    try:
        timestamp = _to_timestamp(data['received_at'])
        ret = {"time_str": _convert_timestamp(timestamp),
               "timestamp": timestamp,
               "program": data['program'],
               "source_name": data['source_name'],
               "severity": data['severity'],
               "hostname": data['hostname'],
               "facility": data['facility'],
               "source_ip": data['source_ip'],
               "source_id": data['source_id'],
               "message": data['message'],
               "subCategory": data['meta']['category'],
               "mid": data['id']}
        return ret

    except Exception as e:
        logger.error('Invalid Log entry : ' + str(e))


def _parse_message(data):

    message_match = False
    for (regexp, fmt) in zip(regexps, date_time_fmts):
        match = regexp.match(data['message'])
        if match:
            log_level, date_obj, message = _get_parts(match, fmt)
            data['severity'] = log_level
            data['received_at'] = date_obj
            data['message'] = message
            message_match = True

    if not message_match:
        data['received_at'] = _parse_rcvd_time(data['received_at'])

    return _extract_event_data(data)

def data_sink(data):
    try:
        data = _parse_message(data)
        _queue.put(data)
    except Exception, e:
        logger.info("Data %s " % str(data))
        logger.error("Exception %s " % str(e))

def get_data():
   return _queue.get()


def _print_grps(match, fmt):
    print "Group --"
    log_level, date_obj, message = _get_parts(match, fmt)
    print "Log level is %s " % (log_level)
    print "Date is %s" % str(date_obj)
    print "Message is %s" % message


def _test2():
    f = open('./input-test.txt')
    for l in f.readlines():
        l = l.strip()
        print _parse_message(loads(l.replace("\'", "\"")))


if __name__ == "__main__":
    _test2()
