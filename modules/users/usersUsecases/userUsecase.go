package usersUsecases

import (
	"fmt"

	"github.com/guatom999/Ecommerce-Go/config"
	"github.com/guatom999/Ecommerce-Go/modules/users"
	"github.com/guatom999/Ecommerce-Go/modules/users/usersRepositories"
	"github.com/guatom999/Ecommerce-Go/pkg/authen"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetPassport(req *users.UserCredential) (*users.UserPassport, error)
	RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error)
	DeleteOauth(oauthId string) error
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

func (u *userUsecase) GetPassport(req *users.UserCredential) (*users.UserPassport, error) {
	user, err := u.userRepository.FindOneUserByEmail(req.Email)

	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("password incorrect")
	}

	accessToken, err := authen.NewAuth(authen.Access, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})

	refreshToken, err := authen.NewAuth(authen.Refresh, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})

	// passport := users.UserPassport{
	// 	User: &users.User{

	// 	},

	// }

	passport := &users.UserPassport{
		User: &users.User{
			Id:       user.Id,
			Email:    user.Email,
			Username: user.Username,
			RoleId:   user.RoleId,
		},
		Token: &users.UserToken{
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}

	if err := u.userRepository.InsertOauth(passport); err != nil {
		return nil, fmt.Errorf("insert oauth failed : %v", err)
	}

	return passport, nil

}

func (u *userUsecase) RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error) {
	claims, err := authen.ParseToken(u.cfg.Jwt(), req.RefreshToken)

	if err != nil {
		return nil, err
	}

	oauth, err := u.userRepository.FindOneOauth(req.RefreshToken)

	if err != nil {
		return nil, err
	}

	profile, err := u.userRepository.GetProfile(oauth.UserId)

	if err != nil {
		return nil, fmt.Errorf("get user profile failed : %v", err)
	}

	newClaims := &users.UserClaims{
		Id:     profile.Id,
		RoleId: profile.RoleId,
	}

	accessToken, err := authen.NewAuth(
		authen.Access,
		u.cfg.Jwt(),
		newClaims,
	)

	if err != nil {
		return nil, err
	}

	refreshToken := authen.RepeatToken(
		u.cfg.Jwt(),
		newClaims,
		claims.ExpiresAt.Unix(),
	)

	passport := &users.UserPassport{
		User: profile,
		Token: &users.UserToken{
			Id:           oauth.Id,
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken,
		},
	}

	if err := u.userRepository.UpdateOauth(passport.Token); err != nil {
		return nil, err
	}

	return passport, nil

}

func (u *userUsecase) DeleteOauth(oauthId string) error {
	if err := u.userRepository.DeleteOauth(oauthId); err != nil {
		return err
	}

	return nil
}
