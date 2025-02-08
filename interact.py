# simple python script to send commands and receive output.
import requests, time, json

conf = json.load(open('config.json', 'r'))

while True:
    try:
        cmd = input('Coyote> ').strip()

        req = requests.post(f'http://127.0.0.1/{conf['sendCommandEp']}', json = {"command": cmd})
        sendCmdInfo = req.json()

        if sendCmdInfo['success']:

            print('Sent Command!')
            time.sleep(3)

            req = requests.get(f'http://127.0.0.1/{conf['outputEp']}')

            outputData = req.json()
            print(f'Output : {outputData['output'].strip()}')

    except (KeyboardInterrupt):
        exit(0)