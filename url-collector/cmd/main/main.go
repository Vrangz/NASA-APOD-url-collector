package main

import (
	"log"
	"url-collector/internal/config"
	"url-collector/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("config could not be loaded.", err)
	}

	s := server.New(cfg)
	log.Fatal(s.Start())
}
