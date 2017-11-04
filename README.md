# TCP concurrent connection tester

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/b211244c4a674049864d45020aa8e883)](https://www.codacy.com/app/chadell/check-max-tcp-connections?utm_source=github.com&utm_medium=referral&utm_content=dachad/check-max-tcp-connections&utm_campaign=badger) [![Build Status](https://travis-ci.org/dachad/check-max-tcp-connections.svg?branch=master)](https://travis-ci.org/dachad/check-max-tcp-connections)

## TL;DR

Tool to test concurrent connections towards a server listening to a TCP port
Script to test concurrent connections towards a server listening to a TCP port

## Description/Script steps

- Given a hostname, port, the number of connections (100 by default), 
a delay between connections (10ms by default) and an interval between stats
updates to the standard output...
- It will use goroutines to open a tcp connection and try to read from it
- The tool will exit once all connections have been dialed (successfully or not)

## Usage

```bash
% ./tcpMaxConn --help
Usage of ./tcpMaxConn:
  -y, --assume-yes        Force execution without asking for confirmation
  -c, --connections int   Number of connections you want to open (default 100)
  -d, --debug             Print debugging information to the standard error
  -h, --host string       Host you want to open tcp connections against (Required)
  -i, --interval int      Interval, in seconds, between stats updates (default 1)
  -p, --port int          Port you want to open tcp connections against (Required)
  -s, --sleep int         Time you want to sleep between connections, in ms (default 10)
```

## Example

```bash
% ./tcpMaxConn --host myhttpsamplehost.com --port 80 --connections 10 --sleep 1000 -y
Total: 10, Dialing: 0, Established: 0, Closed: 0, Error: 0, NotInitiated: 10
Total: 10, Dialing: 0, Established: 1, Closed: 0, Error: 0, NotInitiated: 9
Total: 10, Dialing: 1, Established: 2, Closed: 0, Error: 0, NotInitiated: 7
Total: 10, Dialing: 1, Established: 3, Closed: 0, Error: 0, NotInitiated: 6
Total: 10, Dialing: 1, Established: 4, Closed: 0, Error: 0, NotInitiated: 5
Total: 10, Dialing: 1, Established: 5, Closed: 0, Error: 0, NotInitiated: 4
Total: 10, Dialing: 1, Established: 6, Closed: 0, Error: 0, NotInitiated: 3
Total: 10, Dialing: 1, Established: 7, Closed: 0, Error: 0, NotInitiated: 2
Total: 10, Dialing: 1, Established: 8, Closed: 0, Error: 0, NotInitiated: 1
Total: 10, Dialing: 1, Established: 9, Closed: 0, Error: 0, NotInitiated: 0
Total: 10, Dialing: 0, Established: 10, Closed: 0, Error: 0, NotInitiated: 0
Total: 10, Dialing: 0, Established: 10, Closed: 0, Error: 0, NotInitiated: 0
```

## TO-DO

- Timeout configuration (max-duration?) so this can be reused for CI tests (if you cannot open X concurrent requests in 1 second, thats potentially a problem) 
- Keepalive connections / reopen closed connections to keep this number of concurrent connections during an specific time (max-duration?)
- more test coverage
- "auto-incremental" mode; it opens connections at an specific rate until it fails or it times-out, giving you an idea of the max concurrency your service supports
- distributed executions; several daemons may be able to collaborate to measure the capacity of an specific target
- Docker image and OS packages for the most common distros
- See also the issues in this Github repo!!
