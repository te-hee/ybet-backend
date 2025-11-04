package model

type Output struct {
	Success bool `json:"success"`
}

type OutputError struct {
	Output
	Error string `json:"error"`
}

type OutputGetHistory struct {
	Output
	Messages []Message `json:"messages"`
}

type OutputSendMessege struct {
	Output
}

func NewOutput(success bool) Output {
	return Output{Success: success}
}

func NewOutputError(errorMessage string) OutputError {
	return OutputError{Output: NewOutput(false), Error: errorMessage}
}

func NewOutputGetHistory(messages []Message) OutputGetHistory {
	return OutputGetHistory{Output: NewOutput(true), Messages: messages}
}

func NewOutputSendMessage() OutputSendMessege {
	return OutputSendMessege{Output: NewOutput(true)}
}
