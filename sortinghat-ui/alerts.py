# Copyright (c) 2015 Platform9 Systems Inc.

__author__ = 'roopak'

import threading
import time
from dbConnect import MongoConnect
from log_handler import logging
import requests
from bson.json_util import dumps,loads

logger = logging.getLogger(__name__)

ALERT_TIMEOUT = 60 * 15
SLACK_URL = "https://hooks.slack.com/services/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"

class AlertMongoConnect(MongoConnect):

    def __init__(self):
        self._collection_name = "pf9_alerts"
        super(AlertMongoConnect, self).__init__()
        self._initCollection()


    def _initCollection(self):
        try:
            self._alert_collection = self._db[self._collection_name]
            logger.info('Connected to MongoDB collections: %s', self._collection_name)
        except Exception as e:
            logger.error('Could not connect to: %s %s', self._collection_name, str(e))

    def get_alert_time(self):
        data = self._alert_collection.find_one({'alert_id':1})
        if data:
            return data['alert_time']
        return 0

    def save_alert_time(self, ts):
        data = self._alert_collection.find_one_and_update({'alert_id':1},
                                                         {"$set":{'alert_id':1, "alert_time": ts}},
                                                         upsert=True)

def send_alert(data):
    payload = {'text': data, "mrkdwn": True}
    r = requests.post(SLACK_URL, data=dumps(payload))

def _format_bucket_data(prefix, bucket_data):
    return ">>> _%s_ *#*: %d  *account:* `%s` \n *message:* _%s_ " % \
           (prefix, bucket_data['category_info']['prod']['count'], bucket_data['category_info']['prod']['hostnames'], bucket_data['message'])

def _get_mute_buckets(mute_settings):
    muted_buckets = set()
    for mute_setting in mute_settings:
        muted_buckets.add(mute_setting['bucket_id'])
    return muted_buckets

def send_alerts(alert_db, last_alert_ts, prefix):
    mute_settings = alert_db.db_get_mute_settings()
    muted_buckets = _get_mute_buckets(mute_settings)
    buckets = alert_db.get_buckets_by_category('prod')

    for k in buckets.keys():
        bucket = buckets[k]
        if not bucket['_id'] in muted_buckets:
            if bucket['ts'] > last_alert_ts :
                send_alert(_format_bucket_data(prefix, bucket))

def alerting_thread():
    """
    Look for last time we sent alerting information

    :return:
    """
    alert_db = AlertMongoConnect()
    count = 0
    while True:
        alert_time = alert_db.get_alert_time()
        current_ts = time.time()
        prefix = "incr"
        if (current_ts - alert_time) >= ALERT_TIMEOUT:
            if count % 4 == 0:
                # Send all the data every 4th iteration
                alert_time = 0
                prefix = "full"

            #send_alerts(alert_db, alert_time, prefix)

        time.sleep(ALERT_TIMEOUT)
        count = count + 1
        alert_db.save_alert_time(current_ts)


def start_alerting():
    t = threading.Thread(target=alerting_thread)
    t.daemon = True
    t.start()
