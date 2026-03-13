package usersUsecases

import (
	"errors"
	"testing"

	"github.com/guatom999/Ecommerce-Go/config"
	"github.com/guatom999/Ecommerce-Go/modules/users"
	"github.com/guatom999/Ecommerce-Go/modules/users/usersRepositories"
	"github.com/guatom999/Ecommerce-Go/pkg/authen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// --- Mock Config ---

type mockJwtConfig struct{}

func (m *mockJwtConfig) SecretKey() []byte      { return []byte("test-secret-key-32-bytes-long!!") }
func (m *mockJwtConfig) AdminKey() []byte       { return []byte("test-admin-key-32-bytes-long!!!") }
func (m *mockJwtConfig) ApiKey() []byte         { return []byte("test-api-key-32-bytes-long!!!!!") }
func (m *mockJwtConfig) AccessExpiresAt() int   { return 3600 }
func (m *mockJwtConfig) RefreshExpireAt() int   { return 86400 }
func (m *mockJwtConfig) SetAccessExpires(t int) {}
func (m *mockJwtConfig) SetRefreshExpire(t int) {}

type mockConfig struct{}

func (m *mockConfig) App() config.IAppConfig { return nil }
func (m *mockConfig) Db() config.IDbConfig   { return nil }
func (m *mockConfig) Jwt() config.IJwtConfig { return &mockJwtConfig{} }

// --- Helper ---

func newUserUsecaseWithMock() (IUserUsecase, *usersRepositories.MockUserRepository) {
	mockRepo := new(usersRepositories.MockUserRepository)
	uc := UsersUsecase(&mockConfig{}, mockRepo)
	return uc, mockRepo
}

// --- InsertCustomer ---
// BcryptHashing แก้ไข req.Password ก่อนส่งเข้า repo → ใช้ mock.Anything สำหรับ req

func TestInsertCustomer_InsertUserFailed(t *testing.T) {
	uc, mockRepo := newUserUsecaseWithMock()

	req := &users.UserRegisterReq{Email: "test@test.com", Password: "password123", Username: "testuser"}
	mockRepo.On("InsertUser", mock.Anything, false).Return((*users.UserPassport)(nil), errors.New("insert failed"))

	result, err := uc.InsertCustomer(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestInsertCustomer_Success(t *testing.T) {
	uc, mockRepo := newUserUsecaseWithMock()

	req := &users.UserRegisterReq{Email: "test@test.com", Password: "password123", Username: "testuser"}
	expected := &users.UserPassport{User: &users.User{Id: "u1", Email: "test@test.com"}}
	mockRepo.On("InsertUser", mock.Anything, false).Return(expected, nil)

	result, err := uc.InsertCustomer(req)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

// --- InsertAdmin ---

func TestInsertAdmin_InsertUserFailed(t *testing.T) {
	uc, mockRepo := newUserUsecaseWithMock()

	req := &users.UserRegisterReq{Email: "admin@test.com", Password: "adminpass", Username: "admin"}
	mockRepo.On("InsertUser", mock.Anything, true).Return((*users.UserPassport)(nil), errors.New("insert failed"))

	result, err := uc.InsertAdmin(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestInsertAdmin_Success(t *testing.T) {
	uc, mockRepo := newUserUsecaseWithMock()

	req := &users.UserRegisterReq{Email: "admin@test.com", Password: "adminpass", Username: "admin"}
	expected := &users.UserPassport{User: &users.User{Id: "u2", Email: "admin@test.com", RoleId: 2}}
	mockRepo.On("InsertUser", mock.Anything, true).Return(expected, nil)

	result, err := uc.InsertAdmin(req)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

// --- GetPassport ---
// 3 จุด error: FindOneUserByEmail / bcrypt ไม่ตรง / InsertOauth

func TestGetPassport_FindUserFailed(t *testing.T) {
	uc, mockRepo := newUserUsecaseWithMock()

	req := &users.UserCredential{Email: "test@test.com", Password: "password123"}
	mockRepo.On("FindOneUserByEmail", "test@test.com").Return((*users.UserCredentialCheck)(nil), errors.New("user not found"))

	result, err := uc.GetPassport(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestGetPassport_PasswordIncorrect(t *testing.T) {
	uc, mockRepo := newUserUsecaseWithMock()

	// hash ของ "actualpassword" ไม่ตรงกับ "wrongpassword"
	hash, _ := bcrypt.GenerateFromPassword([]byte("actualpassword"), bcrypt.MinCost)
	req := &users.UserCredential{Email: "test@test.com", Password: "wrongpassword"}
	mockRepo.On("FindOneUserByEmail", "test@test.com").Return(
		&users.UserCredentialCheck{Id: "u1", Email: "test@test.com", Password: string(hash), RoleId: 1}, nil,
	)

	result, err := uc.GetPassport(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestGetPassport_InsertOauthFailed(t *testing.T) {
	uc, mockRepo := newUserUsecaseWithMock()

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	req := &users.UserCredential{Email: "test@test.com", Password: "password123"}
	mockRepo.On("FindOneUserByEmail", "test@test.com").Return(
		&users.UserCredentialCheck{Id: "u1", Email: "test@test.com", Password: string(hash), RoleId: 1}, nil,
	)
	mockRepo.On("InsertOauth", mock.Anything).Return(errors.New("oauth insert failed"))

	result, err := uc.GetPassport(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestGetPassport_Success(t *testing.T) {
	uc, mockRepo := newUserUsecaseWithMock()

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	req := &users.UserCredential{Email: "test@test.com", Password: "password123"}
	mockRepo.On("FindOneUserByEmail", "test@test.com").Return(
		&users.UserCredentialCheck{Id: "u1", Email: "test@test.com", Password: string(hash), Username: "testuser", RoleId: 1}, nil,
	)
	mockRepo.On("InsertOauth", mock.Anything).Return(nil)

	result, err := uc.GetPassport(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "u1", result.User.Id)
	assert.NotEmpty(t, result.Token.AccessToken)
	assert.NotEmpty(t, result.Token.RefreshToken)
	mockRepo.AssertExpectations(t)
}

// --- RefreshPassport ---
// 4 จุด error: ParseToken / FindOneOauth / GetProfile / UpdateOauth
// ใช้ mockJwtConfig เพื่อ generate token จริงใน test

func generateTestRefreshToken(claims *users.UserClaims) string {
	jwtCfg := &mockJwtConfig{}
	auth, _ := authen.NewAuth(authen.Refresh, jwtCfg, claims)
	return auth.SignToken()
}

func TestRefreshPassport_ParseTokenFailed(t *testing.T) {
	uc, _ := newUserUsecaseWithMock()

	req := &users.UserRefreshCredential{RefreshToken: "invalid.token.string"}

	result, err := uc.RefreshPassport(req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRefreshPassport_FindOneOauthFailed(t *testing.T) {
	uc, mockRepo := newUserUsecaseWithMock()

	token := generateTestRefreshToken(&users.UserClaims{Id: "u1", RoleId: 1})
	req := &users.UserRefreshCredential{RefreshToken: token}
	mockRepo.On("FindOneOauth", token).Return((*users.Oauth)(nil), errors.New("oauth not found"))

	result, err := uc.RefreshPassport(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestRefreshPassport_GetProfileFailed(t *testing.T) {
	uc, mockRepo := newUserUsecaseWithMock()

	token := generateTestRefreshToken(&users.UserClaims{Id: "u1", RoleId: 1})
	req := &users.UserRefreshCredential{RefreshToken: token}
	mockRepo.On("FindOneOauth", token).Return(&users.Oauth{Id: "oauth1", UserId: "u1"}, nil)
	mockRepo.On("GetProfile", "u1").Return((*users.User)(nil), errors.New("profile not found"))

	result, err := uc.RefreshPassport(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestRefreshPassport_UpdateOauthFailed(t *testing.T) {
	uc, mockRepo := newUserUsecaseWithMock()

	token := generateTestRefreshToken(&users.UserClaims{Id: "u1", RoleId: 1})
	req := &users.UserRefreshCredential{RefreshToken: token}
	mockRepo.On("FindOneOauth", token).Return(&users.Oauth{Id: "oauth1", UserId: "u1"}, nil)
	mockRepo.On("GetProfile", "u1").Return(&users.User{Id: "u1", RoleId: 1}, nil)
	mockRepo.On("UpdateOauth", mock.Anything).Return(errors.New("update failed"))

	result, err := uc.RefreshPassport(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestRefreshPassport_Success(t *testing.T) {
	uc, mockRepo := newUserUsecaseWithMock()

	token := generateTestRefreshToken(&users.UserClaims{Id: "u1", RoleId: 1})
	req := &users.UserRefreshCredential{RefreshToken: token}
	mockRepo.On("FindOneOauth", token).Return(&users.Oauth{Id: "oauth1", UserId: "u1"}, nil)
	mockRepo.On("GetProfile", "u1").Return(&users.User{Id: "u1", Email: "test@test.com", RoleId: 1}, nil)
	mockRepo.On("UpdateOauth", mock.Anything).Return(nil)

	result, err := uc.RefreshPassport(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "u1", result.User.Id)
	assert.NotEmpty(t, result.Token.AccessToken)
	mockRepo.AssertExpectations(t)
}

// --- DeleteOauth ---

func TestDeleteOauth_Failed(t *testing.T) {
	uc, mockRepo := newUserUsecaseWithMock()

	mockRepo.On("DeleteOauth", "oauth1").Return(errors.New("delete failed"))

	err := uc.DeleteOauth("oauth1")

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteOauth_Success(t *testing.T) {
	uc, mockRepo := newUserUsecaseWithMock()

	mockRepo.On("DeleteOauth", "oauth1").Return(nil)

	err := uc.DeleteOauth("oauth1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// --- GetUserProfile ---

func TestGetUserProfile_Failed(t *testing.T) {
	uc, mockRepo := newUserUsecaseWithMock()

	mockRepo.On("GetProfile", "u1").Return((*users.User)(nil), errors.New("not found"))

	result, err := uc.GetUserProfile("u1")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestGetUserProfile_Success(t *testing.T) {
	uc, mockRepo := newUserUsecaseWithMock()

	expected := &users.User{Id: "u1", Email: "test@test.com", Username: "testuser"}
	mockRepo.On("GetProfile", "u1").Return(expected, nil)

	result, err := uc.GetUserProfile("u1")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}
