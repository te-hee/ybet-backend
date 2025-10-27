package model

type InputMessage struct{
	Conetnt string `json:"content"`
}

type InputHistory struct{
	Limit uint32 `json:"limit"`
}