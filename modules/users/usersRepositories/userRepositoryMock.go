package usersRepositories

import (
	"github.com/guatom999/Ecommerce-Go/modules/users"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error) {
	args := m.Called(req, isAdmin)
	return args.Get(0).(*users.UserPassport), args.Error(1)
}

func (m *MockUserRepository) FindOneUserByEmail(email string) (*users.UserCredentialCheck, error) {
	args := m.Called(email)
	return args.Get(0).(*users.UserCredentialCheck), args.Error(1)
}

func (m *MockUserRepository) InsertOauth(req *users.UserPassport) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockUserRepository) FindOneOauth(refreshToken string) (*users.Oauth, error) {
	args := m.Called(refreshToken)
	return args.Get(0).(*users.Oauth), args.Error(1)
}

func (m *MockUserRepository) UpdateOauth(req *users.UserToken) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockUserRepository) GetProfile(userId string) (*users.User, error) {
	args := m.Called(userId)
	return args.Get(0).(*users.User), args.Error(1)
}

func (m *MockUserRepository) DeleteOauth(userId string) error {
	args := m.Called(userId)
	return args.Error(0)
}
