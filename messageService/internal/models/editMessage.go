package models

type EditMessage struct {
	UserId    string `json:"user_id,omitempty"`
	MessageId string `json:"message_id,omitempty"`
	Content   string `json:"content,omitempty"`
}

type NatsEditMessage struct {
	MessageId string `json:"message_id,omitempty"`
	Content   string `json:"content,omitempty"`
}
