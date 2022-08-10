# Webstress
[![build](https://github.com/d-Rickyy-b/webstress/actions/workflows/release_build.yml/badge.svg)](https://github.com/d-Rickyy-b/webstress/actions/workflows/release_build.yml)

While developing my tool [certstream-server-go](https://github.com/d-Rickyy-b/certstream-server-go), I was searching for a tool to stress test my websocket server.
Not by sending requests, but by immitating the behaviour of a client that's on the receiving end.
I came across the python tool [wsstat](https://github.com/Fitblip/wsstat) by [Fitblip](https://github.com/Fitblip).

Sadly I ran into troubles installing the tool, so I decided to create my own.

## Screenshot
![webstress screenshot](https://github.com/d-Rickyy-b/webstress/blob/master/docs/img/webstress_impression.png?raw=true)

## Usage
Using webstress is simple. You can either download and compile the code by yourself or use one of our [precompiled binaries](https://github.com/d-Rickyy-b/webstress/releases). Since it's written in Go, you should be able to run it on any mayor OS.

```
Usage of webstress:
  -pingInterval int
        number of seconds between pings (default 30)
  -remoteAddr string
        remote address to connect to (default "ws://localhost:8080/")
  -workerCount int
        number of workers to start (default 30)
```

