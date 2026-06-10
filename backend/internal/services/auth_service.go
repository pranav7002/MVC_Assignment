package services

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

type UserRepositoryInterface interface {
	InsertUser(username, hash string) error
	GetAttributeFromUsername(username, column string) (string, error)
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
    hash, err := authService.UserRepo.GetAttributeFromUsername(username, "password_hash")
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

func (authService *AuthService) CreateToken(username string) (string, error) {
	userID, err := authService.UserRepo.GetAttributeFromUsername(username, "user_id")
	if err != nil {
		return "", err
	}

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
        jwt.MapClaims{ 
        "user_id": userID, 
        "exp": time.Now().Add(time.Hour * 24).Unix(), 
        })

    tokenString, err := token.SignedString(secretKey)
    if err != nil {
    return "", err
    }

 return tokenString, nil
}