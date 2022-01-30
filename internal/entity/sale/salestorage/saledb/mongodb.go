package saledb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"nprn/internal/entity/sale/salemodel"
	"nprn/pkg/logging"
)

type SaleDB struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func NewCollection(database *mongo.Database, collection string, logger *logging.Logger) *SaleDB {
	return &SaleDB{
		collection: database.Collection(collection),
		logger:     logger,
	}
}

func (s SaleDB) Create(ctx context.Context, sale salemodel.Sale) (string, error) {
	return "", nil
}

func (s SaleDB) GetOne(ctx context.Context, id string) (salemodel.Sale, error) {
	return salemodel.Sale{}, nil
}

func (s SaleDB) GetAll(ctx context.Context) ([]salemodel.Sale, error) {
	return nil, nil
}

func (s SaleDB) Update(ctx context.Context, sale salemodel.Sale) error {
	return nil
}

func (s SaleDB) Delete(ctx context.Context, id string) error {
	return nil
}
