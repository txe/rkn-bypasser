# rkn-bypasser

This is [roskomnadzor's](https://eng.rkn.gov.ru/) blocker bypasser.

# How it works

It loads blocked IPs from [roskomsvoboda](http://reestr.rublacklist.net/api/ips) and run local [SOCKS5 proxy](https://github.com/armon/go-socks5) with special dialer. This dialer checks requested IP if it blocked. If IP blocked it uses dialer through the [TOR](https://www.torproject.org/). If not it uses default net dialer. Here's how evil roskomnadzor blocking machine is bypassed.

# How to use it

## Using docker

You need [docker](https://www.docker.com/community-edition). Just open project root in terminal and enter:
```
$ docker-compose up --build -d
``` 
It will build docker and run it. You only need to configure your browser 
(and any other soft) to use SOCKS5 proxy at `127.0.0.1:8000`.

## Manually build

You can build it by your own using [go](https://golang.org/dl/):
```
$ go install github.com/dimuls/rkn-bypasser
```

Now you can use it:
```
$ rkn-bypasser --bind-addr 127.0.0.1:8000
```

## Does it work on Windows?

Sure. All this setups possible to do on Windows. You can write me if you need any help.
 
# Contacts

Feel free to contact me:

* Email: dimuls@yandex.ru
* Telegram: @dimuls