package handlers

import (
	"encoding/json"
	"github.com/andre2ar/go-products/internal/dto"
	"github.com/andre2ar/go-products/internal/entity"
	"github.com/andre2ar/go-products/internal/infra/database"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UserRepository database.UserRepositoryInterface
	Jwt            *jwtauth.JWTAuth
	JwtExpiresIn   int
}

func NewUserHandler(userRepository database.UserRepositoryInterface) *UserHandler {
	return &UserHandler{UserRepository: userRepository}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	err = h.UserRepository.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
