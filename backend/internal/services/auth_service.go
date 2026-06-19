package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepositoryInterface interface {
	InsertUser(username, hash string) error
	GetAttributeFromUsername(username, column string) (string, error)
}

type AuthService struct {
	UserRepo  UserRepositoryInterface
	SecretKey []byte
}

func (s *AuthService) RegisterUser(username, password string) error {
	hash, err := getHashPassword(password)
	if err != nil {
		return err
	}

	err = s.UserRepo.InsertUser(username, hash)
	return err
}

func (s *AuthService) LoginUser(username, password string) (bool, error) {
	hash, err := s.UserRepo.GetAttributeFromUsername(username, "password_hash")
	if err != nil {
		return false, err
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

func (s *AuthService) CreateToken(username string) (string, error) {
	userID, err := s.UserRepo.GetAttributeFromUsername(username, "id")
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(s.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
