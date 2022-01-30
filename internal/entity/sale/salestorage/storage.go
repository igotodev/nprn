package salestorage

import (
	"context"
	"nprn/internal/entity/sale/salemodel"
)

type SaleStorage interface {
	Create(ctx context.Context, sale salemodel.Sale) (string, error)
	GetOne(ctx context.Context, id string) (salemodel.Sale, error)
	GetAll(ctx context.Context) ([]salemodel.Sale, error)
	Update(ctx context.Context, sale salemodel.Sale) error
	Delete(ctx context.Context, id string) error
}
