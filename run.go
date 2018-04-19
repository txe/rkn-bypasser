package main

import (
	"log"

	"github.com/getlantern/errors"
	"github.com/someanon/rkn-bypasser/proxy"
	"github.com/urfave/cli"
)

func run(c *cli.Context) error {
	initLog()

	if !c.IsSet("addr") {
		log.Fatal("[ERR] Set ADDR environment variable or -addr flag")
		return errors.New("addr is not set")
	}

	addr := c.String("addr")
	torAddr := c.String("tor")

	torStr := "without tor"
	if torAddr != "" {
		torStr = "with tor " + torAddr
	}
	log.Printf("Running at %s %s\n", addr, torStr)

	proxy.Run(addr, torAddr)
	return nil
}
