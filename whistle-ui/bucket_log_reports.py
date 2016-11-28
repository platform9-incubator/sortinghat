# Copyright (c) 2015 Platform9 Systems Inc.

__author__ = 'roopak'

import dbConnect
from mute import MuteMongoConnect
from collections import OrderedDict
from log_handler import logging
from token_match import token_ratio
from datetime import datetime
from fuzzy_match import recompute

logger = logging.getLogger(__name__)

bucket_db = dbConnect.BucketMongoConnect()
logs_db = dbConnect.LogsMongoConnect()
mute_db = MuteMongoConnect()


def get_result(category=None):
    logger.info("Called get result")
    mute_buckets = mute_db.db_get_mute_settings()
    mute_account_buckets = mute_db.db_get_all_account_mute()
    aggregate_result = logs_db.get_all_bucket_aggregate_info(category, mute_buckets, mute_account_buckets, None, None)
    buckets_arr = bucket_db.get_all_buckets()
    bucket_dict = {}
    for bucket in buckets_arr:
        canonical_message = bucket['message']
        user_message = bucket['message']
        if 'canonical_message' in bucket:
           canonical_message = bucket['canonical_message']
        if 'user_message' in bucket:
           user_message = bucket['user_message']
        bucket_dict[bucket['_id']] = {'canonical_message': canonical_message, 'user_message': user_message}

    # Add 'user_message' or 'canonical_message' to the aggregate information.
    for r in aggregate_result:
        if 'canonical_message' in bucket_dict[r['_id']]:
            r['canonical_message'] = bucket_dict[r['_id']]['canonical_message']
        if 'user_message' in bucket_dict[r['_id']]:
            r['user_message'] = bucket_dict[r['_id']]['user_message']

    return aggregate_result


# Get host specific results
def get_host_result(hostname):
    logger.info("get_host_result")
    return logs_db.get_host_logs(hostname)

def get_category_result(category):
    logger.info("get category result")
    return get_result(category)

def reset_buckets_mutes():
    logger.info("reset bucket mutes")
    recompute()
    return get_result()

def update_bucket_message(bucket_id, msg):
    logger.info("Adding Bucket info")
    bucket = {'_id': bucket_id, 'user_messagex': msg}
    return bucket_db.update_bucket(bucket)

def get_bucket_info(bucket_id):
    logger.info("get bucket info")
    return bucket_db.get_bucket(bucket_id)

def get_bucket_details(bucket_id):
    logger.info("Request for bucket_id %s" % str(bucket_id))

    mute_account_buckets = mute_db.db_get_account_mute_by_bucket(bucket_id)
    muted_accounts = map(lambda x:  x['account_id'], mute_account_buckets)

    bucket_info = bucket_db.get_bucket(bucket_id)
    bucket_aggregated_info = logs_db.get_all_bucket_aggregate_info(None, None,mute_account_buckets, None, bucket_id)
    bucket_messages = logs_db.get_logs_for_bucket(bucket_id, mute_account_buckets)
    ret = {'bucket': bucket_info, 'bucket_aggregate': bucket_aggregated_info, 'muted_accounts': muted_accounts, 'bucket_messages': bucket_messages}
    return ret
