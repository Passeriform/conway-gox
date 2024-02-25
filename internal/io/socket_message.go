package io

type SocketMessage struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload,omitempty"`
}
