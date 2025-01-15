package http

import (
	"github.com/lts1379/ticketing-system/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ITestHandler interface {
	Test(c *gin.Context)
}

type TestHandler struct {
	TestUsecase usecase.ITestUsecase
}

func NewTestHandler(testUsecase usecase.ITestUsecase) ITestHandler {
	return &TestHandler{TestUsecase: testUsecase}
}

func (testHandler *TestHandler) Test(c *gin.Context) {
	res := testHandler.TestUsecase.Test(c.Request.Context())
	c.JSON(http.StatusOK, res)
}
