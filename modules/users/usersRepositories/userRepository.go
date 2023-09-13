package usersRepositories

import (
	"github.com/guatom999/Ecommerce-Go/modules/users"
	"github.com/guatom999/Ecommerce-Go/modules/users/usersPattern"
	"github.com/jmoiron/sqlx"
)

type IUserRepository interface {
	InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error)
}

type userRepository struct {
	db *sqlx.DB
}

func UserRepository(db *sqlx.DB) IUserRepository {
	return &userRepository{db}
}

func (r *userRepository) InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error) {

	result := usersPattern.InsertUser(r.db, req, isAdmin)

	var err error

	if isAdmin {
		result, err = result.Admin()
		if err != nil {
			return nil, err
		}

	} else {
		result, err = result.Customer()
		if err != nil {
			return nil, err
		}

	}

	user, err := result.Result()

	if err != nil {
		return nil, err
	}

	return user, nil
}
