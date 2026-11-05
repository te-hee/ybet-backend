package e2e

type signInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signInResponse struct {
	Id           string `json:"id"`
	RefreshToken string `json:"refresh_token"`
	AuthToken    string `json:"auth_token"`
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
