#!/usr/bin/python

import re
import sys
import fileinput
import requests

def check_webhooks(host):
	url = "http://%s:8082/v1/webhooks/?debugOnly=1" % host
	r = requests.get(url)
	response = r.text
	lines = response.split('\n')
	last_line = lines[-1]
	pattern = ".*Running&quot;:(\d+),&quot;Pending&quot;:(\d+).*"
	prog = re.compile(pattern, flags=re.MULTILINE)
	result = prog.match(last_line)
	if result is None:
		return 1

	groups = result.groups()
	running = groups[0]
	pending = groups[1]	
	print "running: %s" % running
	print "pending: %s" % pending
	max_pending = 100
	if int(pending) > max_pending:
		print "Failing because count of pending jobs > %d" % max_pending
		return 1
	return 0

def main():
	rc = 0
	for line in fileinput.input():
		line = line.rstrip('\n')
		if not line:
			continue
		host = line
		rc += check_webhooks(host)
	sys.exit(rc)

if __name__ == "__main__":
	main()

