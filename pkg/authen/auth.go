package authen

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/guatom999/Ecommerce-Go/config"
	"github.com/guatom999/Ecommerce-Go/modules/users"
)

type tokenType string

const (
	Access  tokenType = "access"
	Refresh tokenType = "refresh"
	Admin   tokenType = "admin"
	ApiKey  tokenType = "apikey"
)

type IAuth interface {
	SignToken() string
}

type IAdmin interface {
	SignToken() string
}

type IApiKey interface {
	SignToken() string
}

type auth struct {
	mapClaims *authMapClaims
	cfg       config.IJwtConfig
}

type admin struct {
	*auth
}

type apikey struct {
	*auth
}

type authMapClaims struct {
	Claims *users.UserClaims `json:"claims" `
	jwt.RegisteredClaims
}

func (a *auth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	signString, _ := token.SignedString(a.cfg.SecretKey())

	return signString
}

func (a *admin) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	signString, _ := token.SignedString(a.cfg.AdminKey())

	return signString
}

func (a *apikey) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	signString, _ := token.SignedString(a.cfg.ApiKey())

	return signString
}

func ParseToken(cfg config.IJwtConfig, tokenString string) (*authMapClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &authMapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("sign method not match algorithm")
		}

		return cfg.SecretKey(), nil

	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token is expired")
		} else {
			return nil, fmt.Errorf("parse token failed %v", err)
		}
	}

	if claims, ok := token.Claims.(*authMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is not authMapClaims")
	}

}

func ParseAdminToken(cfg config.IJwtConfig, tokenString string) (*authMapClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &authMapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("sign method not match algorithm")
		}

		return cfg.AdminKey(), nil

	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token is expired")
		} else {
			return nil, fmt.Errorf("parse token failed %v", err)
		}
	}

	if claims, ok := token.Claims.(*authMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is not authMapClaims")
	}

}

func ParseApikey(cfg config.IJwtConfig, tokenString string) (*authMapClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &authMapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("sign method not match algorithm")
		}

		return cfg.ApiKey(), nil

	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("key format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("key is expired")
		} else {
			return nil, fmt.Errorf("parse api key failed %v", err)
		}
	}

	if claims, ok := token.Claims.(*authMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is not apikeyMapClaims")
	}

}

func jwtTimeDurationCal(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func RepeatToken(cfg config.IJwtConfig, claims *users.UserClaims, exp int64) string {
	obj := &auth{
		mapClaims: &authMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "Ecommerce-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeRepeatAdapter(exp),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
		cfg: cfg,
	}

	return obj.SignToken()
}

func NewAuth(tokenType tokenType, cfg config.IJwtConfig, claims *users.UserClaims) (IAuth, error) {
	switch tokenType {
	case Access:
		return newAccessToken(cfg, claims), nil
	case Refresh:
		return newRefreshToken(cfg, claims), nil
	case Admin:
		return newAdminToken(cfg), nil
	case ApiKey:
		return newApiKey(cfg), nil
	default:
		return nil, fmt.Errorf("unknow accesstoken type")
	}
}

func newAccessToken(cfg config.IJwtConfig, claims *users.UserClaims) IAuth {
	return &auth{
		cfg: cfg,
		mapClaims: &authMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "Ecommerce-api",
				Subject:   "access-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.AccessExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newRefreshToken(cfg config.IJwtConfig, claims *users.UserClaims) IAuth {
	return &auth{
		mapClaims: &authMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "Ecommerce-api",
				Subject:   "access-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.RefreshExpireAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
		cfg: cfg,
	}
}

func newAdminToken(cfg config.IJwtConfig) IAuth {
	return &admin{
		&auth{
			mapClaims: &authMapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "Ecommerce-api",
					Subject:   "admin-token",
					Audience:  []string{"admin"},
					ExpiresAt: jwtTimeDurationCal(300),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
			cfg: cfg,
		},
	}
}

func newApiKey(cfg config.IJwtConfig) IApiKey {
	return &apikey{
		&auth{
			mapClaims: &authMapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "Ecommerce-api",
					Subject:   "api-key",
					Audience:  []string{"admin", "customer"},
					ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(2, 0, 0)),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
			cfg: cfg,
		},
	}
}
