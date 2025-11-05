package models

type ListUpdateAction = string
type MessageType = string

const (
	Disconnect ListUpdateAction = "disconnect"
	Connect    ListUpdateAction = "connect"

	SystemMessageType  MessageType = "systemMessage"
	UserMessageType    MessageType = "userMessage"
	UserListUpdateType MessageType = "userListUpdate"
)

type Message struct {
	Type    MessageType `json:"type"`
	Payload any         `json:"payload"`
}

type UserMessage struct {
	Id        string `json:"id"`
	Content   string `json:"content"`
	Timestamp uint64 `json:"timestamp"`
}

type SystemMessage struct {
	Content string `json:"content"`
}

type UserListUpdate struct {
	Action ListUpdateAction `json:"action"`
	Uuid   string           `json:"uuid"`
}
