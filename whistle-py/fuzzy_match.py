# Copyright (c) 2015 Platform9 Systems Inc.
import threading
from datacollector import get_data
import dbConnect
from collections import OrderedDict
from log_handler import logging
from token_match import token_ratio
from datetime import datetime
import time, threading

logger = logging.getLogger(__name__)
skip_words = ['DEBUG']

MIN_FUZZ_RATIO = 65
MAX_FUZZ_RATIO = 80

bucket_db = dbConnect.BucketMongoConnect()
logs_db = dbConnect.LogsMongoConnect()


def _do_they_match(key, message):
    ratio, l1, l2 = token_ratio(key, message)
    max_len = max(l1, l2)
    fuzz_ratio = MIN_FUZZ_RATIO
    if max_len > 15:
        fuzz_ratio = MAX_FUZZ_RATIO
    if 15 > max_len and max_len > 4:
        fuzz_ratio = int(1.7 * max_len) + MIN_FUZZ_RATIO

    return fuzz_ratio <= ratio



def _bucketize(buckets, flds):
    '''
    1. Check if the key 'fuzzy matches' to message dictionary
    }
    '''
    did_match = False
    for bucket in buckets:
        key = bucket['message']
        if _do_they_match(key, flds['message']):
            logger.info('Match Found')
            # Update MongoDB
            bucket_db.update_bucket(bucket)
            logs_db.update_log(flds, bucket['_id'])
            did_match = True
            break

    if not did_match:
        logger.info('No Match Found')
        # Implement 2B
        bucket = {'message': flds['message']}
        # Update MongoDB
        doc = bucket_db.update_bucket(bucket)
        bucket['_id'] = doc['_id']
        logs_db.update_log(flds, doc['_id'])

    return buckets



def _ingest():
    while True:
        events = get_data()
        buckets = bucket_db.get_all_buckets()
        mapped_data = events
        '''
        Add check for skipping DEBUG messages
        '''
        if any(sWord in mapped_data['message'] for sWord in skip_words):
            logger.info('Skipping DEBUG/INFO message')
        elif mapped_data['message'].isspace():
            logger.info('Skipping NULL message')
        elif mapped_data['message'] == '':
            logger.info('Skipping empty message')
        else:
            # Add entry to loga
            logger.info('Inserted log entry: ')
            id = logs_db.insert_log(mapped_data)
            mapped_data['_id'] = id
            buckets = _bucketize(buckets, mapped_data)

def recompute():
    logs = logs_db.get_all_logs()
    for log in logs:
        buckets = bucket_db.get_all_buckets()
        _bucketize(buckets, log)


def start_processing():
    t = threading.Thread(target=_ingest)
    t.daemon = True
    t.start()

def start_scrubber():
    logs_db.delete_old_logs()
    threading.Timer(12*60*60, start_scrubber).start()