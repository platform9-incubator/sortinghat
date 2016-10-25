# Copyright (c) 2015 Platform9 Systems Inc.
import threading

from Queue import Queue
import requests
__author__ = 'roopak'

OUTGOING_URL='https://hooks.slack.com/services/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'

out_queue = Queue()

def push_message(message):
    out_queue.put(message)



def _push_webhook():
    while True:
        message = out_queue.get()
        payload = {'payload': message}
        requests.post(OUTGOING_URL, params=payload)

def start_webhook_processor():
    t = threading.Thread(target= _push_webhook)
    t.daemon = True
    t.start()