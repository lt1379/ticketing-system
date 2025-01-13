package repository

import (
	"context"
	"my-project/domain/dto"
	"my-project/domain/model"
)

type ITicketRepository interface {
	Create(ctx context.Context, ticket model.Ticket) (int64, error)
	GetAll(ctx context.Context, pagination dto.RequestPagination) ([]model.Ticket, int64, error)
}
