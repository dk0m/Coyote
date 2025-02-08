import requests
# a simple request to change the server's issued command, for testing purposes.
requests.post('http://127.0.0.1/sendCommand', json = {"command": "echo hello!"})
