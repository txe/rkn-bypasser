package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"strings"

	"github.com/armon/go-socks5"
	"golang.org/x/net/proxy"
)

var bindAddr = os.Getenv("BIND_ADDR")
var torProxyAddr = os.Getenv("TOR_PROXY")

func dial(ctx context.Context, network, addr string) (net.Conn, error) {
	ip := strings.Split(addr, ":")[0]
	if blockedIPs.Has(ip) {
		dialer, err := proxy.SOCKS5(network, torProxyAddr, nil, proxy.Direct)
		if err != nil {
			return nil, errors.New("failed to create TOR SOCK5 proxy server dialer: " + err.Error())
		}
		conn, err := dialer.Dial(network, addr)
		if err != nil {
			return nil, errors.New("failed to dial to TOR SOCK5 proxy server: " + err.Error())
		}
		return conn, nil
	}
	return net.Dial(network, addr)
}

func main() {
	initLog()

	if bindAddr == "" {
		log.Fatalln("[ERR] BIND_ADDR environment variable is not set")
	}

	if torProxyAddr == "" {
		log.Fatalln("[ERR] TOR_PROXY environment variable is not set")
	}

	initBlockedIPs()

	server, err := socks5.New(&socks5.Config{Dial: dial})
	if err != nil {
		log.Fatalln("[ERR] Fail to create SOCK5 proxy server: " + err.Error())
	}

	if err := server.ListenAndServe("tcp", bindAddr); err != nil {
		log.Fatalln("[ERR] Fail to start SOCK5 proxy server: " + err.Error())
	}
}
