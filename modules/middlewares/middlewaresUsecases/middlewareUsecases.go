package middlewaresUsecases

import "github.com/guatom999/Ecommerce-Go/modules/middlewares/middlewaresRepositories"

type IMiddlewaresUsecase interface {
}

type middlewaresUsecase struct {
	middlewareRepository middlewaresRepositories.IMiddlewareRepository
}

func MiddlewaresUsecase(middlewareRepository middlewaresRepositories.IMiddlewareRepository) IMiddlewaresUsecase {
	return &middlewaresUsecase{middlewareRepository}
}
