#!/usr/bin/env python

"""
This is an utility script that can be used to launch python processes in daemon mode (via grizzled.os)
while also redirecting their stdout and stderr to a log file rotating utility (rotatelogs in this example).
"""

import os, sys, time
from socket import gethostname
from optparse import OptionParser
from grizzled.os import daemonize

PYTHON_BINARY = "python2.6"
PATH_TO_PYTHON_BINARY = "/usr/bin/%s" % PYTHON_BINARY
ROTATELOGS_CMD = "/usr/sbin/rotatelogs"
LOGDIR = "/opt/tornado/logs"
LOGDURATION = 86400

def main():
    usage = "usage: %prog options"
    parser = OptionParser(usage=usage)
    parser.add_option("-a", "--action", dest="action", 
         help="action to be performed for the service (e.g. stop|start|restart)")
    parser.add_option("-l", "--logger", dest="logger", 
         help="full path to log rotating program (e.g. /opt/apache2/bin/rotatelogs)")
    parser.add_option("-d", "--logdir", dest="logdir", 
         help="log directory (e.g. /opt/tornado/logs)")
    parser.add_option("-s", "--service", dest="service", 
         help="name of service module (e.g. evite.profileweb)")
    parser.add_option("-p", "--port", dest="port", 
         help="port to run service on (e.g. 9000)")
    parser.add_option("-x", "--xargs", dest="xargs", 
         help="extra arguments to be passed to the service (comma-separated list)")

    (options, args) = parser.parse_args()
    action = options.action
    logger = options.logger
    logdir = options.logdir
    service = options.service
    port = options.port
    xargs = options.xargs
    logger = options.logger

    if not action or not service:
        parser.print_help() 
        sys.exit(1)

    if action and action not in ['start', 'stop', 'restart']:
        parser.print_help() 
        sys.exit(1)

    if not logdir:
        logdir = LOGDIR

    if not logger:
        logger = ROTATELOGS_CMD

    hostname = gethostname()

    execve_args = [PYTHON_BINARY, "-m", service]
    logfile = "%s_%s_log.%%Y-%%m-%%d" % (service, hostname)
    pidfile = "%s/%s.pid" % (logdir, service)
    if port:
        logfile = "%s_%s_%s_log.%%Y-%%m-%%d" % (service, hostname, port)
        pidfile = "%s/%s_%s.pid" % (logdir, service, port)
        execve_args.append("--port=%s" % port)
    if xargs:
        xarglist = xargs.split(',')
        for xarg in xarglist:
            execve_args.append("--%s" % xarg)

    logpipe ="%s %s/%s %d" % (logger, logdir, logfile, LOGDURATION)
    execve_path = PATH_TO_PYTHON_BINARY

    if action == 'start':
        start(logpipe, execve_args, execve_path, pidfile)
    elif action == 'stop':
        stop(pidfile)
    elif action == 'restart':
        stop(pidfile)
        start(logpipe, execve_args, execve_path, pidfile)
        
def start(logpipe, execve_args, execve_path, pidfile):
    # open the pipe to ROTATELOGS
    #so = se = open("my.log", 'w', 0)
    so = se = os.popen(logpipe, 'w')

    # re-open stdout without buffering
    sys.stdout = os.fdopen(sys.stdout.fileno(), 'w', 0)

    # redirect stdout and stderr to the log file opened above
    os.dup2(so.fileno(), sys.stdout.fileno())
    os.dup2(se.fileno(), sys.stderr.fileno())

    daemonize(no_close=True, pidfile=pidfile)
    print execve_path
    os.execv(execve_path, execve_args)

def stop(pidfile):
    time.sleep(1)
    if not os.path.isfile(pidfile):
        return
        
    pid = open(pidfile).readlines()[0]
    pid = pid.rstrip()
    cmd = "kill -9 %s" % pid
#    print "running: '%s'" % cmd
    os.system(cmd)
    os.unlink(pidfile)


if __name__ == "__main__":
    main()
