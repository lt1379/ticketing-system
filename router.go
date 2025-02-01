package main

import (
	"github.com/lts1379/ticketing-system/domain/repository"
	httpHandler "github.com/lts1379/ticketing-system/interfaces/http"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitiateRouter(userHandler httpHandler.IUserHandler, testHandler httpHandler.ITestHandler, tickerHandler httpHandler.ITicketHandler, userRepository repository.IUserRepository) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://tulus.tech"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://tulus.tech"
		},
		MaxAge: 12 * time.Hour,
	}))

	api := router.Group("api")
	//api.Use(middleware.Auth(userRepository))

	router.POST("/login", userHandler.Login)
	router.POST("/register", userHandler.Register)

	api.POST("/tickets/create", tickerHandler.Create)
	api.POST("/tickets", tickerHandler.GetAll)

	router.POST("/healthz", testHandler.Test)

	api.POST("/", func(ctx *gin.Context) {
		res := ctx.Request.Body
		ctx.JSON(http.StatusOK, res)
	})

	return router
}
