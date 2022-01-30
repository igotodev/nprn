package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, host, port, username, password, database, authDB string) (*mongo.Database, error) {
	var mdbURL string
	var isAuth bool

	if username == "" && password == "" {
		mdbURL = fmt.Sprintf("mongodb://%s:%s", host, port)
		isAuth = false
	} else {
		mdbURL = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
		isAuth = true
	}

	clientOptions := options.Client().ApplyURI(mdbURL)
	if isAuth {
		if authDB == "" {
			authDB = database
		}

		credential := options.Credential{
			AuthSource: authDB,
			Username:   username,
			Password:   password,
		}

		clientOptions.SetAuth(credential)
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping mongoDB: %v", err)
	}

	return client.Database(database), nil
}
