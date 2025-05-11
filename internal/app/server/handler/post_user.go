package handler

import (
	"net/http"

	"github.com.br/gibranct/simplified-wallet/internal/app/usecase"
)

type PostUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	CPF      string `json:"cpf"`
}

func (h handler) PostUser(w http.ResponseWriter, r *http.Request) {
	var input PostUserRequest

	err := h.readJSON(w, r, &input)
	if err != nil {
		h.writeJson(w, http.StatusBadRequest, envelope{"error": err.Error()}, nil)
		return
	}

	userID, err := h.createUser.Execute(r.Context(), usecase.CreateUserInput{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		CPF:      input.CPF,
	})

	if err != nil {
		h.writeJson(w, http.StatusUnprocessableEntity, envelope{"error": err.Error()}, nil)
		return
	}

	h.writeJson(w, http.StatusCreated, envelope{"user_id": userID}, nil)

	return

}
