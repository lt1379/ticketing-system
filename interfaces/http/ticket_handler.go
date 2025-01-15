package http

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lts1379/ticketing-system/domain/dto"
	"github.com/lts1379/ticketing-system/domain/model"
	"github.com/lts1379/ticketing-system/infrastructure/logger"
	"github.com/lts1379/ticketing-system/usecase"
	"math"
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
	var res dto.Res

	var req dto.RequestPagination
	if err := context.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().WithField("error", err.Error()).Error("failed to bind request")

		res.ResponseCode = "400"
		res.ResponseMessage = "Bad Request"
		res.Meta = err.Error()

		context.JSON(http.StatusBadRequest, res)
		return
	}

	req.PageSize = max(req.PageSize, 10)
	req.PageSize = min(req.PageSize, 50)

	if (req.PageSize % 10) < 6 {
		req.PageSize -= req.PageSize % 10
	} else {
		req.PageSize += 10 - (req.PageSize % 10)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		logger.GetLogger().WithField("error", err.Error()).Error("failed to validate request")

		res.ResponseCode = "400"
		res.ResponseMessage = "Bad Request"
		res.Meta = err.Error()

		context.JSON(http.StatusBadRequest, res)
		return
	}

	logger.GetLogger().WithField("req", req).Info("processing request")

	response, err := t.ticketUsecase.GetAll(context.Request.Context(), req)
	if err != nil {
		logger.GetLogger().WithField("error", err.Error()).Error("failed to get ticket")

		res.ResponseCode = "500"
		res.ResponseMessage = "Internal Server Error"
		res.Meta = err.Error()

		context.JSON(http.StatusInternalServerError, res)
		return
	}

	res.ResponseCode = "200"
	res.ResponseMessage = "Success"
	res.Data = response
	res.Meta = dto.Pagination{
		PageNumber:  1,
		PerPage:     req.PageSize,
		TotalPage:   int(math.Ceil(float64(len(response)) / float64(req.PageSize))),
		TotalRecord: len(response),
	}

	context.JSON(http.StatusOK, res)
}

type validationError struct {
	Namespace       string `json:"namespace"` // can differ when a custom TagNameFunc is registered or
	Field           string `json:"field"`     // by passing alt name to ReportError like below
	StructNamespace string `json:"structNamespace"`
	StructField     string `json:"structField"`
	Tag             string `json:"tag"`
	ActualTag       string `json:"actualTag"`
	Kind            string `json:"kind"`
	Type            string `json:"type"`
	Value           string `json:"value"`
	Param           string `json:"param"`
	Message         string `json:"message"`
}

func (t TicketHandler) Create(c *gin.Context) {
	var res dto.Res
	var req dto.RequestTicketDto

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().WithField("error", err).Error("failed to bind request")
		res.ResponseCode = "400"
		res.ResponseMessage = "Bad Request"
		res.Meta = err.Error()
		c.JSON(http.StatusBadRequest, res)
		return
	}

	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		var currentErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			e := validationError{
				Namespace:       err.Namespace(),
				Field:           err.Field(),
				StructNamespace: err.StructNamespace(),
				StructField:     err.StructField(),
				Tag:             err.Tag(),
				ActualTag:       err.ActualTag(),
				Kind:            fmt.Sprintf("%v", err.Kind()),
				Type:            fmt.Sprintf("%v", err.Type()),
				Value:           fmt.Sprintf("%v", err.Value()),
				Param:           err.Param(),
				Message:         err.Error(),
			}

			indent, err := json.MarshalIndent(e, "", "  ")
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			if e.Field == "Title" {
				e.Message = "Title must be at least 10 characters long"
			} else if e.Field == "Message" {
				e.Message = "Message must be at least 100 characters long"
			} else if e.Field == "UserId" {
				e.Message = "UserId is required"
			} else {
				e.Message = "Unknown error"
			}
			currentErrors = append(currentErrors, string(e.Message))

			fmt.Println(string(indent))
		}

		res.ResponseCode = "400"
		res.ResponseMessage = "Bad Request"
		res.Meta = currentErrors

		c.JSON(http.StatusBadRequest, res)
		return
	}

	ticket := model.Ticket{
		Title:   req.Title,
		Message: req.Message,
		UserId:  req.UserId,
	}

	response, err := t.ticketUsecase.Create(c.Request.Context(), ticket)
	if err != nil {
		logger.GetLogger().WithField("error", err.Error()).Error("ticket create")
		res.ResponseCode = "500"
		res.ResponseMessage = "Internal Server Error"
		res.Meta = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res.ResponseCode = "201"
	res.ResponseMessage = "Created"
	res.Data = response

	c.JSON(http.StatusCreated, res)
}

func NewTicketHandler(ticketUsecase usecase.ITicketUsecase) ITicketHandler {
	return &TicketHandler{ticketUsecase: ticketUsecase}
}
