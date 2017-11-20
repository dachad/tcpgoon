<p align="center">
<img src="http://www.confusedcoders.com/wp-content/uploads/2016/10/golang-1.jpg" alt="tcpgoon" title="tcpgoon" />
</p>

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/b211244c4a674049864d45020aa8e883)](https://www.codacy.com/app/dachad/tcpgoon?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=dachad/tcpgoon&amp;utm_campaign=Badge_Grade)
[![Build Status](https://travis-ci.org/dachad/tcpgoon.svg?branch=master)](https://travis-ci.org/dachad/tcpgoon)
[![Go Report Card](https://goreportcard.com/badge/github.com/dachad/tcpgoon)](https://goreportcard.com/report/github.com/dachad/tcpgoon)
[![](https://images.microbadger.com/badges/image/dachad/tcpgoon.svg)](https://microbadger.com/images/dachad/tcpgoon "Get your own image badge on microbadger.com")
[![](https://images.microbadger.com/badges/version/dachad/tcpgoon.svg)](https://microbadger.com/images/dachad/tcpgoon "Get your own version badge on microbadger.com")
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

## Examples

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
Timing stats for 10 established connections min/avg/max/dev = 17.929ms/19.814ms/29.811ms/3.353ms
% echo $?
0
```

Partially succeeded execution (mix of successes and errors against the target):
```bash
% ./tcpgoon --host myhttpsamplehost.com --port 8080 --connections 10 --sleep 999 -y
Total: 10, Dialing: 0, Established: 0, Closed: 0, Error: 0, NotInitiated: 10
Total: 10, Dialing: 0, Established: 1, Closed: 0, Error: 0, NotInitiated: 9
Total: 10, Dialing: 0, Established: 2, Closed: 0, Error: 0, NotInitiated: 8
Total: 10, Dialing: 1, Established: 2, Closed: 0, Error: 0, NotInitiated: 7
Total: 10, Dialing: 2, Established: 2, Closed: 0, Error: 0, NotInitiated: 6
Total: 10, Dialing: 3, Established: 2, Closed: 0, Error: 0, NotInitiated: 5
Total: 10, Dialing: 4, Established: 2, Closed: 0, Error: 0, NotInitiated: 4
Total: 10, Dialing: 4, Established: 2, Closed: 0, Error: 1, NotInitiated: 3
Total: 10, Dialing: 5, Established: 2, Closed: 0, Error: 2, NotInitiated: 1
Total: 10, Dialing: 5, Established: 2, Closed: 0, Error: 3, NotInitiated: 0
Total: 10, Dialing: 4, Established: 2, Closed: 0, Error: 4, NotInitiated: 0
Total: 10, Dialing: 3, Established: 2, Closed: 0, Error: 5, NotInitiated: 0
Total: 10, Dialing: 3, Established: 2, Closed: 0, Error: 5, NotInitiated: 0
Total: 10, Dialing: 2, Established: 2, Closed: 0, Error: 6, NotInitiated: 0
Total: 10, Dialing: 1, Established: 2, Closed: 0, Error: 7, NotInitiated: 0
Total: 10, Dialing: 0, Established: 2, Closed: 0, Error: 8, NotInitiated: 0
--- myhttpsamplehost.com:8080 tcp test statistics ---
Total: 10, Dialing: 0, Established: 2, Closed: 0, Error: 8, NotInitiated: 0
Timing stats for 2 established connections min/avg/max/dev = 1.914ms/2.013ms/2.113ms/99µs
Timing stats for 8 failed connections min/avg/max/dev = 5.000819s/5.002496s/5.004758s/1.448ms
% echo $?
2
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
Timing stats for 10 failed connections min/avg/max/dev = 5.00025s/5.001741s/5.00317s/908µs
% echo $?
2
```

## Executing the tests

You can use the standard go test command, or use our scripts we also run as CI.

Main tests execution:
```bash
% ./_script/test
```

Emulation of a travis job execution using docker (of course, it needs docker):
```bash
% ./_script/cibuild-docker
```

And also emulating a travis job deployment (it publishes new binaries
providing successful tests and the right credentials):
```bash
% ./_script/cibuild-docker -d
```

## TO-DO

We do use [Github issues](/issues/) to track bugs, improvements and feature requests. Do not hesitate
to raise new ones ;)
