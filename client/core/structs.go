package core

import (
	"crypto/rsa"
	"net/http"
	"time"
)

type Identification struct {
	Id        string        `json:"id"`
	Hostname  string        `json:"hostname"`
	PublicKey rsa.PublicKey `json:"pubkey"`
}

type IdentificationResult struct {
	PublicKey rsa.PublicKey `json:"pubkey"`
}

type ServerIdentity = IdentificationResult

type ClientConfig struct {
	ServerHost      string
	CommunicationEp string
	RecognitionEp   string
}

type Security struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  rsa.PublicKey
}

type Client struct {
	Id       string
	Hostname string

	HttpClient http.Client

	Server ServerIdentity

	Config   ClientConfig
	Security Security
	History  []Command
}

type CommandRequest struct {
	Id string
}

type EncryptedCommand struct {
	Id               string `json:"commandId"`
	EncryptedContent []byte `json:"encryptedContent"`
}

type EncryptedCommandResult struct {
	Id              string
	EncryptedOutput []byte
}

type Command struct {
	Id         string
	Content    string
	ReceivedAt time.Time
}
