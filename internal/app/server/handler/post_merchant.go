package handler

import (
	"net/http"

	"github.com.br/gibranct/simplified-wallet/internal/app/usecase"
	"github.com.br/gibranct/simplified-wallet/internal/domain/vo"
)

type PostMerchantRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	CNPJ     string `json:"cnpj"`
}

func (h handler) PostMerchant(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.otel.Start(r.Context(), "PostMerchant")
	defer span.End()

	var input PostMerchantRequest

	err := h.readJSON(w, r, &input)
	if err != nil {
		err = h.writeJson(w, http.StatusBadRequest, envelope{"error": err.Error()}, nil)
		if err != nil {
			h.logger.Println(err)
		}
		return
	}

	userID, err := h.createUser.Execute(ctx, usecase.CreateUserInput{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Document: input.CNPJ,
		UserType: vo.MerchantUserType,
	})

	if err != nil {
		err = h.writeJson(w, http.StatusUnprocessableEntity, envelope{"error": err.Error()}, nil)
		if err != nil {
			h.logger.Println(err)
		}
		return
	}

	err = h.writeJson(w, http.StatusCreated, envelope{"user_id": userID}, nil)
	if err != nil {
		err = h.writeJson(w, http.StatusInternalServerError, envelope{"error": "failed to write response"}, nil)
		if err != nil {
			h.logger.Println(err)
		}
		return
	}
}
