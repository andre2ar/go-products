package handlers

import (
	"encoding/json"
	"github.com/andre2ar/go-products/internal/dto"
	"github.com/andre2ar/go-products/internal/entity"
	"github.com/andre2ar/go-products/internal/infra/database"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
	"time"
)

type Error struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UserRepository database.UserRepositoryInterface
}

func NewUserHandler(userRepository database.UserRepositoryInterface) *UserHandler {
	return &UserHandler{UserRepository: userRepository}
}

func (h *UserHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("Jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("JwtExpiresIn").(int)

	var loginCredentials dto.LoginCredentialsInput
	err := json.NewDecoder(r.Body).Decode(&loginCredentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.UserRepository.FindByEmail(loginCredentials.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	if !user.ValidatePassword(loginCredentials.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, tokenString, _ := jwt.Encode(map[string]interface{}{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})

	accessToken := dto.AuthResponse{AccessToken: tokenString}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
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
