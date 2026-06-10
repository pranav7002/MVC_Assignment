package controller

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type AuthServiceInterface interface {
	RegisterUser(username, password string) error
	LoginUser(username, password string) (bool, error)
	CreateToken(userID string) (string, error)
}

type AuthController struct {
	AuthService AuthServiceInterface
}

type UserRequestBody struct {
	Username string `json:"username"`
    Password string `json:"password"`
}

func (authController *AuthController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
    userReqBody := new(UserRequestBody)

	err := json.NewDecoder(r.Body).Decode(userReqBody)
	if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Please provide the correct input!!"))
        return
	}

	err = authController.AuthService.RegisterUser(userReqBody.Username, userReqBody.Password)
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Something bad happened on the server :/"))
		return
	}

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("User registered successfully!"))
}

func (authController *AuthController) LoginHandler(w http.ResponseWriter, r *http.Request) {
    userReqBody := new(UserRequestBody)	

	err := json.NewDecoder(r.Body).Decode(userReqBody)
	if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Please provide the correct input!!"))
        return
	}

	isAuthenticated, err := authController.AuthService.LoginUser(userReqBody.Username, userReqBody.Password)
	if err == pgx.ErrNoRows {
		w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("Incorrect username entered!!"))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Something bad happened on the server :/"))
		return
	} 

	if !isAuthenticated {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Incorrect password please check again"))
        return
	}

	tokenString, err := authController.AuthService.CreateToken(userReqBody.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Something bad happened on the server :/"))
		return
	}

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(tokenString))
}