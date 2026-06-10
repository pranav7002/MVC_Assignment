package controller

import "net/http"

type AuthServiceInterface interface {
	RegisterUser(username, password string) error
}

type AuthController struct {
	AuthService AuthServiceInterface
}

func (authController *AuthController) RegisterHandler(w http.ResponseWriter, r *http.Request) {

}