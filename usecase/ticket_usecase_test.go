package usecase_test

import (
	"context"
	"github.com/lts1379/ticketing-system/domain/model"
	"github.com/lts1379/ticketing-system/mocks/repomocks"
	"github.com/lts1379/ticketing-system/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestTicketUsecase_Create(t *testing.T) {
	ticketRepository := &repomocks.ITicketRepository{}
	ticketRepository.On("Create", mock.Anything, mock.AnythingOfType("model.Ticket")).Return(int64(1), nil).Once()

	ticketUsecase := usecase.NewTicketUsecase(ticketRepository)
	response, err := ticketUsecase.Create(context.Background(), model.Ticket{
		Id:      1,
		Title:   "Title",
		Message: "Message",
		UserId:  1,
		Status:  string(model.Open),
	})

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, response.Id, int64(1))
}
