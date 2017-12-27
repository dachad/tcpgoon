<p align="center">
<img src="https://raw.githubusercontent.com/dachad/tcpgoon/master/_imgs/tcpgoontransparent.png" alt="tcpgoon" title="tcpgoon" width="400"/>
</p>
<p align="center">
<img src="https://raw.githubusercontent.com/dachad/tcpgoon/master/_imgs/coollogo_com-290231302.png" alt="tcpgoon" title="tcpgoon" width="438"/>
</p>

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/b211244c4a674049864d45020aa8e883)](https://www.codacy.com/app/dachad/tcpgoon?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=dachad/tcpgoon&amp;utm_campaign=Badge_Grade)
[![Build Status](https://travis-ci.org/dachad/tcpgoon.svg?branch=master)](https://travis-ci.org/dachad/tcpgoon)
[![Go Report Card](https://goreportcard.com/badge/github.com/dachad/tcpgoon)](https://goreportcard.com/report/github.com/dachad/tcpgoon)
[![](https://images.microbadger.com/badges/image/dachad/tcpgoon.svg)](https://microbadger.com/images/dachad/tcpgoon "Get your own image badge on microbadger.com")
[![](https://images.microbadger.com/badges/version/dachad/tcpgoon.svg)](https://microbadger.com/images/dachad/tcpgoon "Get your own version badge on microbadger.com")
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/dachad/tcpgoon/blob/master/LICENSE)

**[TLDR](#tldr)** . **[Description](#description)** . **[Usage](#usage)** . **[Help](#help)** . **[Examples](#examples)** . **[Execution using Docker](#execution-using-docker)** . **[Extra project information](#extra-project-information)** . **[Why do I want to test tcp connections?](#why-do-i-want-to-test-tcp-connections)** . **[Where does the project name come from?](#where-does-the-project-name-come-from)** . **[Authors](#authors)** . **[Especial thanks to...](#especial-thanks-to)** . **[Development information](#development-information)** . **[TO-DO](#to-do)** . **[README maintenance](#readme-maintenance)** . **[Testing locally](#testing-locally)** . 

## TL;DR

Tool to test concurrent connections towards a server listening to a TCP port

## Description

* Given a hostname, port, the number of connections (100 by default), 
a delay between connections (10ms by default) and an interval between stats
updates to the standard output...
* It will use goroutines to open a tcp connection and try to read from it
* The tool will exit once all connections have been dialed (successfully or not)
* Exit status different from 0 represent executions where all connections were not 
established successfully, facilitating the integration with test suites.

## Usage

### Help

```bash
% ./tcpgoon --help
tcpgoon tests concurrent connections towards a server listening on a TCP port

Usage:
  tcpgoon [flags] <host> <port>

Flags:
  -y, --assume-yes         Force execution without asking for confirmation
  -c, --connections int    Number of connections you want to open (default 100)
  -t, --dial-timeout int   Connection dialing timeout, in ms (default 5000)
  -h, --help               help for tcpgoon
  -i, --interval int       Interval, in seconds, between stats updates (default 1)
  -s, --sleep int          Time you want to sleep between connections, in ms (default 10)
  -d, --debug              Print debugging information to the standard error
```

### Examples

Successful execution (connections were opened as expected):
```bash
% ./tcpgoon myhttpsamplehost.com 80 --connections 10 --sleep 999 -y 
Total: 10, Dialing: 0, Established: 0, Closed: 0, Error: 0, NotInitiated: 10
Total: 10, Dialing: 1, Established: 1, Closed: 0, Error: 0, NotInitiated: 8
Total: 10, Dialing: 1, Established: 2, Closed: 0, Error: 0, NotInitiated: 7
...
Total: 10, Dialing: 1, Established: 8, Closed: 0, Error: 0, NotInitiated: 1
Total: 10, Dialing: 1, Established: 9, Closed: 0, Error: 0, NotInitiated: 0
Total: 10, Dialing: 0, Established: 10, Closed: 0, Error: 0, NotInitiated: 0
--- myhttpsamplehost.com:80 tcp test statistics ---
Total: 10, Dialing: 0, Established: 10, Closed: 0, Error: 0, NotInitiated: 0
Response time stats for 10 established connections min/avg/max/dev = 17.929ms/19.814ms/29.811ms/3.353ms
% echo $?
0
```

Partially succeeded execution (mix of successes and errors against the target):
```bash
% ./tcpgoon myhttpsamplehost.com 8080 --connections 10 --sleep 999 -y
Total: 10, Dialing: 0, Established: 0, Closed: 0, Error: 0, NotInitiated: 10
Total: 10, Dialing: 0, Established: 1, Closed: 0, Error: 0, NotInitiated: 9
Total: 10, Dialing: 0, Established: 2, Closed: 0, Error: 0, NotInitiated: 8
Total: 10, Dialing: 1, Established: 2, Closed: 0, Error: 0, NotInitiated: 7
...
Total: 10, Dialing: 2, Established: 2, Closed: 0, Error: 6, NotInitiated: 0
Total: 10, Dialing: 1, Established: 2, Closed: 0, Error: 7, NotInitiated: 0
Total: 10, Dialing: 0, Established: 2, Closed: 0, Error: 8, NotInitiated: 0
--- myhttpsamplehost.com:8080 tcp test statistics ---
Total: 10, Dialing: 0, Established: 2, Closed: 0, Error: 8, NotInitiated: 0
Response time stats for 2 established connections min/avg/max/dev = 1.914ms/2.013ms/2.113ms/99µs
Time to error stats for 8 failed connections min/avg/max/dev = 5.000819s/5.002496s/5.004758s/1.448ms
% echo $?
2
```

Unsuccessful execution (unable to open connections against the destination host:port):
```bash
% ./tcpgoon myhttpsamplehost.com 81 --connections 10 --sleep 999 -y
Total: 10, Dialing: 0, Established: 0, Closed: 0, Error: 0, NotInitiated: 10
Total: 10, Dialing: 2, Established: 0, Closed: 0, Error: 0, NotInitiated: 8
Total: 10, Dialing: 3, Established: 0, Closed: 0, Error: 0, NotInitiated: 7
...
Total: 10, Dialing: 3, Established: 0, Closed: 0, Error: 6, NotInitiated: 1
Total: 10, Dialing: 3, Established: 0, Closed: 0, Error: 7, NotInitiated: 0
Total: 10, Dialing: 2, Established: 0, Closed: 0, Error: 8, NotInitiated: 0
Total: 10, Dialing: 1, Established: 0, Closed: 0, Error: 9, NotInitiated: 0
--- myhttpsamplehost.com:81 tcp test statistics ---
Total: 10, Dialing: 0, Established: 0, Closed: 0, Error: 10, NotInitiated: 0
Time to error stats for 10 failed connections min/avg/max/dev = 5.00025s/5.001741s/5.00317s/908µs
% echo $?
2
```

### Execution using Docker

See our [public docker image and its documentation](https://hub.docker.com/r/dachad/tcpgoon/). The image
is being updated continuously; you can bind to specific versions, or to the "latest" tag.


## Extra project information

### Why do I want to test tcp connections?
 
Stressing TCP connections against a server/application facilitates the detection of
bottlenecks/issues limiting the capacity of this server/application to accept/keep an specific
(and potentially) large number of parallel connections. Some examples of typical (configuration) 
issues:

* OS configuration (tcp backlog, network drivers buffers),
* number of file descriptors/processes the server can use,
* application listener properties...

These limitations may pass unnoticed in actual application-l7 stress tests, given 
other bottlenecks can arise earlier than these limitations, or degradation scenarios and/or special
conditions may not be reproduced during the stress tests execution, but real life (lots of connections
queued because of a dependency taking longer to reply than it usually does?)
 
hping is not an actual option for this use case given it won't complete a 3-way handshake,
so the connection will not reach the accept() syscall of your server/application or fill up your
tcp backlog.

### Where does the project name come from?
```
Goon: /ɡuːn/ noun informal; noun: goon; plural noun: goons ;
...
2.
NORTH AMERICAN
a bully or thug, especially a member of an armed or security force.
...
```
<p align="center">
<img src="https://raw.githubusercontent.com/dachad/tcpgoon/master/_imgs/thegoon.jpg" alt="thegoon" title="thegoon" width="250"/>
</p>


### Authors

* [Christian Adell](https://github.com/chadell)
* [Daniel Caballero](https://github.com/dcaba)

### Especial thanks to...

* [Linafm design](https://www.facebook.com/linafmdisegni/), for our custom and nice Goon Gopher


## Development information

### TO-DO

We do use [Github issues](https://github.com/dachad/tcpgoon/issues) to track bugs, improvements and feature requests. Do not hesitate
to raise new ones, or solve them for us by raising PRs ;)

### README maintenance

Do not edit README.md, as its autogenerated and your changes will be lost. Consider README.src.md. Dockerhub README requires manual
maintenance, bringing relevant aspects from here and adapting cmdusage (by docker run...)

### Testing locally

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


