package usecase_test

import (
	"context"
	"crypto/md5"
	"database/sql"
	"fmt"
	"github.com/lts1379/ticketing-system/domain/model"
	"github.com/lts1379/ticketing-system/mocks/repomocks"
	"github.com/lts1379/ticketing-system/usecase"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserUsecase_RegisterSuccess(t *testing.T) {
	userRepository := &repomocks.IUserRepository{}
	userRepository.On("CreateUser", context.Background(), mock.AnythingOfType("model.User")).Return(nil).Once()

	userUsecase := usecase.NewUserUsecase(userRepository)
	response := userUsecase.Register(context.Background(), model.ReqRegister{
		Name:     "Lambok Tulus Simamora",
		UserName: "lamboktulus1379",
		Password: "MyPassword_123",
	})

	assert.NotNil(t, response)
	assert.Equal(t, response.ResponseCode, "200")
}

func TestUserUsecase_RegisterError(t *testing.T) {
	userRepository := &repomocks.IUserRepository{}
	userRepository.On("CreateUser", context.Background(), mock.AnythingOfType("model.User")).Return(sql.ErrNoRows).Once()

	userUsecase := usecase.NewUserUsecase(userRepository)
	response := userUsecase.Register(context.Background(), model.ReqRegister{
		Name:     "Lambok Tulus Simamora",
		UserName: "lamboktulus1379",
		Password: "MyPassword_123",
	})

	assert.NotNil(t, response)
	assert.Equal(t, response.ResponseCode, "500")
}

func TestUserUsecase_LoginSuccess(t *testing.T) {
	userRepository := &repomocks.IUserRepository{}
	md5Req := fmt.Sprintf("%x", md5.Sum([]byte("MyPassword_123")))
	userRepository.On("GetByUserName", context.Background(), mock.Anything).Return(model.User{
		ID:        1,
		Name:      "Lambok Tulus Simamora",
		UserName:  "lamboktulus1379",
		Password:  md5Req,
		CreatedAt: time.Now(),
		CreatedBy: 0,
		UpdatedAt: time.Now(),
		UpdatedBy: 0,
	}, nil).Once()

	userUsecase := usecase.NewUserUsecase(userRepository)

	response := userUsecase.Login(context.Background(), model.ReqLogin{
		UserName: "lamboktulus1379",
		Password: "MyPassword_123",
	})

	assert.NotNil(t, response)
	assert.Equal(t, "200", response.ResponseCode)
}

func TestUserUsecase_LoginUserNotFound(t *testing.T) {
	userRepository := &repomocks.IUserRepository{}
	userRepository.On("GetByUserName", context.Background(), mock.Anything).Return(model.User{}, sql.ErrNoRows).Once()

	userUsecase := usecase.NewUserUsecase(userRepository)

	md5Req := fmt.Sprintf("%x", md5.Sum([]byte("MyPassword_123")))
	response := userUsecase.Login(context.Background(), model.ReqLogin{
		UserName: "lamboktulus1379",
		Password: md5Req,
	})

	assert.NotNil(t, response)
	assert.Equal(t, "401", response.ResponseCode)
}

func TestUserUsecase_LoginUserWrongPassword(t *testing.T) {
	userRepository := &repomocks.IUserRepository{}
	md5Req := fmt.Sprintf("%x", md5.Sum([]byte("MyPassword_123")))
	userRepository.On("GetByUserName", context.Background(), mock.Anything).Return(model.User{
		ID:        1,
		Name:      "Lambok Tulus Simamora",
		UserName:  "lamboktulus1379",
		Password:  md5Req,
		CreatedAt: time.Now(),
		CreatedBy: 0,
		UpdatedAt: time.Now(),
		UpdatedBy: 0,
	}, nil).Once()

	userUsecase := usecase.NewUserUsecase(userRepository)

	response := userUsecase.Login(context.Background(), model.ReqLogin{
		UserName: "lamboktulus1379",
		Password: "MyPassword_124",
	})

	assert.NotNil(t, response)
	assert.Equal(t, "401", response.ResponseCode)
}
