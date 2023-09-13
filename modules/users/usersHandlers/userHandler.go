package usersHandlers

import (
	"github.com/guatom999/Ecommerce-Go/config"
	"github.com/guatom999/Ecommerce-Go/modules/users/usersUsecases"
)

type IusersHandler interface {
}

type usersHandler struct {
	cfg         config.IConfig
	userUsecase usersUsecases.IUserUsecase
}

func UsersHandler(cfg config.IConfig, userUsecase usersUsecases.IUserUsecase) IusersHandler {
	return &usersHandler{cfg, userUsecase}
}
