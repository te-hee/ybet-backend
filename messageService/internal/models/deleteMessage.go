package models

type DeleteMessage struct {
	UserId    string `json:"user_id,omitempty"`
	MessageId string `json:"message_id,omitempty"`
}

type NatsDdeleteMessage struct {
	MessageId string `json:"message_id,omitempty"`
}
