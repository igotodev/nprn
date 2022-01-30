package salestorage

import (
	"context"
	"nprn/internal/entity/sale/salemodel"
)

type SaleStorage interface {
	Create(ctx context.Context, sales salemodel.Sale) (string, error)
	GetOne(ctx context.Context, username string, password string) (salemodel.Sale, error)
	GetAll(ctx context.Context) ([]salemodel.Sale, error)
	Update(ctx context.Context, user salemodel.Sale) error
	Delete(ctx context.Context, id string) error
}
