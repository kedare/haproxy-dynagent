package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/kardianos/service"
)

// Service structure
type HAProxyDynAgent struct {
	service service.Service
	cmd     *exec.Cmd
}

// Startup hook
func (p *HAProxyDynAgent) Start(s service.Service) error {
	configuration, err := loadConfiguration()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}
	go processAgent(configuration)
	return nil
}

// Shutdown hook
func (p *HAProxyDynAgent) Stop(s service.Service) error {
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}
