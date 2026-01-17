package model

type DeleteMessageRequest struct {
	MessageId string `json:"message_id,omitempty"`
	UserId    string
}
