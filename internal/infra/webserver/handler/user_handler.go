package handler

import (
	"encoding/json"
	"net/http"

	"github.com/FreitasGabriel/fullcycle-api/internal/dto"
	"github.com/FreitasGabriel/fullcycle-api/internal/entity"
	"github.com/FreitasGabriel/fullcycle-api/internal/infra/database"
	"github.com/go-chi/jwtauth"
)

type UserHandler struct {
	UserDB       database.UserInterface
	JWT          *jwtauth.JWTAuth
	JWTExpiresIn int
}

func NewUserHandler(userDB database.UserInterface) *UserHandler {
	return &UserHandler{
		UserDB: userDB,
	}
}

func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = uh.UserDB.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}
