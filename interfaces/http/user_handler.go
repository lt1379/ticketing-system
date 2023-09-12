package http

import (
	"crypto/md5"
	"fmt"
	"log"
	"my-project/domain/model"
	"my-project/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IUserHandler interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

type UserHandler struct {
	userUsecase usecase.IUserUsecase
}

func NewUserHandler(userUsecase usecase.IUserUsecase) IUserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func (userHandler *UserHandler) Login(c *gin.Context) {
	var req model.ReqLogin

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("An error occurred: %v", err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("An error occurred: %v", err.Error()))
	}

	res := userHandler.userUsecase.Login(c.Request.Context(), req)

	c.JSON(http.StatusOK, res)
}

func (userHandler *UserHandler) Register(c *gin.Context) {
	var req model.ReqRegister

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("An error occurred: %v", err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("An error occurred: %v", err.Error()))
		return
	}
	data := []byte(req.Password)
	req.Password = fmt.Sprintf("%x", md5.Sum(data))
	res := userHandler.userUsecase.Register(c.Request.Context(), req)

	c.JSON(http.StatusOK, res)
}
