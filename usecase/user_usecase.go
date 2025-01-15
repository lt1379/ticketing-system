package usecase

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/lts1379/ticketing-system/domain/dto"
	"github.com/lts1379/ticketing-system/domain/model"
	"github.com/lts1379/ticketing-system/domain/repository"
	"github.com/lts1379/ticketing-system/infrastructure/configuration"
	"github.com/lts1379/ticketing-system/infrastructure/logger"
	"github.com/lts1379/ticketing-system/infrastructure/utils"
	"time"
)

type IUserUsecase interface {
	Login(ctx context.Context, req model.ReqLogin) dto.ResLogin
	Register(ctx context.Context, req model.ReqRegister) dto.ResRegister
}

type UserUsecase struct {
	userRepository repository.IUser
}

func NewUserUsecase(userRepository repository.IUser) IUserUsecase {
	return &UserUsecase{userRepository: userRepository}
}

func (userUsecase *UserUsecase) Login(ctx context.Context, req model.ReqLogin) dto.ResLogin {
	var res dto.ResLogin

	user, err := userUsecase.userRepository.GetByUserName(ctx, req.UserName)
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while Getting username")
		res.ResponseCode = "401"
		res.ResponseMessage = "Unautorized."
		return res
	}
	md5Req := fmt.Sprintf("%x", md5.Sum([]byte(req.Password)))

	if md5Req != user.Password {
		logger.GetLogger().WithField("request_password", md5Req).Error("Password not matching")
		res.ResponseCode = "401"
		res.ResponseMessage = "Unautorized."
		return res
	}

	secretKey := configuration.C.App.SecretKey

	// Create the Claims
	expiration := time.Now().Add(5 * time.Minute)

	claims := make(map[string]interface{})
	claims["user_name"] = user.UserName
	claims["exp"] = expiration.Unix()
	claims["is"] = fmt.Sprint(user.ID)

	accessToken, err := utils.GenerateToken(claims, secretKey)
	if err != nil {
		logger.GetLogger().WithField("error", err).Info("Error while Signed string")
		res.ResponseCode = "401"
		res.ResponseMessage = "Unautorized"
		return res
	}
	res.ResponseCode = "200"
	res.ResponseMessage = "Success"
	res.Data.AccessToken = accessToken
	res.Data.ExpiresAt = expiration.Unix()

	return res
}

func (userUcase *UserUsecase) Register(ctx context.Context, req model.ReqRegister) dto.ResRegister {
	var res dto.ResRegister

	reqUser := model.User{
		Name:     req.Name,
		UserName: req.UserName,
		Password: req.Password,
	}
	err := userUcase.userRepository.CreateUser(ctx, reqUser)
	if err != nil {
		res.Data = nil
		res.ResponseCode = "500"
		res.ResponseMessage = "Internal server error"
		return res
	}
	userDto := dto.UserDto{
		Name:     req.Name,
		UserName: req.UserName,
	}
	res.Data = userDto
	res.ResponseCode = "200"
	res.ResponseMessage = "Success"

	return res
}
