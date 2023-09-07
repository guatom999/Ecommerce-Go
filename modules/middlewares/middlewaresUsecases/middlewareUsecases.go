package middlewaresUsecases

type IMiddlewaresUsecase interface {
}

type middlewaresUsecase struct {
	middleawresRepository middlewaresRepositories.IMiddlewaresRepository
}

func MiddlewaresUsecase(r middlewaresRepositories.IMiddlewaresRepository) IMiddlewaresUsecase {
	return &middlewaresUsecase{middleawresRepository: r}
}
