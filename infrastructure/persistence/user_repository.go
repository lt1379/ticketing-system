package persistence

import (
	"context"
	"database/sql"

	"github.com/lts1379/ticketing-system/domain/model"
	"github.com/lts1379/ticketing-system/domain/repository"
	"github.com/lts1379/ticketing-system/infrastructure/logger"
)

type UserRepository struct {
	sqlDB *sql.DB
}

func NewUserRepository(sqlDB *sql.DB) repository.IUser {
	return &UserRepository{sqlDB}
}

func (userRepository *UserRepository) GetById(ctx context.Context, id int) (model.User, error) {
	var user model.User

	statement, err := userRepository.sqlDB.PrepareContext(ctx, `SELECT u.id, u.name, u.user_name, u.password, u.created_at, u.updated_at 
	FROM public.user AS u 
	WHERE u.id = $1`)

	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while prepare statement")
		return user, err
	}
	defer statement.Close()

	result := statement.QueryRow(id)
	err = result.Scan(&user.ID, &user.Name, &user.UserName, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while query")
		return user, err
	}

	return user, nil
}

func (userRepository *UserRepository) GetByUserName(ctx context.Context, userName string) (model.User, error) {
	var user model.User

	statement, err := userRepository.sqlDB.PrepareContext(ctx, `SELECT u.id, u.name, u.user_name, u.password, u.created_at, u.updated_at 
	FROM public.user AS u 
	WHERE u.user_name = $1`)

	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while prepare statement")
		return user, err
	}
	defer statement.Close()

	result := statement.QueryRow(userName)
	err = result.Scan(&user.ID, &user.Name, &user.UserName, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while query")
		return user, err
	}

	return user, nil
}

func (userRepository *UserRepository) CreateUser(ctx context.Context, user model.User) error {
	statement, err := userRepository.sqlDB.PrepareContext(ctx, `INSERT INTO public.user (name, user_name, password) VALUES ($1, $2, $3)`)

	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while prepare statement")
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(user.Name, user.UserName, user.Password)
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error execute query")
		return err
	}

	return nil
}
