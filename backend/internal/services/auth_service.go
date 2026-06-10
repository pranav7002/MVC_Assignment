package services

import "golang.org/x/crypto/bcrypt"

type UserRepositoryInterface interface {
	InsertUser(username, hash string) error
}

type AuthService struct {
	UserRepo UserRepositoryInterface
}

func (authService *AuthService) RegisterUser(username, password string) error {
    hash, err := getHashPassword(password)
    if err != nil {
        return err
    }

	err = authService.UserRepo.InsertUser(username, hash)
	return err
}

func getHashPassword(password string) (string, error) {
    bytePassword := []byte(password)
    hash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}