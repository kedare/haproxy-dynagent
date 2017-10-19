package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/kardianos/service"
)

// HAProxyDynAgent represents the service
type HAProxyDynAgent struct {
	service service.Service
	cmd     *exec.Cmd
}

// Start defines the service startup hook
func (p *HAProxyDynAgent) Start(s service.Service) error {
	configuration, err := loadConfiguration()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}
	go processAgent(configuration)
	return nil
}

// Stop define the service shutdown hook
func (p *HAProxyDynAgent) Stop(s service.Service) error {
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}
