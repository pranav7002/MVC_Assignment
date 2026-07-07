package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/pranav7002/MVC_Assignment/internal/middleware"
)

type AuthServiceInterface interface {
	RegisterUser(username, password string) error
	LoginUser(username, password string) (bool, error)
	CreateToken(userID string) (string, error)
	CreateWSTicket(userID string) (string, error)
}

type AuthController struct {
	AuthService AuthServiceInterface
}

type UserRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *AuthController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	userReqBody := new(UserRequestBody)

	err := json.NewDecoder(r.Body).Decode(userReqBody)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Please provide the correct input!!")
		return
	}

	userReqBody.Username = strings.TrimSpace(userReqBody.Username)

	if userReqBody.Username == "" {
		WriteError(w, http.StatusBadRequest, "Username can't be empty!")
		return
	}
	if len(userReqBody.Username) > 50 {
		WriteError(w, http.StatusBadRequest, "Username longer than 50 characters!")
		return
	}

	if userReqBody.Password == "" {
		WriteError(w, http.StatusBadRequest, "Password can't be empty!")
		return	
	}

	err = c.AuthService.RegisterUser(userReqBody.Username, userReqBody.Password)

	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Something bad happened on the server :/")
		return
	}

	WriteJSON(w, http.StatusOK, "User registered successfully!")
}

func (c *AuthController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	userReqBody := new(UserRequestBody)

	err := json.NewDecoder(r.Body).Decode(userReqBody)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Please provide the correct input!!")
		return
	}

	userReqBody.Username = strings.TrimSpace(userReqBody.Username)

	isAuthenticated, err := c.AuthService.LoginUser(userReqBody.Username, userReqBody.Password)
	if err == pgx.ErrNoRows {
		WriteError(w, http.StatusUnauthorized, "Incorrect username entered!!")
		return
	} else if err != nil {
		WriteError(w, http.StatusInternalServerError, "Something bad happened on the server :/")
		return
	}

	if !isAuthenticated {
		WriteError(w, http.StatusBadRequest, "Incorrect password please check again")
		return
	}

	tokenString, err := c.AuthService.CreateToken(userReqBody.Username)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Something bad happened on the server :/")
		return
	}

	WriteJSON(w, http.StatusCreated, tokenString)
}

func (c *AuthController) GenerateWSTicketHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.ContextKey("user_id")).(string)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	WSTicketString, err := c.AuthService.CreateWSTicket(userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Something bad happened on the server :/")
		return
	}

	WriteJSON(w, http.StatusCreated, WSTicketString)
}
