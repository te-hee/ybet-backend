package model

type InputMessage struct {
	Content  string `json:"content"`
	UserId   string
	Username string
}

type InputHistory struct {
	Limit uint32 `json:"limit"`
}
