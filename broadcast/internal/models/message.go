package models

type Message struct {
	Target  string `json:"target"`
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type UserMessage struct {
	Id        string `json:"id,omitempty"`
	Message   string `json:"message"`
	Timestamp uint64 `json:"timestamp,omitempty"`
}
