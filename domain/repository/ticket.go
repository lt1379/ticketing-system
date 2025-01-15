package repository

import (
	"context"
	"github.com/lts1379/ticketing-system/domain/dto"
	"github.com/lts1379/ticketing-system/domain/model"
)

type ITicketRepository interface {
	Create(ctx context.Context, ticket model.Ticket) (int64, error)
	GetAll(ctx context.Context, pagination dto.RequestPagination) ([]model.Ticket, int64, error)
}
