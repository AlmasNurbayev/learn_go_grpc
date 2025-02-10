package authGrpc

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type LoginRequestForValidate struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=50"`
	Type     string `json:"type" validate:"required"`
	AppId    int32  `json:"app_id" validate:"required,gte=0"`
}

type RegisterRequestForValidate struct {
	// omitempty необходим чтобы работало следующее после него условие проверки
	Email    string `json:"email" validate:"required_without=Phone,omitempty,email"`
	Phone    string `json:"phone" validate:"required_without=Email,omitempty,e164"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}

type IsAdminForValidate struct {
	UserId int64 `json:"user_id" validate:"required,gte=0"`
}

func ValidateStruct(T any) error {

	validate := validator.New()
	err := validate.Struct(T)

	if err != nil {
		fmt.Printf("%+v", err)
		return err
	}

	return nil
}
