package models

type Message struct {
	Id        string `json:"id,omitempty"`
	Message   string `json:"message"`
	Timestamp uint64 `json:"timestamp,omitempty"`
}
