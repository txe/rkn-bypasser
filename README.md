# rkn-bypasser
This is [roskomnadzor's](https://eng.rkn.gov.ru/) blocker bypasser.

# How it works
It loads blocked IPs from [roskomsvoboda](http://reestr.rublacklist.net/api/ips) and run local [SOCK5 proxy](https://github.com/armon/go-socks5) with special dialer. This dialer checks requested IP if it blocked. If IP blocked it uses [tapdance](https://github.com/sergeyfrolov/gotapdance) dialer. If not it uses default net dialer. Here's how evil roskomnadzor blocking machine is bypassed.

# How to use it

## Using docker

You need [docker](https://www.docker.com/community-edition) and gnu-make util. Just open project root in terminal and enter:
```
$ make
``` 
It will build docker and run it. You only need to configure your browser (and any other soft) to use SOCK5 proxy at `127.0.1.1:8000`.

## Manually build

You can build it by your own using [go](https://golang.org/dl/):
```
$ go install github.com/someanon/rkn-baypasser
```

Now you can use it:
```
$ rkn-baypasser -addr 127.0.1.1:8000
```

# TOR reserve

For reliability you can use TOR proxy as reserve if main tapdance dialer fails. 

## Using docker 

Open `tor-proxy` folder and simply enter:

```
$ make
```

It will build and run TOR proxy server. Now you can use it in your main proxy server.

Restart already installed main proxy server:

```
make restart args=-with-tor
```

Or install it with TOR right away:

```
make args=-with-tor
```

## Manually

Just use `-tor` argument with address to your TOR proxy server:

```
$ rkn-baypasser -addr 127.0.1.1:8000 -tor 127.0.1.1:9150
```

# Docker commands list

Root and `tor-proxy` have `gnu-make` based docker control commands:

* `install`
* `build`
* `start`
* `stop`
* `rebuild`
* `restart`
* `log`

Inspect makefiles to be sure what each command doing.  

## Does it work on Windows?

Sure. All this setups possible to do on Windows. You can write me if you need any help. 