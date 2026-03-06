package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/RohitSadawarte79/go-http-framework/internal/domain"
	"github.com/RohitSadawarte79/go-http-framework/internal/service"
)

func urlParam(r *http.Request, name string) string {
	params, ok := r.Context().Value(domain.ParamKey).(map[string]string)

	if !ok {
		return ""
	}

	return params[name]
}

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := urlParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid User Id.", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetByID(id)

	if errors.Is(err, domain.ErrUserNotFound) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	JSON(w, http.StatusOK, user)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	userList, err := h.service.List()

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	JSON(w, http.StatusOK, userList)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}

	err := h.service.Create(&user)

	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			http.Error(w, "validation failed", http.StatusBadRequest)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	JSON(w, http.StatusCreated, user)
}
