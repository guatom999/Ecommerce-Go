package usersRepositories

import (
	"fmt"

	"github.com/guatom999/Ecommerce-Go/modules/users"
	"github.com/guatom999/Ecommerce-Go/modules/users/usersPattern"
	"github.com/jmoiron/sqlx"
)

type IUserRepository interface {
	InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error)
	FindOneUserByEmail(email string) (*users.UserCredentialCheck, error)
}

type userRepository struct {
	db *sqlx.DB
}

func UsersRepository(db *sqlx.DB) IUserRepository {
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

func (r *userRepository) FindOneUserByEmail(email string) (*users.UserCredentialCheck, error) {
	query := `
		SELECT 
			"id" , 
			"email" , 
			"password" , 
			"username" ,
			"role_id" 
		FROM "users" 
		WHERE "email" = $1;
	`

	user := new(users.UserCredentialCheck)

	if err := r.db.Get(user, query, email); err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}
