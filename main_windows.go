// +build windows

package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Usage = "RNK bypasser proxy server"
	app.Version = "0.5"
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
		{
			Name:  "service",
			Usage: "Windows service control",
			Subcommands: []cli.Command{
				{
					Name:   "install",
					Usage:  "install and run Windows service",
					Action: installService,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "addr",
							Usage:  "bind address",
							EnvVar: "ADDR",
							Value:  "127.0.1.1:8000",
						},
						cli.StringFlag{
							Name:   "tor",
							Usage:  "reserve tor proxy server address",
							EnvVar: "TOR",
						},
					},
				},
				{
					Name:   "uninstall",
					Usage:  "uninstall Windows service",
					Action: uninstallService,
				},
				{
					Name:   "start",
					Usage:  "start stopped service Windows",
					Action: startService,
				},
				{
					Name:   "run",
					Usage:  "intended to be used by service manager",
					Action: runService,
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
				},
				{
					Name:   "stop",
					Usage:  "stop service immediately",
					Action: stopService,
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
