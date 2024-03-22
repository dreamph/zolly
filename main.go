package main

import (
	"flag"
	"log"
)

func main() {
	configFile := flag.String("c", "config.yml", "config file")
	config, err := LoadConfig(*configFile)
	if err != nil {
		log.Fatalln(err)
	}

	err = Start(config)
	if err != nil {
		log.Fatalln(err)
	}
}
