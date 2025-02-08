
# Coyote

A Windows HTTP Reverse Shell With E2EE. Using RSA Encryption.

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
We can send commands to a Coyote client with ``interact.py``.

## Output

```
$ python interact.py
Coyote> echo hello!
Sent Command!
Output : hello!
Coyote> cd
Sent Command!
Output : C:\Users\pseco\OneDrive\Desktop\Coyote\client
Coyote> 
```
## Improving It
You can use this project as a base for a secure HTTP-Based remote administration tool, Right now it can only execute simple windows shell commands, You can also make a web panel for communication instead of using a python script, As long as you can send requests to the Coyote server, You can do everything.
