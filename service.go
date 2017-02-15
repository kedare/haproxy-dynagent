package main

import (
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
	configuration := loadConfiguration()
	port := configuration.ListenPort
	adminPort := configuration.AdminPort
	defaultState := configuration.DefaultState
	reportDynamicWeight := configuration.ReportDynamicWeight
	go processAgent(port, adminPort, defaultState, reportDynamicWeight)
	return nil
}

// Shutdown hook
func (p *HAProxyDynAgent) Stop(s service.Service) error {
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}
