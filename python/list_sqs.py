#!/usr/bin/env python

import os, sys
import re
import time
import optparse
import simplejson as json
from boto.sqs.connection import SQSConnection
from boto.sqs.message import Message

AWS_ACCESS_KEY = "YOUR_AWS_ACCESS_KEY"
AWS_SECRET_KEY = "YOUR_AWS_SECRET_KEY"
SQS_QUEUE_NAME = 'EMAIL_STATUS_PROD'

sqs_conn = SQSConnection(AWS_ACCESS_KEY, AWS_SECRET_KEY)
sqs = sqs_conn.create_queue(SQS_QUEUE_NAME)
#sqs_conn.delete_queue(sqs)

print "Queues in production:"

rs = sqs_conn.get_all_queues()
for q in rs:
    print q.id

print "Contents of %s queue:" % SQS_QUEUE_NAME

#msgid = "XZJSYCJUAVXAUSTPWMYU"
#status = "Sent (n8L4kwMM014545 Message accepted for delivery)"
#msg_info = {
#    "trackedMessageId": msgid,
#    "sendmailLogStatus": status
#}
#sqs_msg_body = json.dumps(msg_info)
#m = Message()
#m.set_body(sqs_msg_body)
#sqs.write(m)

while True:
    m = sqs.read(5)
    if not m:
        break
    print m.get_body()
    #sqs.delete_message(m)
