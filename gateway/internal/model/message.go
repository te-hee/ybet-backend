package model

type Message struct {
	MessageId string `json:"message_id"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
	UserId    string `json:"user_id"`
	Username  string `json:"username"`
}
