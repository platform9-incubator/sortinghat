#!/usr/bin/env python
# Copyright (c) 2015 Platform9 Systems Inc.
from flask import Flask,redirect, abort, request, jsonify
import json
from alerts import start_alerting
from fuzzy_match import start_processing, start_scrubber
from bucket_log_reports import get_result, get_host_result, get_category_result, update_bucket_message
from bucket_log_reports import get_bucket_info, reset_buckets_mutes, get_bucket_details
from mute import get_mute_settings, mute_bucket, remove_mute_bucket, mute_account_bucket, remove_mute_account_bucket
from datacollector import data_sink
from log_handler import logging
from bson.json_util import dumps,loads

APP_NAME = "Whistle"

logger = logging.getLogger(__name__)

start_processing()
start_alerting()
start_scrubber()

app = Flask(APP_NAME)
app.debug = True


@app.route('/')
def incoming():
    return redirect("/sortinghat-ui/static/index.html", code=302)

@app.route('/data', methods=['POST'])
def ingest():
    data = request.get_json()
    try:
        for i in range(0,len(data['payload']['events'])):
            try:
                metaTag = data['meta'][i]
            except:
                logger.error('Invalid MetaTag Entry ', data['meta'])
                metaTag = ''

            data['payload']['events'][i]['meta'] = metaTag
            data_sink(data['payload']['events'][i])
    except:
        logger.error('Invalid Log Entry: ',data)

    return ""


@app.route('/alerts', methods=['GET'])
def get_alerts():
    result = get_result()
    return dumps(result)


@app.route('/host/<hostname>', methods=['GET'])
def get_host_alerts(hostname):
    result = get_host_result(hostname)
    return dumps(result)


@app.route('/category/<category>', methods=['GET'])
def get_category_alerts(category):
    result = get_category_result(category)
    return dumps(result)


@app.route('/alerts', methods=['POST'])
def rebuild_alerts():
    return dumps(reset_buckets_mutes())


@app.route('/bucket/details/<bucket_id>', methods=['GET'])
def get_bucket(bucket_id):
    return dumps(get_bucket_details(loads(bucket_id)))


@app.route('/bucket/summary/<bucket_id>', methods=['GET'])
def get_bucket_summary(bucket_id):
    return dumps(get_bucket_info(loads(bucket_id)))

@app.route('/bucket/<bucket_id>/<msg>', methods=['POST'])
def add_message(bucket_id, msg):
    return dumps(update_bucket_message(loads(bucket_id), msg))

@app.route('/settings/mute', methods=['GET'])
def mute_settings():
    return dumps(get_mute_settings())


@app.route('/mute/<bucket_id>/<msg>', methods=['POST'])
def add_mute_settings(bucket_id, msg):
    return dumps(mute_bucket(loads(bucket_id), msg))


@app.route('/unmute/<bucket_id>', methods=['POST'])
def umute(bucket_id):
    return dumps(remove_mute_bucket(loads(bucket_id)))


@app.route('/mute-account-bucket/<bucket_id>', methods=['POST'])
def add_account_mute_settings(bucket_id):
    data = request.get_json()
    return dumps(mute_account_bucket(loads(bucket_id), data['accounts'], data['msg']))


@app.route('/unmute-account-bucket/<bucket_id>/<account_id>/', methods=['POST'])
def unmute_account_mute_settings(bucket_id, account_id):
    return dumps(remove_mute_account_bucket(bucket_id, account_id))

if __name__ == '__main__':
    app.run(host="0.0.0.0", port=8080)
