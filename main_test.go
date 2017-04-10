package main

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func buildTestConfiguration(t *testing.T) (configuration Configuration) {
	configuration = loadConfiguration()
	return
}

func TestValidStates(t *testing.T) {
	if isValidState("WRONG") {
		t.Error("WRONG is not a valid state")
	}

	if !isValidState("up") {
		t.Error("up should be a valid state")
	}
}

func TestDynamicWeightWorker(t *testing.T) {
	configuration := buildTestConfiguration(t)
	ticker := time.NewTicker(time.Duration(configuration.DynamicWeightCPUAverageOnSeconds) * time.Second)
	dynamicWeight := 100.0
	go dynamicWeightWorker(*ticker, configuration, &dynamicWeight)
	time.Sleep(time.Duration(3) * time.Second)
}

func TestDaemonWorkflow(t *testing.T) {
	configuration := buildTestConfiguration(t)
	go main()
	time.Sleep(time.Duration(200) * time.Millisecond)
	processClient(configuration.AdminPort, "up")
	conn, _ := net.Dial("tcp", "127.0.0.1:8888")
	defer conn.Close()
	fmt.Fprint(conn, "")
}
