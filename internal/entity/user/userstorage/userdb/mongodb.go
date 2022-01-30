package userdb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"nprn/internal/entity/user/usermodel"
	"nprn/pkg/logging"
)

type UserDB struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func NewCollection(database *mongo.Database, collection string, logger *logging.Logger) *UserDB {
	return &UserDB{
		collection: database.Collection(collection),
		logger:     logger,
	}
}

func (u *UserDB) Create(ctx context.Context, user usermodel.UserInternal) (string, error) {
	result, err := u.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %v", err)
	}

	u.logger.Trace(fmt.Sprintf("user <%s> is created", user.Username))

	objID, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return objID.Hex(), nil
	}

	return "", fmt.Errorf("failed to convert objectID[%v] to Hex", objID)
}

func (u *UserDB) GetOne(ctx context.Context, username string, password string) (usermodel.UserTransfer, error) {
	//objID, err := primitive.ObjectIDFromHex(id)
	//if err != nil {
	//	return usermodel.UserTransfer{}, fmt.Errorf("failed to convert Hex[%v] to objectID", objID)
	//}

	filter := bson.M{"username": username, "password": password}

	result := u.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return usermodel.UserTransfer{}, fmt.Errorf("failed to find user with username[%s] and/or password[%s]", username, password)
	}

	var user usermodel.UserTransfer

	err := result.Decode(&user)
	if err != nil {
		return usermodel.UserTransfer{}, fmt.Errorf("failed to decode user: %v", err)
	}

	return user, nil
}

func (u *UserDB) Update(ctx context.Context, user usermodel.UserInternal) error {

	objID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return fmt.Errorf("failed to convert user id[%v] to objectID: %v", user.ID, err)
	}

	filter := bson.M{"_id": objID}

	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marhsal user to bytes: %v", err)
	}

	var updateUserObj bson.M

	err = bson.Unmarshal(userBytes, &updateUserObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user bytes: %v", err)
	}

	delete(updateUserObj, "_id") // for not to overwrite id

	update := bson.M{"$set": updateUserObj}

	result, err := u.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update user: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user is not found: %v", err)
	}

	u.logger.Tracef("matched %d documents and modified %d documents", result.MatchedCount, result.ModifiedCount)

	return nil
}

func (u *UserDB) Delete(ctx context.Context, id string) error {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert user id to objectID: %v", err)
	}

	filter := bson.M{"_id": objID}

	result, err := u.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute delete user: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("user is not found")
	}

	u.logger.Tracef("deleted %d documents", result.DeletedCount)

	return nil
}
