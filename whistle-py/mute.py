# Copyright (c) 2015 Platform9 Systems Inc.

__author__ = 'roopak'

from dbConnect import MongoConnect
from log_handler import logging
from pymongo import ReturnDocument

logger = logging.getLogger(__name__)


class MuteMongoConnect(MongoConnect):
    def __init__(self):
        super(MuteMongoConnect, self).__init__()
        self._initCollection()

    def _initCollection(self):
        try:
            self._mute_settings = self._db['pf9_mute_settings']
            self._account_mute_settings = self._db['pf9_account_mute_settings']
            logger.info('Connected to MongoDB collections: %s')
        except Exception as e:
            logger.error('Could not connect to: pf9_mute_settings %s', str(e))


    def db_removeMuteSettings(self):
        try:
            self._mute_settings.remove({})
            return True

        except Exception as e:
            logger.error('Error removing collections: ' + e)

    def db_get_mute_settings(self):
        try:
            doc = self._mute_settings.find({})
            return doc
        except Exception as e:
            logger.error('Could not get Mute Settings: ' + e)

    def db_add_mute(self, bucket_id, msg):
        try:
            doc = self._mute_settings.find_one_and_update({'bucket_id': bucket_id},
                                                          {"$set": {'bucket_id': bucket_id, 'msg': msg}},
                                                          projection={'_id': False},
                                                          upsert=True,
                                                          return_document=ReturnDocument.AFTER)
            return doc

        except Exception as e:
            logger.error("Couldn't update mute settings " + e)

    def db_remove_mute(self, bucket_id):
        try:
            self._mute_settings.remove({'bucket_id': bucket_id})
        except Exception as e:
            logger.error("Couldn't update mute settings " + e)

    def db_get_all_account_mute(self):
        try:
            return self._account_mute_settings.find()
        except Exception as e:
            logger.error("Couldn't find account mute " + str(e))

    def db_add_account_mute(self, bucket_id, account_id, msg):
        try:
            doc = self._account_mute_settings.find_one_and_update({'bucket_id': bucket_id,
                                                                   'account_id': account_id},
                                                                  {
                                                                      '$set': {
                                                                          'bucket_id': bucket_id,
                                                                          'account_id': account_id,
                                                                          'msg': msg
                                                                      }
                                                                  },
                                                                  projection={
                                                                      '_id': True
                                                                  },
                                                                  upsert=True,
                                                                  return_document=ReturnDocument.AFTER)
            return doc
        except Exception as e:
            logger.error("Couldn't update account mute settings " + e)

    def db_get_account_mute_by_bucket(self, bucket_id):
        try:
            return self._account_mute_settings.find({'bucket_id': bucket_id})
        except Exception as e:
            logger.error("Couldn't update mute settings " + e)

    def db_remove_account_mute(self, bucket_id, account_id):
        try:
            self._account_mute_settings.remove({'bucket_id': bucket_id, 'account_id': account_id})
        except Exception as e:
            logger.error("Couldn't update mute settings " + e)


mute_db = MuteMongoConnect()


def get_mute_settings():
    return mute_db.db_get_mute_settings()


def mute_bucket(bucket_id, msg):
    return mute_db.db_add_mute(bucket_id, msg)


def remove_mute_bucket(bucket_id):
    return mute_db.db_remove_mute(bucket_id)


def mute_account_bucket(bucket_id, accounts, msg):
    for account in accounts:
        mute_db.db_add_account_mute(bucket_id, account, msg)


def remove_mute_account_bucket(bucket_id, account_id):
    mute_db.remove_account_mute(bucket_id, account_id)
