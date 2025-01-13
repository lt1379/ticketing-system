package repository

import (
	"context"
	"my-project/domain/model"
)

type ITicketRepository interface {
	Create(ctx context.Context, ticket model.Ticket) (int64, error)
}
