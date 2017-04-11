package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kardianos/service"

	"github.com/shirou/gopsutil/cpu"
)

// Main entrypoint
func main() {
	version := 2
	log.Printf("HAProxy DynAgent v%v", version)
	configuration := loadConfiguration()
	flag.Parse()

	if len(flag.Args()) < 1 {
		serviceConfig := &service.Config{
			Name:        "HAProxyDynAgent",
			DisplayName: "HAProxy DynAgent",
			Description: "State management service for HAProxy",
			Arguments:   []string{},
		}

		haproxyDynAgentService := &HAProxyDynAgent{}

		service, err := service.New(haproxyDynAgentService, serviceConfig)
		if err != nil {
			log.Fatalf("Error setting up service: %v\n", err)
		}
		service.Run()
	} else {
		log.Println("Sending new state to the agent")
		processClient(configuration.AdminPort, flag.Args()[0])
	}

}

// Run when the binary is running without any parameter, meaning we are using it as agent/service
func processAgent(configuration Configuration) {
	listenAddress := fmt.Sprintf("0.0.0.0:%v", configuration.ListenPort)
	log.Printf("HAPROXY DynAgent listening on port %d\n", configuration.ListenPort)
	log.Printf("Administrative interface on port %d\n", configuration.AdminPort)
	listener, err := net.Listen("tcp", listenAddress)
	log.Printf("Default state to '%s'", configuration.DefaultState)
	state := configuration.DefaultState
	dynamicWeight := 100.0
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	go startAdministrativeInterface(configuration, &state)
	if configuration.ReportDynamicWeight {
		log.Println("Dynamic Weight Reporting enabled")
		ticker := time.NewTicker(time.Duration(configuration.DynamicWeightCPUAverageOnSeconds) * time.Second)
		go dynamicWeightWorker(*ticker, configuration, &dynamicWeight)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go routeRequest(conn, &state, &dynamicWeight, configuration)
	}
}

// Run when the binary is running with a parameter, meaning we are using it as client
func processClient(adminPort int, newState string) {
	state := newState
	if isValidState(state) {
		adminURL := fmt.Sprintf("http://127.0.0.1:%v/", adminPort)
		payload := url.Values{}
		payload.Add("state", state)
		res, err := http.PostForm(adminURL, payload)
		if err != nil {
			log.Fatalln("Failed to set state", err)
		} else if res.StatusCode != 200 {
			log.Fatalln("Got invalid response code", res.StatusCode)
		} else {
			log.Println("New state configuration done")
		}
	} else {
		log.Fatalln("Invalid state provided")
	}
}

// Route the incoming TCP request either to the HAPROXY mode to send the status
// Or to the administrative mode to set the state
func routeRequest(conn net.Conn, state *string, dynamicWeight *float64, configuration Configuration) {
	defer conn.Close()
	handleHaproxyRequest(conn, state, dynamicWeight, configuration)
}

// Handle the HAProxy connection, returning it the state
func handleHaproxyRequest(conn net.Conn, state *string, dynamicWeight *float64, configuration Configuration) {
	log.Printf("Got health request from %s\n", conn.RemoteAddr().String())
	var response string
	var ratio float64
	if configuration.ReportDynamicWeight {
		ratio = *dynamicWeight
	} else {
		ratio = 100
	}

	if *state == "up" || *state == "ready" {
		response = fmt.Sprintf("%s %.0f%%\n", *state, ratio)
	} else {
		response = fmt.Sprintf("%s\n", *state)
	}

	log.Printf("Replying with current active state '%s'\n", strings.TrimSpace(response))

	conn.Write([]byte(response))
}

// Check that the state entered in the administrative state is valid
func isValidState(state string) bool {
	ValidStates := []string{
		"ready",
		"up",
		"drain",
		"maint",
		"down",
		"failed",
		"stopped",
	}
	for _, validState := range ValidStates {
		if validState == state {
			return true
		}
	}
	return false
}

// This function is called at regular interval to update the dynamic weight value
func dynamicWeightWorker(ticker time.Ticker, configuration Configuration, dynamicWeight *float64) {
	for {
		select {
		case <-ticker.C:
			cpuUsage, _ := cpu.Percent(time.Duration(configuration.DynamicWeightCPUAverageOnSeconds)*time.Second, false)
			*dynamicWeight = 100 - cpuUsage[0]/2
			log.Printf("Calculated dynamic weight for the last %v seconds: %.0f%%", configuration.DynamicWeightCPUAverageOnSeconds, *dynamicWeight)
		}
	}
}
