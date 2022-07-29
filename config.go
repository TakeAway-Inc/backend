package main

import (
	"os"

	"backend/db"
	"backend/server"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DB     db.Config     `yaml:"db"`
	Server server.Config `yaml:"server"`
}

func ReadConfig(fileName string) (Config, error) {
	var cnf Config
	data, err := os.ReadFile(fileName)
	if err != nil {
		return Config{}, err
	}
	err = yaml.Unmarshal(data, &cnf)
	if err != nil {
		return Config{}, err
	}
	return cnf, nil
}
