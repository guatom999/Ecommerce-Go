package middlewaresUsecases

import "github.com/guatom999/Ecommerce-Go/modules/middlewares/middlewaresRepositories"

type IMiddlewaresUsecase interface {
	FindAccessToken(userId string, accessToken string) bool
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
