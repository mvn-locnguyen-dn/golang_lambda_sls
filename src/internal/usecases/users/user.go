package users

import (
	"golang_lambda_boilerplate/src/internal/forms"
	"golang_lambda_boilerplate/src/pkg/helpers"
	"golang_lambda_boilerplate/src/pkg/utils"
)

type IUser interface {
	Detail() (*forms.DetailUserResponse, error)
	List() ([]forms.DetailUserResponse, error)
}

type User struct {
	Utils utils.IUtil
}

func (us *User) Detail() (*forms.DetailUserResponse, error) {
	u := helpers.GenerateUsers(1)

	response := forms.DetailUserResponse{}
	if err := us.Utils.ConvertStruct(u[0], &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (us *User) List() ([]forms.DetailUserResponse, error) {
	u := helpers.GenerateUsers(10)

	response := []forms.DetailUserResponse{}
	if err := us.Utils.ConvertStruct(u, &response); err != nil {
		return nil, err
	}

	return response, nil
}
