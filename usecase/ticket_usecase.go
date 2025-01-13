package usecase

import (
	"context"
	"my-project/domain/model"
	"my-project/domain/repository"
)

type ITicketUsecase interface {
	Create(ctx context.Context, ticket model.Ticket) (model.Ticket, error)
}

type TicketUsecase struct {
	ticketRepository repository.ITicketRepository
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
