package main

import (
	"flag"
	"log"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to configuration yml")
	flag.Parse()

	if len(configPath) == 0 {
		log.Fatal("Please provide a configuration path")
	}
	config, err := readConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	start(config)
}
