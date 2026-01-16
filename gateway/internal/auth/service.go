package auth

type Service interface {
	GenerateToken(username string) (string, error)
}

type MinimalService struct {
	client Client
}

func NewMinimalService(client Client) *MinimalService {
	return &MinimalService{
		client: client,
	}
}

func (s MinimalService) GenerateToken(username string) (string, error) {
	return s.client.GenerateToken(username)
}
