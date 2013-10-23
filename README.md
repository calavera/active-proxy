# Active-Proxy

Active-Proxy is a dynamic reverse proxy written in Go.

It uses [etcd](https://github.com/coreos/etcd) to dynamically switch application upstreams without restarting.

This is just an experiment, so any feature that you could think of is probably missed. Pull requests are really appreciated.

## Dependencies

This repository includes a vagrant file with a box that already includes the latest stable version of etcd.

Use `vagrant up` to boot the box. The server listens then at localhost:8080.

It uses [gopack](https://github.com/d2fn/gopack) to install other required dependencies when the project builds.

## Installation

You can build it running:

```
$ ./gopack build
```

Or run it from source using:

```
$ script/server
```

## Usage

Imagine that you have an applications called `foo` and `bar` that you want to server from the same port.
`foo` runs on `http://localhost:4000/foo` and `bar` runs on `http://localhost:5000/bar`.

You can make Active-Proxy to run both on `:8080/foo` and `:8080/bar` registering them in the etcd cluster:

```
$ curl -L http://127.0.0.1:8080/v1/keys/applications/foo/port -d value=4000
$ curl -L http://127.0.0.1:8080/v1/keys/applications/bar/port -d value=5000
```

Imagine then that you want to redeploy `foo`. You can boot the new code in a different port, tell Active-Proxy to serve new requests through that port and kill the old application.
There are no complicated hotdeploy strategies required and it allows you to test that the application boots correctly before switching ports.

You only need to register the new port in the etcd cluster and Active-Proxy will do the switch for you without restarting the proxy either:

```
$ curl -L http://127.0.0.1:8080/v1/keys/applications/foo/port -d value=4004
```
