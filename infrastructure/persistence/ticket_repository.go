package persistence

import (
	"context"
	"database/sql"
	"github.com/lts1379/ticketing-system/domain/dto"
	"github.com/lts1379/ticketing-system/domain/model"
	"github.com/lts1379/ticketing-system/domain/repository"
	"github.com/lts1379/ticketing-system/infrastructure/logger"
	"math"
	"strings"
	"sync"
	"time"
)

type TicketRepository struct {
	sqlDB *sql.DB
}

func (t *TicketRepository) GetAll(ctx context.Context, pagination dto.RequestPagination) ([]model.Ticket, int64, error) {
	query, params, err := generateQuery(pagination)
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while generate query")
		return nil, 0, err
	}
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
		err := rows.Scan(&ticket.Id, &ticket.Title, &ticket.Message, &ticket.UserId, &ticket.Status, &ticket.CreatedAt)
		if err != nil {
			logger.GetLogger().WithField("error", err).Error("Error while scan row")
		}
		ticket.Status = model.MapStatusTicket[model.Status(ticket.Status)]
		tickets = append(tickets, ticket)
	}
	err = rows.Err()
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while query rows")
	}

	return tickets, int64(len(tickets)), nil
}

func (t *TicketRepository) Create(ctx context.Context, ticket model.Ticket) (int64, error) {
	statement, err := t.sqlDB.PrepareContext(ctx, `INSERT INTO ticket (title, message, user_id, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`)

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

	res, err := statement.Exec(ticket.Title, ticket.Message, ticket.UserId, ticket.Status, ticket.CreatedAt, ticket.CreatedAt)
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error execute query")
		return 0, err
	}

	return res.LastInsertId()
}

func (t *TicketRepository) WorkerGetAll(ctx context.Context, pagination dto.RequestPagination) ([]model.Ticket, int64, error) {
	resultI := worker(pagination, func(ctx context.Context, pagination dto.RequestPagination) ([]model.Ticket, int64, error) {
		return t.GetAll(ctx, pagination)
	})
	result, ok := resultI.([]model.Ticket)
	if !ok {
		logger.GetLogger().WithField("error", "Error while casting result").Error("Error while casting result")
		return nil, 0, nil
	}
	return result, int64(len(result)), nil
}

func generateQuery(pagination dto.RequestPagination) (string, []interface{}, error) {
	var queryBuilder strings.Builder
	var params []interface{}
	queryBuilder.WriteString("SELECT ticket.Id, ticket.title, ticket.message, ticket.user_id, ticket.status, ticket.created_at FROM ticket")
	queryBuilder.WriteString(" ")
	//Where
	if pagination.Filter != nil {
		queryBuilder.WriteString(" WHERE ")
		filter := *pagination.Filter
		layout := "2006-01-02"
		if filter.Type == "before" {
			queryBuilder.WriteString("ticket.created_at < ?")
			parseTime, err := time.Parse(layout, filter.Value)
			if err != nil {
				logger.GetLogger().WithField("error", err).Error("error parsing ticket.created_at")
				return "", nil, err
			}
			params = append(params, parseTime)
		} else if filter.Type == "after" {
			queryBuilder.WriteString("ticket.created_at > ?")
			parseTime, err := time.Parse(layout, filter.Value)
			if err != nil {
				logger.GetLogger().WithField("error", err).Error("error parsing ticket.created_at")
				return "", nil, err
			}
			params = append(params, parseTime)
		} else {
			queryBuilder.WriteString("ticket.created_at BETWEEN ? AND ?")
			parseTime, err := time.Parse(layout, filter.Value)
			if err != nil {
				logger.GetLogger().WithField("error", err).Error("error parsing ticket.created_at")
				return "", nil, err
			}
			parseTime2, err := time.Parse(layout, filter.Value2)
			if err != nil {
				logger.GetLogger().WithField("error", err).Error("error parsing ticket.created_at")
				return "", nil, err
			}
			params = append(params, parseTime)
			params = append(params, parseTime2)
		}
	}
	queryBuilder.WriteString(" ")
	//Order By
	if pagination.Sort != nil {
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
	}
	queryBuilder.WriteString(" ")

	//LIMIT
	queryBuilder.WriteString("LIMIT ? OFFSET ?")
	params = append(params, pagination.PageSize)
	params = append(params, pagination.PageNumber)

	return queryBuilder.String(), params, nil
}

func worker(pagination dto.RequestPagination, operation func(ctx context.Context, pagination dto.RequestPagination) ([]model.Ticket, int64, error)) interface{} {
	wg := &sync.WaitGroup{}

	var queries []dto.RequestPagination

	perPage := 10

	pageSize := pagination.PageSize
	numberOfPage := math.Ceil(float64(pageSize) / float64(perPage))
	numberOfWorker := int(numberOfPage)
	errChan := make(chan error, numberOfWorker)

	for i := 0; i < numberOfWorker; i++ {
		query := pagination
		query.PageSize = perPage
		if pageSize%perPage != 0 && i == numberOfWorker-1 {
			query.PageSize = int(pageSize % perPage)
		}
		query.PageNumber = i * perPage
		queries = append(queries, query)
	}

	ticketResults := make([][]model.Ticket, numberOfWorker)

	for i := 0; i < numberOfWorker; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			query := queries[i]
			query.PageSize = perPage
			if pageSize%perPage != 0 && i == numberOfWorker-1 {
				query.PageSize = int(pageSize % perPage)
			}
			tickets, _, err := operation(context.Background(), query)
			if err != nil {
				logger.GetLogger().WithField("error", err).Error("Error while operation")
				errChan <- err
				return
			}
			ticketResults[i] = tickets
		}(i)
	}

	wg.Wait()

	close(errChan)

	var tickets []model.Ticket
	for _, ticketResult := range ticketResults {
		tickets = append(tickets, ticketResult...)
	}

	return tickets
}

func NewTicketRepository(db *sql.DB) repository.ITicketRepository {
	return &TicketRepository{sqlDB: db}
}
