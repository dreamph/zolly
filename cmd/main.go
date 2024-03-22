package main

import (
	"flag"
	"github.com/dreamph/zolly"
	"log"
)

func main() {
	configFile := flag.String("c", "config.yml", "config file")
	config, err := zolly.LoadConfig(*configFile)
	if err != nil {
		log.Fatalln(err)
	}

	err = zolly.Start(config)
	if err != nil {
		log.Fatalln(err)
	}
}
