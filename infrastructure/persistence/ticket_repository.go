package persistence

import (
	"context"
	"database/sql"
	"my-project/domain/model"
	"my-project/domain/repository"
	"my-project/infrastructure/logger"
)

type TicketRepository struct {
	sqlDB *sql.DB
}

func (t *TicketRepository) Create(ctx context.Context, ticket model.Ticket) (int64, error) {
	statement, err := t.sqlDB.PrepareContext(ctx, `INSERT INTO ticket (title, message, user_id) VALUES (?, ?, ?)`)

	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while prepare statement")
		return 0, err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			logger.GetLogger().WithField("error", err).Error("Error while close statement")
		}
	}(statement)

	res, err := statement.Exec(ticket.Title, ticket.Message, ticket.UserId)
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error execute query")
		return 0, err
	}

	return res.LastInsertId()
}

func NewTicketRepository(db *sql.DB) repository.ITicketRepository {
	return &TicketRepository{sqlDB: db}
}
