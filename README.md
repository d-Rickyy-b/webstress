# Webstress
[![build](https://github.com/d-Rickyy-b/webstress/actions/workflows/release_build.yml/badge.svg)](https://github.com/d-Rickyy-b/webstress/actions/workflows/release_build.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/d-Rickyy-b/webstress/.svg)](https://pkg.go.dev/github.com/d-Rickyy-b/webstress/)

While developing my tool [certstream-server-go](https://github.com/d-Rickyy-b/certstream-server-go), I was searching for a tool to stress test my websocket server.
Not by sending requests, but by immitating the behaviour of a client that's on the receiving end.
I came across the python tool [wsstat](https://github.com/Fitblip/wsstat) by [Fitblip](https://github.com/Fitblip).

Sadly I ran into troubles installing the tool, so I decided to create my own.

Webstress connects to a websocket server, receives messages sent by the server and counts them.
It is able to connect via an arbitrary amount of workers, which will receive messages from the server in parallel.
This helped my stress test my websocket server and monitor how it behaved under load.

For a more advanced tool check out [websocat](https://github.com/vi/websocat).

## Screenshot
![webstress screenshot](https://github.com/d-Rickyy-b/webstress/blob/master/docs/img/webstress_impression.png?raw=true)

## Usage
Using webstress is simple. You can either download and compile the code by yourself or use one of our [precompiled binaries](https://github.com/d-Rickyy-b/webstress/releases). Since it's written in Go, you should be able to run it on any mayor OS.

```
Usage of webstress:
usage: webstress [-h|--help] -a|--remote-addr "<value>" [-r|--recover]
                 [-p|--ping-interval <integer>] [-w|--worker-count <integer>]
                 [-l|--ratelimit <integer>]

                 Websocket stress tool developed in Go

Arguments:

  -h  --help           Print help information
  -a  --remote-addr    remote address to connect to
  -r  --recover        recover from certain errors. Default: true
  -p  --ping-interval  number of seconds between pings. Default: 30
  -w  --worker-count   number of workers to start. Default: 30
  -l  --ratelimit      rate limit in messages per second per websocket. Default: 0
```

