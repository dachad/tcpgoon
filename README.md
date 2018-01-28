<p align="center">
<img src="https://raw.githubusercontent.com/dachad/tcpgoon/master/_imgs/tcpgoontransparent.png" alt="tcpgoon" title="tcpgoon" width="380"/>
</p>
<p align="center">
<img src="https://raw.githubusercontent.com/dachad/tcpgoon/master/_imgs/coollogo_com-290231302.png" alt="tcpgoon" title="tcpgoon" width="420"/>
</p>

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/b211244c4a674049864d45020aa8e883)](https://www.codacy.com/app/dachad/tcpgoon?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=dachad/tcpgoon&amp;utm_campaign=Badge_Grade)
[![Build Status](https://travis-ci.org/dachad/tcpgoon.svg?branch=master)](https://travis-ci.org/dachad/tcpgoon)
[![Go Report Card](https://goreportcard.com/badge/github.com/dachad/tcpgoon)](https://goreportcard.com/report/github.com/dachad/tcpgoon)
[![Coverage Status](https://coveralls.io/repos/github/dachad/tcpgoon/badge.svg?branch=Issue45_CoverallsReport)](https://coveralls.io/github/dachad/tcpgoon?branch=Issue45_CoverallsReport)
[![](https://images.microbadger.com/badges/image/dachad/tcpgoon.svg)](https://microbadger.com/images/dachad/tcpgoon "Get your own image badge on microbadger.com")
[![](https://images.microbadger.com/badges/version/dachad/tcpgoon.svg)](https://microbadger.com/images/dachad/tcpgoon "Get your own version badge on microbadger.com")
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/dachad/tcpgoon/blob/master/LICENSE)

**[TLDR](#tldr)** . **[Description](#description)** . **[Usage](#usage)** . **[Execution using Docker](#execution-using-docker)** . **[Installation](#installation)** . **[Help](#help)** . **[Examples](#examples)** . **[Extra project information](#extra-project-information)** . **[Why do I want to test TCP connections?](#why-do-i-want-to-test-tcp-connections)** . **[Where does the project name come from?](#where-does-the-project-name-come-from)** . **[Authors](#authors)** . **[Especial thanks to...](#especial-thanks-to)** . **[Development information](#development-information)** . **[TO-DO](#to-do)** . **[Project structure](#project-structure)** . **[README maintenance](#readme-maintenance)** . **[Testing locally](#testing-locally)** . 

## TL;DR

Tool to test concurrent connections towards a server listening to a TCP port

## Description

* Given a hostname, port, the number of connections (100 by default), 
a delay between connections (10ms by default) and an interval between stats
updates to the standard output...
* It will use goroutines to open TCP connections and try to read from them
* The tool will exit once all connections have been dialed (successfully or not)
* Exit status different from 0 represent executions where all connections were not 
established successfully, facilitating the integration in test suites.

## Usage

### Execution using Docker

See our [public docker image and its documentation](https://hub.docker.com/r/dachad/tcpgoon/). The image
is being updated continuously; you can bind to specific versions, or to the "latest" tag.

### Installation

```bash
% go get github.com/dachad/tcpgoon
```

### Help

```bash
% tcpgoon --help
tcpgoon tests concurrent connections towards a server listening on a TCP port

Usage:
  tcpgoon [command]

Available Commands:
  help        Help about any command
  run         Run tcpgoon test
  version     Show tcpgoon version

Flags:
  -h, --help   help for tcpgoon

Use "tcpgoon [command] --help" for more information about a command.
% tcpgoon run --help
Run tcpgoon test

Usage:
  tcpgoon run [flags] <host> <port>

Flags:
  -y, --assume-yes         Force execution without asking for confirmation
  -c, --connections int    Number of connections you want to open (default 100)
  -d, --debug              Print debugging information to the standard error
  -t, --dial-timeout int   Connection dialing timeout, in ms (default 5000)
  -h, --help               help for run
  -i, --interval int       Interval, in seconds, between stats updates (default 1)
  -s, --sleep int          Time you want to sleep between connections, in ms (default 10)
```

### Examples

> *Note: depending on the number of connections you want to open, 
you may need to increase the number of file descriptors your user supports*

Successful execution (connections were opened as expected):
```bash
% tcpgoon run myhttpsamplehost.com 80 --connections 4 --sleep 999 -y
Total: 4, Dialing: 0, Established: 0, Closed: 0, Error: 0, NotInitiated: 4
Total: 4, Dialing: 1, Established: 1, Closed: 0, Error: 0, NotInitiated: 2
Total: 4, Dialing: 1, Established: 2, Closed: 0, Error: 0, NotInitiated: 1
Total: 4, Dialing: 1, Established: 3, Closed: 0, Error: 0, NotInitiated: 0
Total: 4, Dialing: 0, Established: 4, Closed: 0, Error: 0, NotInitiated: 0
--- myhttpsamplehost.com(216.58.201.131):80 tcp test statistics ---
Total: 4, Dialing: 0, Established: 4, Closed: 0, Error: 0, NotInitiated: 0
Response time stats for 4 established connections min/avg/max/dev = 59.842ms/70.798ms/100.766ms/17.312ms

% echo $?
0
```

Unsuccessful execution (unable to open connections against the destination host:port):
```bash
% tcpgoon run myhttpsamplehost.com 81 --connections 4 --sleep 999 -t 1 -y
Total: 4, Dialing: 0, Established: 0, Closed: 0, Error: 0, NotInitiated: 4
Total: 4, Dialing: 1, Established: 0, Closed: 0, Error: 1, NotInitiated: 2
Total: 4, Dialing: 0, Established: 0, Closed: 0, Error: 3, NotInitiated: 1
Total: 4, Dialing: 0, Established: 0, Closed: 0, Error: 4, NotInitiated: 0
Total: 4, Dialing: 0, Established: 0, Closed: 0, Error: 4, NotInitiated: 0
--- myhttpsamplehost.com(216.58.201.131):81 tcp test statistics ---
Total: 4, Dialing: 0, Established: 0, Closed: 0, Error: 4, NotInitiated: 0
Time to error stats for 4 failed connections min/avg/max/dev = 1.085ms/1.107ms/1.148ms/24µs

% echo $?
2
```

## Extra project information

### Why do I want to test TCP connections?
 
Stressing TCP connections against a server/application facilitates the detection of
bottlenecks/issues limiting the capacity of this server/application to accept/keep a specific
(and potentially) large number of parallel connections. Some examples of typical (configuration) 
issues:

* OS configuration (TCP backlog, network drivers buffers),
* number of file descriptors/processes the server can use,
* application listener properties...

These limitations may pass unnoticed in actual application-l7 stress tests, given 
other bottlenecks can arise earlier than these limitations, or degradation scenarios and/or special
conditions may not be reproduced during the stress tests execution, but in real life (lots of connections
queued because of a dependency taking longer to reply than it usually does?)
 
hping is not an actual option for this use case given it won't complete a 3-way handshake,
so the connection will not reach the accept() syscall of your server/application or fill up your
TCP backlog.

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

### Project structure

This project uses a layered topology, where *cmd* (complemented by *cmdutil*) takes care of commands/flags/arguments and uses
*mtcpclient*, which owns and knows everything about "multiple TCP connections" (including reporting), while *tcpclient*
only cares about managing single TCP connections. *tcpserver* is just there as a dependency for the other packages' tests.

<p align="center">
<img src="https://raw.githubusercontent.com/dachad/tcpgoon/master/_imgs/godepgraph.png" width="450"/>
</p>

A shared package (*debugging*) is also supplied just as a basic mechanism to control debug output.

### README maintenance

Do not edit README.md directly, as your changes will be lost. Consider README.src.md and 
the execution (requires [godepgraph](https://github.com/kisielk/godepgraph) and [Graphviz](http://graphviz.org/)) of:
```bash
% ./_script/readme_generator
```

Samples injected in the readme can be found in the `_script/readme_generator_samples/` directory.

Dockerhub README requires manual maintenance, bringing relevant aspects from here and adapting cmdusage (by docker run...)

### Testing locally

You can use the standard go test command, or use our scripts we also run as CI.

Main tests execution:
```bash
% ./_script/test
```

Emulation of a travis job execution using docker (of course, it needs docker):
```bash
% COVERALLS_TOKEN="myToken" ./_script/cibuild-docker
```

And also emulating a travis job deployment (it publishes new binaries
providing successful tests and the right credentials, unless you also use the -r flag / dry-run):
```bash
% COVERALLS_TOKEN="myToken" ./_script/cibuild-docker -d
```


