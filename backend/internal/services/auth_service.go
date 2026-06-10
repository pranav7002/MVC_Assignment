package services

import "golang.org/x/crypto/bcrypt"

type UserRepositoryInterface interface {
	InsertUser(username, hash string) error
	GetPasswordHash(username string) (string, error)
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

func (authService *AuthService) LoginUser(username, password string) (bool, error) {
    hash, err := authService.UserRepo.GetPasswordHash(username)
    if err != nil {
        return false,err
    }

	isSame := checkPassword(hash, password)
	return isSame, nil
}

func getHashPassword(password string) (string, error) {
    bytePassword := []byte(password)
    hash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

func checkPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil 
}