
# Coyote

An HTTP Reverse Shell With E2EE Using RSA Encryption.

## Requirements
You'll need to install these 2 requirements for Coyote to work, Which are [Gin](github.com/gin-gonic/gin) and [Uuid](github.com/google/uuid)
```
go get github.com/gin-gonic/gin
go get github.com/google/uuid
```
## Building & Details

### Client
To build the client, Run this in the client directory
```
go build
```
The client will keep fetching commands in an infinite loop, Executing new ones when they are issued.

### Server 
To build the server, Run this in the server directory
```
go build
```
The server will listen for specific endpoints that will allow the client to identify themselves, fetch commands and send output.

## Sending Commands
We can send commands with a simple python script like this
```python
import requests
requests.post('http://127.0.0.1/sendCommand', json = {"command": "echo hello!"})
```

## Todo
- Allow for shared structs.
- Make the client send the output back (encrypted). 
- Make a program (with nice TUI) to send commands. 
