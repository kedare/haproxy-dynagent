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
		processClient(configuration.AdminPort)
	}

}

// Run when the binary is running with the "-agent" flag
func processAgent(port int, adminPort int, defaultState string, reportDynamicWeight bool) {
	listenAddress := fmt.Sprintf("0.0.0.0:%v", port)
	log.Printf("HAPROXY DynAgent listening on port %d\n", port)
	log.Printf("Administrative interface on port %d\n", adminPort)
	listener, err := net.Listen("tcp", listenAddress)
	log.Printf("Default state to '%s'", defaultState)
	state := defaultState
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	go startAdministrativeInterface(adminPort, &state)

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go routeRequest(conn, &state, reportDynamicWeight)
	}
}

// Run when the binary is running without the "-agent" flag, meaning we are using it as client
func processClient(adminPort int) {
	if len(flag.Args()) < 1 {
		log.Fatal("You need to pass the desired state as parameter")
	} else {
		log.Println("Sending new state to the agent")
	}
	state := flag.Args()[0]
	if isValidState(state) {
		adminUrl := fmt.Sprintf("http://127.0.0.1:%v/", adminPort)
		payload := url.Values{}
		payload.Add("state", state)
		res, err := http.PostForm(adminUrl, payload)
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
func routeRequest(conn net.Conn, state *string, reportDynamicWeight bool) {
	defer conn.Close()
	handleHaproxyRequest(conn, state, reportDynamicWeight)
}

// Handle the HAProxy connection, returning it the state
func handleHaproxyRequest(conn net.Conn, state *string, reportDynamicWeight bool) {
	log.Printf("Got health request from %s\n", conn.RemoteAddr().String())
	var response string
	var ratio float64
	if reportDynamicWeight {
		cpuUsage, _ := cpu.Percent(time.Duration(1)*time.Second, false)
		ratio = 100 - cpuUsage[0]
	} else {
		ratio = 100
	}

	if *state == "up" || *state == "ready" {
		response = fmt.Sprintf("%s,%.0f%%\n", *state, ratio)
	} else {
		response = fmt.Sprintf("%s\n", *state)
	}

	log.Printf("Replying with current active state '%s'\n", strings.TrimSpace(response))

	conn.Write([]byte(response))
}

// Check that the state entered in the administrative state is valid
func isValidState(state string) bool {
	VALID_STATES := []string{
		"ready",
		"up",
		"drain",
		"maint",
		"down",
		"failed",
		"stopped",
	}
	for _, validState := range VALID_STATES {
		if validState == state {
			return true
		}
	}
	return false
}
