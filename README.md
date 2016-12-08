# HAProxy DynAgent

## Introuction

The DynAgent is a simple haproxy agent responding to the requests from the agent-check.

It allows the backend to control its status on the haproxy, to remove itself cleanly without breaking the connections.

The related haproxy documentation can be found here : 

https://cbonte.github.io/haproxy-dconv/configuration-1.5.html#5.2-agent-check

## Usage

You just need to take the .exe and the .ps1, there are no runtime dependencies

You need to first install the service using the ```install.ps1``` script. 

Then, make sure it has been installed successfully and is running.

Your server will by default have the state "up".

To change the state, you need to run the command with as parameter the state you want to put, examples :

```
> .\haproxy-dynagent.exe maint
2016/12/08 16:56:43 Sent new state 'maint' to the agent
> .\haproxy-dynagent.exe down
2016/12/08 16:56:45 Sent new state 'down' to the agent
```

Then make sure you have configured HAProxy to use the agent-check on the specified port (8888)