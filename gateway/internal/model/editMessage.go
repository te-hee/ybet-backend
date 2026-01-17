package model

type EditMessageRequest struct {
	MessageId string `json:"message_id,omitempty"`
	Content   string `json:"content,omitempty"`
	UserId    string
}
