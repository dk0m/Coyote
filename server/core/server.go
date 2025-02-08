package core

import (
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

type ServerConfig struct {
	ServerHost      string `json:"serverHost"`
	Port            int    `json:"port"`
	CommunicationEp string `json:"communicationEp"`
	RecognitionEp   string `json:"recognitionEp"`
}
type Security struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  rsa.PublicKey
}

type Identification struct {
	Id        string        `json:"id"`
	Hostname  string        `json:"hostname"`
	PublicKey rsa.PublicKey `json:"pubkey"`
}

type IdentificationResult struct {
	PublicKey rsa.PublicKey
}
type EncryptedCommand struct {
	EncryptedContent []byte
}

type CommandRequest struct {
	Id string
}

type SendCommandInfo struct {
	Command string `json:"command"`
}

type Command struct {
	Id         string
	Content    string
	ReceivedAt time.Time
}

type Server struct {
	Config    ServerConfig
	Security  Security
	Router    *gin.Engine
	Clients   map[string]Identification
	IssuedCmd Command
}

var FAILURE_MESSAGE = gin.H{"success": false}

func (server *Server) ConfigServer() error {
	configFile, err := os.Open("../config.json")

	if err != nil {
		return err
	}

	gin.Default()
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
	server.IssuedCmd = Command{Content: "N/A"}

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

		if err != nil {
			context.JSON(http.StatusUnauthorized, FAILURE_MESSAGE)
		}

		context.JSON(http.StatusOK, gin.H{"commandId": command.Id, "encryptedContent": encryptedBytes})
	})

	router.POST("/sendCommand", func(context *gin.Context) {
		sendCmdInfo := SendCommandInfo{}

		err := context.BindJSON(&sendCmdInfo)

		if err != nil {
			context.JSON(http.StatusUnauthorized, FAILURE_MESSAGE)
		}

		// in the next version, it will be a queue instead of just a structure
		server.IssuedCmd = Command{
			Id:         uuid.New().String(),
			Content:    sendCmdInfo.Command,
			ReceivedAt: time.Now(),
		}

	})

	server.Router = router

	return nil
}

func (server *Server) StartServer() {
	router := server.Router
	config := server.Config

	router.Run(fmt.Sprintf("%s:%d", config.ServerHost, config.Port))
}
