package model

type SendMessageRequest struct {
	Content  string `json:"content" validate:"required"`
	UserId   string `json:"-"`
	Username string `json:"-"`
}

type GetHistoryRequest struct {
	Limit uint32 `json:"limit" query:"limit" validate:"required,min=1"`
}
