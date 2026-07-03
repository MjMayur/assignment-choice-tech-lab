package main

import (
	"project/pkg/log"
)

func main() {

	err := initConfig()
	if err != nil {
		log.Print(err)
	}

	registry, err := initServer(&config)
	if err != nil {
		log.Print(err)
	}

	err = registry.StartServer()
	if err != nil {
		log.Print(err)
	}
}
