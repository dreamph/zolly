package main

import (
	"flag"

	"log"
)

func main() {
	configFile := flag.String("c", "config.yml", "config file")
	flag.Parse()

	cfg, err := LoadConfig(*configFile)
	if err != nil {
		log.Fatalln(err)
	}

	err = Start(cfg)
	if err != nil {
		log.Fatalln(err)
	}
}
