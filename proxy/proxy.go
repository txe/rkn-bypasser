package proxy

import (
	"context"
	"log"
	"net"
	"strings"

	"github.com/armon/go-socks5"
	"github.com/someanon/gotapdance/tapdance"
	"golang.org/x/net/proxy"
)

type dialer struct {
	torAddr string
}

func (d dialer) torDial(ctx context.Context, network, addr string) (net.Conn, error) {
	dialer, err := proxy.SOCKS5(network, d.torAddr, nil, proxy.Direct)
	if err != nil {
		return nil, err
	}
	conn, err := dialer.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (d dialer) dial(ctx context.Context, network, addr string) (net.Conn, error) {
	ip := strings.Split(addr, ":")[0]
	if blockedIPs.Has(ip) {
		conn, err := tapdance.Dial(network, addr)
		if err != nil && d.torAddr != "" {
			return d.torDial(ctx, network, addr)
		}
		return conn, err
	}
	return net.Dial(network, addr)
}

func Run(addr string, torAddr string) {
	initBlockedIPs()
	d := dialer{torAddr: torAddr}
	server, err := socks5.New(&socks5.Config{Dial: d.dial})
	if err != nil {
		log.Fatalln("[ERR] Fail to create SOCK5 proxy server: " + err.Error())
	}
	if err := server.ListenAndServe("tcp", addr); err != nil {
		log.Fatalln("[ERR] Fail to start SOCK5 proxy server: " + err.Error())
	}
}
