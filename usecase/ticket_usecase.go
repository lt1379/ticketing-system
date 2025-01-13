package usecase

import (
	"context"
	"my-project/domain/dto"
	"my-project/domain/model"
	"my-project/domain/repository"
	"my-project/infrastructure/logger"
)

type ITicketUsecase interface {
	Create(ctx context.Context, ticket model.Ticket) (model.Ticket, error)
	GetAll(ctx context.Context, pagination dto.RequestPagination) ([]model.Ticket, error)
}

type TicketUsecase struct {
	ticketRepository repository.ITicketRepository
}

func (t TicketUsecase) GetAll(ctx context.Context, pagination dto.RequestPagination) ([]model.Ticket, error) {
	res, _, err := t.ticketRepository.GetAll(ctx, pagination)
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
	lastInsertedId, err := t.ticketRepository.Create(ctx, ticket)
	if err != nil {
		return model.Ticket{}, err
	}
	ticket.Id = lastInsertedId
	return ticket, nil
}

func NewTicketUsecase(ticketRepository repository.ITicketRepository) ITicketUsecase {
	return &TicketUsecase{ticketRepository: ticketRepository}
}
