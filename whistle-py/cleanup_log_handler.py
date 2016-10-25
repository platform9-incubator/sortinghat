import logging
import logging.config

LOG_FILENAME = './logs/cleanup.log'
FILE_SIZE = 104857600
FILE_COUNT = 5

logging.config.fileConfig("./cleanup.config", defaults={'logfilename': LOG_FILENAME,
                                                        'filesize': FILE_SIZE,
                                                        'filecount':FILE_COUNT})
