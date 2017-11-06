# TCP concurrent connection tester

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/b211244c4a674049864d45020aa8e883)](https://www.codacy.com/app/dachad/tcpgoon?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=dachad/tcpgoon&amp;utm_campaign=Badge_Grade)
[![Build Status](https://travis-ci.org/dachad/tcpgoon.svg?branch=master)](https://travis-ci.org/dachad/tcpgoon)
[![Go Report Card](https://goreportcard.com/badge/github.com/dachad/tcpgoon)](https://goreportcard.com/report/github.com/dachad/tcpgoon)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/dachad/tcpgoon/blob/master/LICENSE)

## TL;DR

Tool to test concurrent connections towards a server listening to a TCP port

## Description

- Given a hostname, port, the number of connections (100 by default), 
a delay between connections (10ms by default) and an interval between stats
updates to the standard output...
- It will use goroutines to open a tcp connection and try to read from it
- The tool will exit once all connections have been dialed (successfully or not)
- Exit status different from 0 represent executions where all connections were not 
established successfully

## Usage

```bash
% ./tcpgoon --help
Usage of ./tcpgoon:
  -y, --assume-yes         Force execution without asking for confirmation
  -c, --connections int    Number of connections you want to open (default 100)
  -d, --debug              Print debugging information to the standard error
  -t, --dial-timeout int   Connection dialing timeout, in s (default 5)
  -h, --host string        Host you want to open tcp connections against (Required)
  -i, --interval int       Interval, in seconds, between stats updates (default 1)
  -p, --port int           Port you want to open tcp connections against (Required)
  -s, --sleep int          Time you want to sleep between connections, in ms (default 10)
```

## Example

Successful execution (connections were opened as expected):
```bash
% ./tcpgoon --host myhttpsamplehost.com --port 80 --connections 10 --sleep 999 -y 
Total: 10, Dialing: 0, Established: 0, Closed: 0, Error: 0, NotInitiated: 10
Total: 10, Dialing: 1, Established: 1, Closed: 0, Error: 0, NotInitiated: 8
Total: 10, Dialing: 1, Established: 2, Closed: 0, Error: 0, NotInitiated: 7
Total: 10, Dialing: 1, Established: 3, Closed: 0, Error: 0, NotInitiated: 6
Total: 10, Dialing: 1, Established: 4, Closed: 0, Error: 0, NotInitiated: 5
Total: 10, Dialing: 1, Established: 5, Closed: 0, Error: 0, NotInitiated: 4
Total: 10, Dialing: 1, Established: 6, Closed: 0, Error: 0, NotInitiated: 3
Total: 10, Dialing: 1, Established: 7, Closed: 0, Error: 0, NotInitiated: 2
Total: 10, Dialing: 1, Established: 8, Closed: 0, Error: 0, NotInitiated: 1
Total: 10, Dialing: 1, Established: 9, Closed: 0, Error: 0, NotInitiated: 0
Total: 10, Dialing: 0, Established: 10, Closed: 0, Error: 0, NotInitiated: 0
--- myhttpsamplehost.com:80 tcp test statistics ---
Total: 10, Dialing: 0, Established: 10, Closed: 0, Error: 0, NotInitiated: 0
% echo $?
0
```
Unsuccessful execution (unable to open connections against the destination host:port):
```bash
% ./tcpgoon --host myhttpsamplehost.com --port 81 --connections 10 --sleep 999 -y
Total: 10, Dialing: 0, Established: 0, Closed: 0, Error: 0, NotInitiated: 10
Total: 10, Dialing: 2, Established: 0, Closed: 0, Error: 0, NotInitiated: 8
Total: 10, Dialing: 3, Established: 0, Closed: 0, Error: 0, NotInitiated: 7
Total: 10, Dialing: 3, Established: 0, Closed: 0, Error: 1, NotInitiated: 6
Total: 10, Dialing: 3, Established: 0, Closed: 0, Error: 2, NotInitiated: 5
Total: 10, Dialing: 3, Established: 0, Closed: 0, Error: 3, NotInitiated: 4
Total: 10, Dialing: 3, Established: 0, Closed: 0, Error: 4, NotInitiated: 3
Total: 10, Dialing: 3, Established: 0, Closed: 0, Error: 5, NotInitiated: 2
Total: 10, Dialing: 3, Established: 0, Closed: 0, Error: 6, NotInitiated: 1
Total: 10, Dialing: 3, Established: 0, Closed: 0, Error: 7, NotInitiated: 0
Total: 10, Dialing: 2, Established: 0, Closed: 0, Error: 8, NotInitiated: 0
Total: 10, Dialing: 1, Established: 0, Closed: 0, Error: 9, NotInitiated: 0
--- myhttpsamplehost.com:81 tcp test statistics ---
Total: 10, Dialing: 0, Established: 0, Closed: 0, Error: 10, NotInitiated: 0
% echo $?
2
```

## TO-DO

- Timeout configuration (max-duration?) so this can be reused for CI tests (if you cannot open X concurrent requests in 1 second, thats potentially a problem) 
- Keepalive connections / reopen closed connections to keep this number of concurrent connections during an specific time (max-duration?)
- more test coverage
- "auto-incremental" mode; it opens connections at an specific rate until it fails or it times-out, giving you an idea of the max concurrency your service supports
- distributed executions; several daemons may be able to collaborate to measure the capacity of an specific target
- Docker image and OS packages for the most common distros
- See also the issues in this Github repo!!
