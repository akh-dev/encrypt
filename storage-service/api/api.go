package api

import (
	"encoding/json"
)

type Response struct {
	StatusCode    int             `json:"status_code"`
	StatusMessage string          `json:"status_message"`
	Result        json.RawMessage `json:"result,omitempty"`
	Errors        []string        `json:"errors,omitempty"`
}

type IdMessage struct {
	Id      string `json:"id"`
	Payload string `json:"payload"`
}

type Id struct {
	Id string `json:"id"`
}
