package models

type ListUpdateAction = string
type MessageType = string

const (
	Disconnect ListUpdateAction = "disconnect"
	Connect    ListUpdateAction = "connect"

	SystemMessageType  MessageType = "systemMessage"
	UserMessageType    MessageType = "userMessage"
	UserListUpdateType MessageType = "userListUpdate"
	EditMessageType    MessageType = "editMessage"
	DeleteMessageType  MessageType = "deleteMessage"
)

type Message struct {
	Type    MessageType `json:"type"`
	Payload any         `json:"payload"`
}

type UserMessage struct {
	Id        string `json:"message_id"`
	Username  string `json:"username"`
	Message   string `json:"content"`
	Timestamp uint64 `json:"timestamp"`
	UserId    string `json:"user_id"`
}

type EditMessage struct {
	MessageId string `json:"message_id,omitempty"`
	Content   string `json:"content,omitempty"`
}

type DeleteMessage struct {
	MessageId string `json:"message_id,omitempty"`
}

type SystemMessage struct {
	Content string `json:"content"`
}

type UserListUpdate struct {
	Action   ListUpdateAction `json:"action"`
	UserId   string           `json:"user_id"`
	Username string           `json:"username"`
}
