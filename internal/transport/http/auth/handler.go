package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type UseCase interface {
	SignUp(email, login, password string) error
	SignIn(login, password string) (string, error)
}

type Handler struct {
	useCase  UseCase
	validate *validator.Validate
}

func NewHandler(uc UseCase) *Handler {
	return &Handler{useCase: uc, validate: validator.New()}
}

type SignUpInput struct {
	Email    string `json:"email" validate:"required,email,max=254"`
	Login    string `json:"login" validate:"required,min=2,max=72"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type SignInInput struct {
	Email    string `json:"email" validate:"required,email,max=254"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var reqInput SignUpInput

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
	var reqInput SignInInput
	defer r.Body.Close()
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

	token, err := h.useCase.SignIn(reqInput.Email, reqInput.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(map[string]string{"token": token})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
