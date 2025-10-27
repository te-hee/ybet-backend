package model

type InputMessage struct{
	Content string `json:"content"`
}

type InputHistory struct{
	Limit uint32 `json:"limit"`
}
