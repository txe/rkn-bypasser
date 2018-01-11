package main

import (
	"log"
	"os"
	"time"
)

type logWriter struct{}

func (w logWriter) Write(bytes []byte) (int, error) {
	return os.Stdout.WriteString(time.Now().UTC().Format("2006-01-02 15:04:05.999999 MST ") + string(bytes))
}

func initLog() {
	log.SetFlags(0)
	log.SetOutput(logWriter{})
}
