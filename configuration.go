package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Configuration struct {
	ListenPort   int
	DefaultState string
}

func loadConfiguration() Configuration {

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	configurationFilePath := dir + "/config.toml"
	log.Printf("Loading configuration from %v", configurationFilePath)
	configuration := &Configuration{}
	toml.DecodeFile(configurationFilePath, &configuration)
	return *configuration
}
