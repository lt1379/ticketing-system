package http

import (
	"github.com/gin-gonic/gin"
	"my-project/domain/dto"
	"my-project/domain/model"
	"my-project/infrastructure/logger"
	"my-project/usecase"
	"net/http"
)

type ITicketHandler interface {
	Create(*gin.Context)
	GetAll(*gin.Context)
}

type TicketHandler struct {
	ticketUsecase usecase.ITicketUsecase
}

func (t TicketHandler) GetAll(context *gin.Context) {
	var req dto.RequestPagination
	if err := context.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().WithField("error", err.Error()).Error("failed to bind request")
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := t.ticketUsecase.GetAll(context.Request.Context(), req)
	if err != nil {
		logger.GetLogger().WithField("error", err.Error()).Error("failed to get ticket")
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, res)
}

func (t TicketHandler) Create(c *gin.Context) {
	var req dto.RequestTicketDto

	if err := c.ShouldBind(&req); err != nil {
		logger.GetLogger().WithField("error", err).Error("failed to bind request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket := model.Ticket{
		Title:   req.Title,
		Message: req.Message,
		UserId:  req.UserId,
	}

	res, err := t.ticketUsecase.Create(c.Request.Context(), ticket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logger.GetLogger().WithField("error", err.Error()).Error("ticket create")
		return
	}

	c.JSON(http.StatusCreated, res)
}

func NewTicketHandler(ticketUsecase usecase.ITicketUsecase) ITicketHandler {
	return &TicketHandler{ticketUsecase: ticketUsecase}
}
