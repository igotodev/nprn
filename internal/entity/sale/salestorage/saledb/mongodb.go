package saledb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (s *SaleDB) Create(ctx context.Context, sale salemodel.Sale) (string, error) {
	result, err := s.collection.InsertOne(ctx, sale)
	if err != nil {
		return "", fmt.Errorf("failed to create new sale: %v", err)
	}

	objID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert objectID to Hex[%s]", objID.Hex())
	}
	s.logger.Tracef(fmt.Sprintf("sale id=%s is created", objID.Hex()))

	return objID.Hex(), nil
}

func (s *SaleDB) GetOne(ctx context.Context, id string) (salemodel.Sale, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return salemodel.Sale{}, fmt.Errorf("failed to convert objectID to Hex[%s]", objID.Hex())
	}

	filter := bson.M{"_id": objID}

	result := s.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return salemodel.Sale{}, fmt.Errorf("failed to find sale with id=%s", id)
	}

	var sale salemodel.Sale

	err = result.Decode(&sale)
	if err != nil {
		return salemodel.Sale{}, fmt.Errorf("failed to decode sale: %v", err)
	}

	return sale, nil
}

func (s *SaleDB) GetAll(ctx context.Context) ([]salemodel.Sale, error) {

	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get all sales: %v", err)
	}

	var sales []salemodel.Sale

	err = cursor.All(ctx, &sales)
	if err != nil {
		return nil, fmt.Errorf("failed to decode all sales: %v", err)
	}

	return sales, nil
}

func (s *SaleDB) Update(ctx context.Context, sale salemodel.Sale) error {
	objID, err := primitive.ObjectIDFromHex(sale.ID)
	if err != nil {
		return fmt.Errorf("failed to convert sale id=%v to objectID: %v", sale.ID, err)
	}

	filter := bson.M{"_id": objID}

	saleBytes, err := bson.Marshal(sale)
	if err != nil {
		return fmt.Errorf("failed to marhsal sale to bytes: %v", err)
	}

	var updateSaleObj bson.M

	err = bson.Unmarshal(saleBytes, &updateSaleObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal sale bytes: %v", err)
	}

	delete(updateSaleObj, "_id") // for not to overwrite id

	update := bson.M{"$set": updateSaleObj}

	result, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update sale: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("sale is not found: %v", err)
	}

	s.logger.Tracef("matched %d documents and modified %d documents", result.MatchedCount, result.ModifiedCount)

	return nil
}

func (s *SaleDB) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert objectID to Hex[%s]", objID.Hex())
	}

	filter := bson.M{"_id": objID}

	result, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute delete sale: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("sale is not found")
	}

	s.logger.Tracef("deleted %d documents", result.DeletedCount)

	return nil
}
