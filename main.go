package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/kardianos/service"

	"github.com/shirou/gopsutil/cpu"
)

// Main entrypoint
func main() {
	var port = 8888
	var defaultState = "up"
	flag.Parse()

	if len(flag.Args()) < 1 {
		var serviceConfig = &service.Config{
			Name:        "HAProxyDynAgent",
			DisplayName: "HAProxy DynAgent",
			Description: "State management service for HAProxy",
			Arguments: []string{
				"-port", string(port),
				"-state", defaultState,
			},
		}

		var haproxyDynAgentService = &HAProxyDynAgent{}

		var service, err = service.New(haproxyDynAgentService, serviceConfig)
		if err != nil {
			log.Fatalf("Error setting up service: %v\n", err)
		}
		service.Run()
	} else {
		processClient(port)
	}

}

// Run when the binary is running with the "-agent" flag
func processAgent(port int, defaultState string) {
	var listenAddress = fmt.Sprintf("0.0.0.0:%v", port)
	log.Printf("HAPROXY DynAgent listening on %s\n", listenAddress)
	var listener, err = net.Listen("tcp", listenAddress)
	log.Printf("Default state to '%s'", defaultState)
	var state = defaultState
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		var conn, err = listener.Accept()
		if err != nil {
			panic(err)
		}
		go routeRequest(conn, &state)
	}
}

// Run when the binary is running without the "-agent" flag, meaning we are using it as client
func processClient(port int) {
	if len(flag.Args()) < 1 {
		log.Fatal("You need to pass the desired state as parameter")
	}
	var state = flag.Args()[0]
	if isValidState(state) {
		var listenAddress = fmt.Sprintf("127.0.0.1:%v", port)
		var conn, err = net.Dial("tcp", listenAddress)
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		fmt.Fprintf(conn, "%s\n", state)
		log.Printf("Sent new state '%v' to the agent\n", state)
	} else {
		log.Fatalln("Invalid state provided")
	}
}

// Route the incoming TCP request either to the HAPROXY mode to send the status
// Or to the administrative mode to set the state
func routeRequest(conn net.Conn, state *string) {
	defer conn.Close()
	if strings.Contains(conn.RemoteAddr().String(), "127.0.0.1") {
		handleAdministrativeRequest(conn, state)
	} else {
		handleHaproxyRequest(conn, state)
	}
}

// Handle the administrative connection that is allowed to change the state
func handleAdministrativeRequest(conn net.Conn, state *string) {
	log.Printf("Got administrative connection from %s\n", conn.RemoteAddr().String())
	fmt.Fprintf(conn, "Current state: %v\r\n Please enter new state: ", *state)
	var buffer, err = bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		log.Fatalf("Error reading buffer from administrative connection: %v", err)
	}
	var request = strings.Trim(string(buffer), " \r\n")
	if isValidState(request) {
		log.Printf("Switching state to '%v'", request)
		*state = request
	} else {
		log.Printf("Got invalid state '%v'", request)
	}
}

// Handle the HAProxy connection, returning it the state
func handleHaproxyRequest(conn net.Conn, state *string) {
	log.Printf("Got health request from %s\n", conn.RemoteAddr().String())
	var response string
	if *state == "up" {
		// If the state is up, we return a percentage that will tell
		// HAProxy how busy if the server (The higher it is, the most requests it will receive)
		var cpuUsage, _ = cpu.Percent(time.Duration(1)*time.Second, false)
		var ratio = 100 - cpuUsage[0]
		response = fmt.Sprintf("%.0f%%", ratio)
	} else {
		// If not, return the state directly
		response = fmt.Sprintf("%s\n", *state)
	}
	log.Printf("Replying with current active state '%s'\n", response)

	conn.Write([]byte(response))
}

// Check that the state entered in the administrative state is valid
func isValidState(state string) bool {
	var VALID_STATES = []string{
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
