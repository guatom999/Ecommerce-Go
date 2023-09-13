package users

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

//entity

type User struct {
	Id       string `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Username string `db:"username" json:"username"`
	RoleId   int    `db:"role_id" json:"role_id"`
}

type UserRegisterReq struct {
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
	Username string `db:"username" json:"username"`
}

func (obj *UserRegisterReq) BcryptHashing() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(obj.Password), 10)
	if err != nil {
		return fmt.Errorf("hashed password failed: %v", err)
	}

	obj.Password = string(hashedPassword)
	return nil
}

func (obj *UserRegisterReq) Isbool() bool {

	match, err := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, obj.Email)

	if err != nil {
		return false
	}

	return match
}

type UserPassport struct {
	User  User       `json:"user"`
	Token *UserToken `json:"token"`
}

type UserToken struct {
	Id           string `db:"id" json:"id"`
	AccessToken  string `db:"access_token" json:"access_token"`
	RefreshToken string `db:"refresh_token" json:"refresh"`
}
