package main

import (
	"flag"
	"log"

	"github.com/cxrdevelop/optimization_engine/api_server/config"
	"github.com/cxrdevelop/optimization_engine/api_server/server"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "run/config.yml", "path to the config file")
	flag.Parse()

	if cfg, err := config.ReadConfig(configPath); err != nil {
		log.Panicf("error reading config: %s", err)
	} else {
		prometheus.MustRegister(prometheus.NewBuildInfoCollector())
		server.New(cfg).Start()
	}
}
