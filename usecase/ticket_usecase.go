package usecase

import (
	"context"
	"github.com/lts1379/ticketing-system/domain/dto"
	"github.com/lts1379/ticketing-system/domain/model"
	"github.com/lts1379/ticketing-system/domain/repository"
	"github.com/lts1379/ticketing-system/infrastructure/logger"
	"time"
)

type ITicketUsecase interface {
	Create(ctx context.Context, ticket model.Ticket) (model.Ticket, error)
	GetAll(ctx context.Context, pagination dto.RequestPagination) ([]model.Ticket, error)
}

type TicketUsecase struct {
	ticketRepository repository.ITicketRepository
}

func (t TicketUsecase) GetAll(ctx context.Context, pagination dto.RequestPagination) ([]model.Ticket, error) {
	var res []model.Ticket
	var count int64
	var err error

	if pagination.PageSize >= 50 {
		res, count, err = t.ticketRepository.WorkerGetAll(ctx, pagination)
	} else {
		res, count, err = t.ticketRepository.GetAll(ctx, pagination)
	}
	logger.GetLogger().WithField("res", res).WithField("count", count).Info("ticketRepository.GetAll")
	if err != nil {
		logger.GetLogger().WithField("error", err.Error()).Error("ticketRepository.GetAll")
		return nil, err
	}
	tickets := []model.Ticket{}
	for _, ticket := range res {
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}

func (t TicketUsecase) Create(ctx context.Context, ticket model.Ticket) (model.Ticket, error) {
	// Load the location for the timezone
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		logger.GetLogger().WithField("error", err.Error()).Error("time.LoadLocation")
		return model.Ticket{}, err
	}
	now := time.Now().In(loc)

	ticket.Status = string(model.Open)
	ticket.CreatedAt = &now
	lastInsertedId, err := t.ticketRepository.Create(ctx, ticket)
	if err != nil {
		return model.Ticket{}, err
	}
	ticket.Status = model.MapStatusTicket[model.Status(ticket.Status)]
	ticket.Id = lastInsertedId
	return ticket, nil
}

func NewTicketUsecase(ticketRepository repository.ITicketRepository) ITicketUsecase {
	return &TicketUsecase{ticketRepository: ticketRepository}
}
