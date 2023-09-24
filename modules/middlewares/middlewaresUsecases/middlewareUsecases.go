package middlewaresUsecases

import (
	"github.com/guatom999/Ecommerce-Go/modules/middlewares"
	"github.com/guatom999/Ecommerce-Go/modules/middlewares/middlewaresRepositories"
)

type IMiddlewaresUsecase interface {
	FindAccessToken(userId string, accessToken string) bool
	FindRole() ([]*middlewares.Role, error)
}

type middlewaresUsecase struct {
	middlewareRepository middlewaresRepositories.IMiddlewareRepository
}

func MiddlewaresUsecase(middlewareRepository middlewaresRepositories.IMiddlewareRepository) IMiddlewaresUsecase {
	return &middlewaresUsecase{middlewareRepository}
}

func (m *middlewaresUsecase) FindAccessToken(userId string, accessToken string) bool {

	return m.middlewareRepository.FindAccessToken(userId, accessToken)
}

func (m *middlewaresUsecase) FindRole() ([]*middlewares.Role, error) {

	roles, err := m.middlewareRepository.FindRole()
	if err != nil {
		return nil, err
	}

	return roles, nil
}
