package usersUsecases

import (
	"github.com/guatom999/Ecommerce-Go/config"
	"github.com/guatom999/Ecommerce-Go/modules/users"
	"github.com/guatom999/Ecommerce-Go/modules/users/usersRepositories"
)

type IUserUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
}

type userUsecase struct {
	cfg            config.IConfig
	userRepository usersRepositories.IUserRepository
}

func UsersUsecase(cfg config.IConfig, userRepository usersRepositories.IUserRepository) IUserUsecase {
	return &userUsecase{
		cfg,
		userRepository,
	}
}

func (u *userUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {

	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	result, err := u.userRepository.InsertUser(req, false)
	if err != nil {
		return nil, err
	}

	return result, nil
}
