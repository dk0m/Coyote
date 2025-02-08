package core

import (
	"crypto/rsa"
	"time"

	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	ServerHost      string `json:"serverHost"`
	Port            int    `json:"port"`
	CommunicationEp string `json:"communicationEp"`
	RecognitionEp   string `json:"recognitionEp"`
	SendCommandEp   string `json:"sendCommandEp"`
	OutputEp        string `json:"outputEp"`
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

type SendCommandInfo struct {
	Command string `json:"command"`
}

type IdentificationResult struct {
	PublicKey rsa.PublicKey
}
type EncryptedCommand struct {
	EncryptedContent []byte
}

type EncryptedCommandResult struct {
	Id              string
	EncryptedOutput []byte
}

type CommandOutput struct {
	Id     string
	Output string
}

type CommandRequest struct {
	Id string
}

type Command struct {
	Id         string
	Content    string
	ReceivedAt time.Time
}

type Server struct {
	Config     ServerConfig
	Security   Security
	Router     *gin.Engine
	Clients    map[string]Identification
	IssuedCmd  *Command
	LastOutput *CommandOutput
}
