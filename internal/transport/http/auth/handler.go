package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type UseCase interface {
	SignUp(email, login, password string) error
	SignIn(email, login, password string) error
}

type Handler struct {
	useCase  UseCase
	validate *validator.Validate
}

func NewHandler(uc UseCase) *Handler {
	return &Handler{useCase: uc, validate: validator.New()}
}

type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email,max=254"`
	Login    string `json:"name" validate:"required,min=2,max=72"`
	Password string `json:"password" validate:"min=8,max=64"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var reqInput SignUpRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqInput); err != nil {
		http.Error(w, "Error JSON parse", http.StatusBadRequest)
		return
	}

	err := h.validate.Struct(reqInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.useCase.SignUp(reqInput.Email, reqInput.Login, reqInput.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Регистрация"))
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	h.useCase.SignIn("k", "sd", "sdsd")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Вход"))
}
