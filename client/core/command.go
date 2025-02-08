package core

import (
	"time"
)

type EncryptedCommand struct {
	Id               string `json:"commandId"`
	EncryptedContent []byte `json:"encryptedContent"`
}

type Command struct {
	Id         string
	Content    string
	ReceivedAt time.Time
}
