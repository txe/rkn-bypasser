// +build !windows

package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Usage = "RNK bypasser proxy server"
	app.Version = "0.04"
	app.Author = "Vadim Chernov"
	app.Email = "dimuls@yandex.ru"

	app.Commands = []cli.Command{
		{
			Name:  "run",
			Usage: "run proxy server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "addr",
					Usage:  "bind address",
					EnvVar: "ADDR",
				},
				cli.StringFlag{
					Name:   "tor",
					Usage:  "reserve tor proxy server address",
					EnvVar: "TOR",
				},
			},
			Action: run,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
