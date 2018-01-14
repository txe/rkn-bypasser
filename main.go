package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/someanon/rkn-bypasser/proxy"
)

type logWriter struct {
	file *os.File
}

func (w logWriter) Write(bytes []byte) (int, error) {
	return w.file.WriteString(time.Now().Format("2006-01-02 15:04:05.999999 ") + string(bytes))

}

func main() {
	log.SetFlags(0)
	log.SetOutput(logWriter{file: os.Stdout})

	var (
		addr    string
		torAddr string
		withTor bool
	)

	flag.StringVar(&addr, "addr", os.Getenv("ADDR"), "bind address")
	flag.StringVar(&torAddr, "tor", os.Getenv("TOR"), "TOR proxy server address")
	flag.BoolVar(&withTor, "with-tor", false, "use TOR proxy reserve with default address")

	flag.Parse()

	if addr == "" {
		log.Fatal("[ERR] Set ADDR environment variable or -addr flag")
	}

	if torAddr == "" && withTor {
		torAddr = "tor-proxy:9150"
	}

	torStr := "without tor"
	if torAddr != "" {
		torStr = "with tor " + torAddr
	}
	log.Printf("Starting proxy with addr %s and %s\n", addr, torStr)

	proxy.Run(addr, torAddr)
}
