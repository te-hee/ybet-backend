package e2e

type loginRequest struct {
	Username string `json:"username"`
}

type loginResponse struct {
	Token  string `json:"token"`
	UserId string `json:"user_id"`
}

type sendMessageRequest struct {
	Content string `json:"content"`
}

type sendMessageResponse struct {
	Message
}

type Message struct {
	Username  string `json:"username"`
	UserId    string `json:"user_id"`
	Id        string `json:"id"`
	Content   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}
