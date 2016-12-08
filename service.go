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
	var port = 8888
	var defaultState = "up"
	go processAgent(port, defaultState)
	return nil
}

// Shutdown hook
func (p *HAProxyDynAgent) Stop(s service.Service) error {
	if service.Interactive() {
		os.Exit(0)
	} else if p.cmd.ProcessState.Exited() == false {
		p.cmd.Process.Kill()
	}
	return nil
}
