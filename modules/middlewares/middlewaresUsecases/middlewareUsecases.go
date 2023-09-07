package middlewaresUsecases

import "github.com/guatom999/Ecommerce-Go/modules/middlewares/middlewaresRepositories"

type IMiddlewaresUsecase interface {
}

type middlewaresUsecase struct {
	middlewareRepository middlewaresRepositories.IMiddlewareRepository
}

func MiddlewaresUsecase(r middlewaresRepositories.IMiddlewareRepository) IMiddlewaresUsecase {
	return &middlewaresUsecase{r}
}
