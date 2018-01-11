# rkn-bypasser
This is [roskomnadzor's](https://eng.rkn.gov.ru/) blocker bypasser.

# How it works
It loads blocked IPs from [roskomsvoboda](http://reestr.rublacklist.net/api/ips) and run local [SOCK5 proxy](https://github.com/armon/go-socks5) with special dialer. This dialer checks requested IP if it blocked. If IP blocked it proxying request to TOR proxy server. If not it works as simple proxy. Here's how evil roskomnadzor blocking machine is bypassed.

# How to use it

## Using docker

This is easiest way. You need [docker](https://www.docker.com/community-edition) and gnu-make util. Just open project root in terminal and enter:
```
$ make
``` 
It will build dockers and run them. You only need to configure your browser (and any other soft) to use SOCK5 proxy at `127.0.1.1:8000`.

## Manually build

You can build it by your own using [go](https://golang.org/dl/):
```
$ go install github.com/someanon/rkn-baypasser
```
Download and run [TOR proxy](https://dist.torproject.org/). You might need to build from source. Alternatively you may find binary packages in your distro's repo.

Now you can use it:
```
$ BIND_ADDR=127.0.1.1:8000 TOR_PROXY=127.0.0.1:9150 rkn-baypasser
``` 

## Does it work on Windows?

Sure. All this setups possible to do in windows. You can write me if you need any help. 