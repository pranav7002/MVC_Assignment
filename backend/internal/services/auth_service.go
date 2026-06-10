package services

type UserRepositoryInterface interface {
	InsertUser(username, hash string) error
}

type AuthService struct {
	UserRepo UserRepositoryInterface
}

func (authService *AuthService) RegisterUser(username, password string) error {
	return nil
}