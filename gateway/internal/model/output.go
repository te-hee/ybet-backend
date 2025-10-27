package model

type Output struct{
	Success bool
}

type OutputError struct{
	Output
	Error string
}

type OutputGetHistory struct{
	Output
	Messages []Message
}

type OutputSendMessege struct{
	Output
}