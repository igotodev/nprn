package main

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"nprn/internal/config"
	"nprn/internal/entity/sale/salestorage/saledb"
	"nprn/internal/entity/user/userstorage/userdb"
	"nprn/internal/handler"
	"nprn/internal/service"
	"nprn/pkg/client/mongodb"
	"nprn/pkg/logging"
	"nprn/pkg/server"
	"time"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("application is started")

	cfg := config.GetConfig()

	router := httprouter.New()
	myServer := server.NewServer()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	myMongo, err := mongodb.NewClient(ctx,
		cfg.MongoDB.Host, cfg.MongoDB.Port, cfg.MongoDB.Username,
		cfg.MongoDB.Password, cfg.MongoDB.DBName, cfg.MongoDB.AuthDB)
	if err != nil {
		logger.Fatal(err)
	}

	myUsers := userdb.NewCollection(myMongo, cfg.MongoDB.UserCollection, logger)
	mySales := saledb.NewCollection(myMongo, cfg.MongoDB.SaleCollection, logger)

	appService := service.NewService(myUsers, mySales, logger)

	handl := handler.NewHandler(appService, logger)

	handl.RegisterRouting(router)

	logger.Fatal(myServer.Run(router, logger, cfg))
}
