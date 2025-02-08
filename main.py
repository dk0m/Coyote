import requests


requests.post('http://127.0.0.1/sendCommand', json = {"command": "echo hi!!!!"})