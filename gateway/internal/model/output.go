package model

type OutputError struct {
	Error string `json:"error"`
}

type OutputGetHistory struct {
	Messages []Message `json:"messages"`
}

type OutputSendMessege struct {
	MessageId string `json:"message_id"`
	Timestamp int64  `json:"timestamp"`
}

func NewOutputError(errorMessage string) OutputError {
	return OutputError{Error: errorMessage}
}

func NewOutputGetHistory(messages []Message) OutputGetHistory {
	return OutputGetHistory{Messages: messages}
}
