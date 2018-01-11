package main

import (
	"context"
	"log"
	"net"
	"os"
	"strings"

	"github.com/armon/go-socks5"
	"github.com/someanon/gotapdance/tapdance"
)

func dial(ctx context.Context, network, addr string) (net.Conn, error) {
	ip := strings.Split(addr, ":")[0]
	if blockedIPs.Has(ip) {
		return tapdance.Dial(network, addr)
	}
	return net.Dial(network, addr)
}

func main() {
	initLog()

	var addr = os.Getenv("ADDR")
	if addr == "" {
		addr = "127.0.1.1:8000"
	}

	initBlockedIPs()

	server, err := socks5.New(&socks5.Config{Dial: dial})
	if err != nil {
		log.Fatalln("[ERR] Fail to create SOCK5 proxy server: " + err.Error())
	}

	if err := server.ListenAndServe("tcp", addr); err != nil {
		log.Fatalln("[ERR] Fail to start SOCK5 proxy server: " + err.Error())
	}
}
