# HAProxy DynAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/kedare/haproxy-dynagent)](https://goreportcard.com/report/github.com/kedare/haproxy-dynagent)

## Introduction

The DynAgent is a simple haproxy agent responding to the requests from the agent-check.

It allows to report a dynamic weight (This is optional (ReportDynamicWeight)) so nodes with heavy CPU with have a lower priority (The minimum is 50% of the default priority, the maximum being 100%, this is based on the CPU average of a specified time (DynamicWeightCPUAverageOnSeconds))

It allows the backend to control its status on the haproxy, to remove itself cleanly without breaking the connections.

The state control is only available from the same host on the configured AdminPort, you can control it either by using the same binary as client when passing the desired state as paramter, or through HTTP (check admin.go), it also have a tiny web interface for troubleshooting to get/set the state (On the AdminPort)

The related haproxy documentation can be found here :

https://cbonte.github.io/haproxy-dconv/configuration-1.5.html#5.2-agent-check

## Usage

You just need to take the .exe, the .toml and the .ps1, there are no runtime dependencies

You need to first install the service using the ```install.ps1``` script. (On Windows)
Or create a simple systemd configuration file to run the process on Linux

Then, make sure it has been installed successfully and is running.

Your server will by default have the state "up", this can be changed in the .toml configuration file (Check the valid options in the HAProxy documentation)

To change the state, you need to run the command with as parameter the state you want to put, examples :

```
> .\haproxy-dynagent.exe maint
2016/12/08 16:56:43 Sent new state 'maint' to the agent
> .\haproxy-dynagent.exe down
2016/12/08 16:56:45 Sent new state 'down' to the agent
```

Then make sure you have configured HAProxy to use the agent-check on the specified port (8888)
