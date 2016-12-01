import logging
import logging.config

LOG_FILENAME = 'logs/debug.log'
WEB_LOGS = 'logs/access.log'
#100 MB
FILE_SIZE = 104857600

#5 files totalling 500MB max
FILE_COUNT = 5
logging.config.fileConfig("logging.config", defaults={  'logfilename': LOG_FILENAME,
                                                        'accesslogfilename': WEB_LOGS,
                                                        'filesize': FILE_SIZE,
                                                        'filecount':FILE_COUNT})
