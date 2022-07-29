package main

import (
	"flag"

	"backend/db"
	"backend/server"

	log "github.com/sirupsen/logrus"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "config.yaml", "Set path to config file.")
	flag.Parse()
	config, err := ReadConfig(configPath)
	if err != nil {
		log.WithError(err).Fatal("can't configure from config file")
	}

	database, err := db.NewDB(config.DB)
	if err != nil {
		log.WithError(err).Fatal("can't create db")
	}

	s := server.NewServer(server.WithDB(database))
	s.AddHTTPServer(config.Server)
	s.Run()
}
