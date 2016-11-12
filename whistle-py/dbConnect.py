# Copyright (c) 2015 Platform9 Systems Inc.
from pymongo import MongoClient
from datetime import date, datetime, tzinfo, timedelta
import pytz
from pymongo import ReturnDocument
from log_handler import logging

logger = logging.getLogger(__name__)

class MongoConnect(object):
    def __init__(self):
        self.__client = None
        self._db = None
        self._dbRebuild = False
        self.db_connect()
        self.db_initDB()
        self.db_initCollection()

    def db_connect(self, server='mongo', port=27017):
        try:
            self.__client = MongoClient(server, port)
            logger.info('Connected to MongoDB')
        except:
            logger.error('Could not connect to MongoDB server')

    def db_stop(self):
        try:
            self.__client.close()
            logger.info('MongoDB connection closed')
        except:
            logger.error('Could not close the MongoDB connection')

    def db_initDB(self, dbName='pf9_logs'):
        try:
            self._db = self.__client[dbName]
            logger.info('Connected to MongoDB database: ' + dbName)
        except:
            logger.error('Could not create the MongoDB Database')

    def db_initCollection(self):
        try:
            logger.info('Connected to MongoDB collections:')
        except:
            logger.error('Could not connect to mongo db')


    def get_epoch(self, dt):
        try:
            seconds = (dt - datetime(1970, 1, 1)).total_seconds()
            return seconds
        except Exception as e:
            logger.error('Could not convert to epoch: ' + e)




class BucketMongoConnect(MongoConnect):
    def __init__(self):
        self._collection_name = "pf9_bucket_logs"
        super(BucketMongoConnect, self).__init__()
        self._initCollection()

    def _initCollection(self):
        try:
            self._bucket_collection = self._db[self._collection_name]
            logger.info('Connected to MongoDB collections: %s', self._collection_name)
        except Exception as e:
            logger.error('Could not connect to: %s %s', self._collection_name, str(e))


    def remove_all_buckets(self):
        try:
            self._bucket_collection.remove({})
            return True

        except Exception as e:
            logger.error('Error removing collections: ' + e)

    def get_all_buckets(self):
        try:
            return self._bucket_collection.find({})
        except Exception as e:
            logger.error('Cannot get buckets')

    def get_bucket(self, bucket_id):
        try:
            return self._bucket_collection.find({'_id': bucket_id})
        except Exception as e:
            logger.error('Cannot find bucket')

    def update_bucket(self, data):
        '''
        check if there is existing entry with same message id
        if entry
            update the count and host filed [ last affected host will be the last one in the list
            [ list.append() default behavior ]]
        else
            add the new category

        return [ either category updated or new category added ]
        '''
        try:
            query = {'message': data['message']}
            if data.has_key('_id'):
                query = {'_id': data['_id']}
            doc = self._bucket_collection.find_one_and_update(query,
                                                         {"$set":data},
                                                         projection={'_id': True},
                                                         upsert=True,
                                                         return_document=ReturnDocument.AFTER)
            return doc

        except Exception as e:
            logger.error('Error updating category: ' + e)

    def reset_all_bucket_properties(self, properties):
        try:
            doc = self._bucket_collection.update({},{"$set": properties})
            return doc

        except Exception as e:
            logger.error('Error updating category: ' + e)


    def delete_bucket(self, data):
        try:
            return self._bucket_collection.delete_one({'message': data['message']})
        except Exception as e:
            logger.error('Unable to delete the category: ' + e)


    def get_buckets_by_category(self, sub_category):
        try:
            return self._bucket_collection.find({'categories': sub_category})
        except Exception as e:
            logger.error('Could not generate category: ' + str(e))

    def get_buckets_by_category_and_time(self, sub_category, min_timestamp):
        try:
            return self._bucket_collection.find({'categories': sub_category, 'ts':{'$gt', min_timestamp}})
        except Exception as e:
            logger.error('Could not generate category: ' + str(e))




