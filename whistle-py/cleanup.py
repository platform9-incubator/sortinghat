# Copyright (c) 2015 Platform9 Systems Inc.
from fuzzy_match import logs_db

from cleanup_log_handler import logging

logger = logging.getLogger(__name__)


logger.info('Deleted Old Log Entries')
logs_db.delete_old_logs()
#Rebuild the categories
logger.info('Rebuilding buckets')
#bucket  = recompute_buckets()
logger.info('Rebuilding buckets done ')
