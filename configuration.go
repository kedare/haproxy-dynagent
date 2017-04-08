package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Configuration represents the parameter defined
// in the configuration file (or for testing)
type Configuration struct {
	ListenPort                       int
	AdminPort                        int
	DefaultState                     string
	ReportDynamicWeight              bool
	DynamicWeightCPUAverageOnSeconds uint
}

// Loads the configuration file or generate the equivalent for testing
func loadConfiguration() Configuration {
	configuration := &Configuration{}
	if flag.Lookup("test.v") != nil {
		configuration.AdminPort = 8889
		configuration.ListenPort = 8888
		configuration.ReportDynamicWeight = true
		configuration.DynamicWeightCPUAverageOnSeconds = 1
	} else {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		configurationFilePath := dir + "/config.toml"
		log.Printf("Loading configuration from %v", configurationFilePath)
		toml.DecodeFile(configurationFilePath, &configuration)
	}
	return *configuration
}
