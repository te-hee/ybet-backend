package model

type DeleteMessageRequest struct {
	MessageId string `uri:"message_id,omitempty" validate:"required,uuid"`
	UserId    string `json:"-"`
}
