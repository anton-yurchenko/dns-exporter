package main

import (
	"os"

	"github.com/anton-yurchenko/dns-exporter/internal/app"
	log "github.com/sirupsen/logrus"
)

// Version of an application
const Version string = "1.0.8"

func init() {
	log.SetReportCaller(false)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:            false,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		DisableTimestamp:       false,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	app.Entrypoint(Version)
}
