package userstorage

import (
	"context"
	"nprn/internal/entity/user/usermodel"
)

type UserStorage interface {
	Create(ctx context.Context, user usermodel.UserInternal) (string, error)
	GetOne(ctx context.Context, username string, password string) (usermodel.UserTransfer, error)
	Update(ctx context.Context, user usermodel.UserInternal) error
	Delete(ctx context.Context, id string) error
}

type AuthStorage interface {
	Create(ctx context.Context, user usermodel.UserInternal) (string, error)
	GetOne(ctx context.Context, username string, password string) (string, error)
}
