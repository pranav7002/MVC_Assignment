package controller

import (
	"encoding/json"
	"net/http"
)

type AuthServiceInterface interface {
	RegisterUser(username, password string) error
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