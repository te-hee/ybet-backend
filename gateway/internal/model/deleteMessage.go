package model

type DeleteMessageRequest struct {
	MessageId string `json:"message_id,omitempty" validate:"required,uuid"`
	UserId    string `json:"-"`
}
