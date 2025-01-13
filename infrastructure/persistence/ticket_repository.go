package persistence

import (
	"context"
	"database/sql"
	"my-project/domain/dto"
	"my-project/domain/model"
	"my-project/domain/repository"
	"my-project/infrastructure/logger"
	"strings"
	"time"
)

type TicketRepository struct {
	sqlDB *sql.DB
}

func (t *TicketRepository) GetAll(ctx context.Context, pagination dto.RequestPagination) ([]model.Ticket, int64, error) {
	var queryBuilder strings.Builder
	var params []interface{}
	queryBuilder.WriteString("SELECT ticket.title, ticket.message, ticket.user_id, ticket.created_at FROM ticket")
	queryBuilder.WriteString(" ")
	//Where
	if pagination.Filter != nil {
		queryBuilder.WriteString(" WHERE ")
		filter := *pagination.Filter
		if filter.Type == "before" {
			queryBuilder.WriteString("ticket.created_at < ?")
			parseTime, err := time.Parse("2006-01-02", filter.Value)
			if err != nil {
				logger.GetLogger().WithField("error", err).Error("error parsing ticket.created_at")
				return nil, 0, err
			}
			params = append(params, parseTime)
		} else if filter.Type == "after" {
			queryBuilder.WriteString("ticket.created_at > ?")
			parseTime, err := time.Parse("2006-01-02", filter.Value)
			if err != nil {
				logger.GetLogger().WithField("error", err).Error("error parsing ticket.created_at")
				return nil, 0, err
			}
			params = append(params, parseTime)
		} else {
			queryBuilder.WriteString("ticket.created_at BETWEEN ? AND ?")
			parseTime, err := time.Parse("2006-01-02", filter.Value)
			if err != nil {
				logger.GetLogger().WithField("error", err).Error("error parsing ticket.created_at")
				return nil, 0, err
			}
			parseTime2, err := time.Parse("2006-01-02", filter.Value2)
			if err != nil {
				logger.GetLogger().WithField("error", err).Error("error parsing ticket.created_at")
				return nil, 0, err
			}
			params = append(params, parseTime)
			params = append(params, parseTime2)
		}
	}
	queryBuilder.WriteString(" ")
	//Order By
	if pagination.Sort.Name == "created_at" {
		queryBuilder.WriteString("ORDER BY created_at")
	} else {
		queryBuilder.WriteString("ORDER BY user_id")
	}
	queryBuilder.WriteString(" ")
	if pagination.Sort.Dir == "asc" {
		queryBuilder.WriteString("ASC")
	} else {
		queryBuilder.WriteString("DESC")
	}
	queryBuilder.WriteString(" ")

	//LIMIT
	queryBuilder.WriteString("LIMIT ? OFFSET ?")
	params = append(params, pagination.PageSize)
	params = append(params, 0)

	query := queryBuilder.String()
	logger.GetLogger().WithField("query", query).WithField("params", params).Info("query")
	statement, err := t.sqlDB.PrepareContext(ctx, query)
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while prepare statement")
		return nil, 0, err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			logger.GetLogger().WithField("error", err).Error("Error while close statement")
		}
	}(statement)

	rows, err := statement.QueryContext(ctx, params...)
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while query statement")
		return nil, 0, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.GetLogger().WithField("error", err).Error("Error while close rows")
		}
	}(rows)

	var tickets []model.Ticket
	for rows.Next() {
		ticket := model.Ticket{}
		err := rows.Scan(&ticket.Title, &ticket.Message, &ticket.UserId, &ticket.CreatedAt)
		if err != nil {
			logger.GetLogger().WithField("error", err).Error("Error while scan row")
		}
		tickets = append(tickets, ticket)
	}
	err = rows.Err()
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while query rows")
	}

	return tickets, int64(len(tickets)), nil
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
