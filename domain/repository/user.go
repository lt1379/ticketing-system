package repository

import (
	"context"

	"my-project/domain/model"
)

type IUser interface {
	GetById(ctx context.Context, id int) (model.User, error)
	GetByUserName(ctx context.Context, userName string) (model.User, error)
	CreateUser(ctx context.Context, user model.User) error
}
