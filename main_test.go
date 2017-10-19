package main

import (
	"fmt"
	"net"
	"testing"
	"time"
)

// Build the test configuration (Without real configuration file)
func buildTestConfiguration(t *testing.T) (configuration Configuration) {
	configuration, err := loadConfiguration()
	if err != nil {
		panic(err)
	}
	return
}

// Test states checks
func TestValidStates(t *testing.T) {
	if isValidState("WRONG") {
		t.Error("WRONG is not a valid state")
	}

	if !isValidState("up") {
		t.Error("up should be a valid state")
	}
}

// Test the dynamic weight worker
func TestDynamicWeightWorker(t *testing.T) {
	configuration := buildTestConfiguration(t)
	ticker := time.NewTicker(time.Duration(configuration.DynamicWeightCPUAverageOnSeconds) * time.Second)
	dynamicWeight := 100.0
	go dynamicWeightWorker(*ticker, configuration, &dynamicWeight)
	time.Sleep(time.Duration(3) * time.Second)
}

// Test the full workflow
func TestDaemonWorkflow(t *testing.T) {
	configuration := buildTestConfiguration(t)
	go main()
	time.Sleep(time.Duration(200) * time.Millisecond)
	processClient(configuration.AdminPort, "up")
	conn, _ := net.Dial("tcp", "127.0.0.1:8888")
	defer func() {
		err := conn.Close()
		if err != nil {
			t.Fatalf("Error when closing the TCP connection to the daemon: %s", err)
		}
	}()
	fmt.Fprint(conn, "")
}
