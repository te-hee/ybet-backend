package model

type EditMessageRequest struct {
	MessageId string `uri:"message_id,omitempty" validate:"required,uuid"`
	Content   string `json:"content,omitempty" validate:"required"`
	UserId    string `json:"-"`
}
