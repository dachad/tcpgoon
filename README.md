# TCP concurrent connection tester

## TL;DR

Script to test the concurrent connections towards a server listening to a TCP port

## Description/Script steps

- Given a hostname, port, and optionally the number of connections (100 by default) and delay between connections (10ms by default)...
- It will use goroutines to open a tcp connection and try to read from it, waiting the specified delay between each light-thread creation
- The main flow will just wait until all goroutines have finished (that is, when the OS detects the tcp connection as closed)


## Usage

```bash
 % ./tcpMaxConn -h
Usage of ./tcpMaxConn:
  -connections int
        Number of connections you want to open (default 100)
  -delay int
        Number of ms you want to sleep between each connection creation (default 10)
  -host string
        Host you want to open tcp connections against (default "localhost")
  -port int
        Port you want to open tcp connections against (default 8888)
```

## Example

```bash
% ./tcpMaxConn -host ec2-54-229-56-140.eu-west-1.compute.amazonaws.com -port 8080 -connections 5 
Initiating runner # 1
         runner 1 is initiating a connection
Runner 1 initated. Remaining: 4
Initiating runner # 2
         runner 2 is initiating a connection
Runner 2 initated. Remaining: 3
Initiating runner # 3
         runner 3 is initiating a connection
Runner 3 initated. Remaining: 2
Initiating runner # 4
         runner 4 is initiating a connection
Runner 4 initated. Remaining: 1
Initiating runner # 5
         runner 5 is initiating a connection
Runner 5 initated. Remaining: 0
Waiting runners to finish
         runner 2 established the connection
         runner 1 established the connection
         runner 4 established the connection
         runner 3 established the connection
         runner 5 established the connection
         runner 2 got its connection closed
         runner 1 got its connection closed
         runner 4 got its connection closed
         runner 5 got its connection closed
         runner 3 got its connection closed

Terminating Program
```

## TO-DO

- Provide "live stats" of the number of threads and the status of each connection (using channels?)
- Timeout configuration (max-duration?) so this can be reused for CI tests (if you cannot open X concurrent requests in 1 second, thats potentially a problem) 
- Keepalive connections / reopen closed connections to keep this number of concurrent connections during an specific time (max-duration?)
- Tests of this test ;)
- "auto-incremental" mode; it opens connections at an specific rate until it fails or it times-out, giving you an idea of the max concurrency your service supports
- distributed executions; several daemons may be able to collaborate to measure the capacity of an specific target