class LogsMongoConnect(MongoConnect):
    def __init__(self):
        self._collection_name = "pf9_raw_logs"
        super(LogsMongoConnect, self).__init__()
        self._initCollection()

    def _initCollection(self):
        try:
            self._log_collection = self._db[self._collection_name]
            logger.info('Connected to MongoDB collections: %s', self._collection_name)
        except Exception as e:
            logger.error('Could not connect to: %s %s', self._collection_name, str(e))

    def get_all_bucket_aggregate_info(self, category, mute_buckets, mute_account_buckets, lowest_timestamp, bucket_id):
        try:
            match_expr_subcategory = None
            match_bucket_expr = None
            if bucket_id:
                match_bucket_expr = {
                    '$match': {
                        'bucket_id': bucket_id
                    }
                }
            if category:
                match_expr_subcategory = {
                    '$match': {
                        'category': category
                    }
                }
            match_expr_time = None

            if lowest_timestamp:
                match_expr_time = {
                    '$match': {
                        'timestamp': {'$gte': lowest_timestamp}
                    }
                }
            match_expr_mute_buckets = None
            if mute_buckets:
                mute_bucket_ids = map(lambda x: x['bucket_id'], mute_buckets)
                if mute_bucket_ids and len(mute_bucket_ids) > 0:
                    match_expr_mute_buckets = {
                        '$match': {
                            'bucket_id': {
                                '$nin': mute_bucket_ids
                            }
                        }
                    }

            match_expr_mute_account_buckets = None

            if mute_account_buckets :
                mute_account_buckets = map(lambda x: {'bucket_id': x['bucket_id'], 'source_name': x['account_id']},
                                           mute_account_buckets)
                if mute_account_buckets and len(mute_account_buckets) > 0:
                    match_expr_mute_account_buckets = {
                        '$match': {
                            '$nor': mute_account_buckets
                        }
                    }

            group_stage1 = {
                '$group': {
                    '_id': {
                        'bucket_id': '$bucket_id',
                        'source_name': '$source_name'
                    },
                    'source_count': {
                        '$sum': 1
                    },
                    'severity' : {
                      '$first' : '$severity'
                    },
                    'message': {
                        '$first': '$message'
                    },
                    'timestamp': {
                        '$max': '$timestamp'
                    },
                    'time_str': {
                        '$max': '$time_str'
                    }
                }
            }

            group_stage2 = {
                '$group': {
                    '_id': '$_id.bucket_id',
                    'source_names': {
                      '$push': '$_id.source_name'
                    },
                    'sources': {'$push':{
                        'source_name': '$_id.source_name',
                        'source_count': '$source_count'
                    }},
                    'total_count': {
                        '$sum': '$source_count'
                    },
                    'severity' : {
                      '$first' : '$severity'
                    },
                    'message': {
                        '$first': '$message'
                    },
                    'timestamp': {
                        '$max': '$timestamp'
                    },
                    'time_str': {
                        '$max': '$time_str'
                    }
                }
            }

            sort_stage = {
                '$sort' : {'total_count' : -1}
            }
            pipeline = filter(lambda x: x != None,
                              [match_bucket_expr, match_expr_time, match_expr_subcategory, match_expr_mute_buckets, match_expr_mute_account_buckets, group_stage1,
                               group_stage2, sort_stage])
            return self._log_collection.aggregate(pipeline=pipeline)
        except Exception as e:
            logger.error('Could not connect to %s %s', self._collection_name, str(e))

    def get_logs_for_bucket(self, bucket_id, mute_account_buckets):
        try:
            accounts_to_mute = map(lambda x: x['account_id'], mute_account_buckets)
            return self._log_collection.find(
                    {
                        'bucket_id': bucket_id,
                        'source_name': {'$nin': accounts_to_mute}
                    })
        except Exception as e:
            logger.error('Cannot find the items from bucket_id')

    def insert_log(self, data):
        try:
            # insert the entry to pf9_raw_logs
            insert_id = self._log_collection.insert_one(data).inserted_id
            return insert_id
        except Exception as e:
            logger.error('Invalid Log Entry : ' + e)

    def update_log(self, flds, bucket_id):
        try:
            # Update the pf9_raw_logs with bucket_id
            ret = self._log_collection.find_one_and_update({'_id': flds['_id']},
                                                           {'$set': {'bucket_id': bucket_id}})
            logger.info('Raw updated %s', ret)
        except Exception as e:
            logger.error('Cannot update the raw log')

    def delete_old_logs(self, numDays=3):
        try:
            d = datetime.utcnow() - timedelta(days=numDays)
            epoch = self.get_epoch(d)

            logger.info("Deleting entries older than " + str(epoch))
            epoch = long(epoch * 1000) * 1000
            logger.info('Delete old enrtry < ' + str(epoch))
            print('Delete old entry < ' +str(epoch))
            ret = self._log_collection.delete_many({"timestamp": { "$lt": epoch }})
            return ret

        except Exception as e:
            logger.error('Could not delete log entries: '+str(e))


    def get_host_logs(self, hostname):
        try:
            doc = self._log_collection.find({'hostname': hostname})

            return doc
        except Exception as e:
            logger.error('Could not get host alerts: ' + e)

    def get_logs_by_bucket(self, bucket_id):
        try:
            doc = self._log_collection.find({'bucket_id': bucket_id})

            return doc
        except Exception as e:
            logger.error('Could not get bucket alerts: ' + e)


    def get_logs_by_category(self, category):
        try:
            doc = self._log_collection.find({'subCategory': category})

            return doc
        except Exception as e:
            logger.error('Could not get category data: ' + e)

    def get_all_logs(self):
        try:
            doc = self._log_collection.find({})

            return doc
        except Exception as e:
            logger.error('Could not get category data: ' + e)

