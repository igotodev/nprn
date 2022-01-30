package main

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"nprn/internal/config"
	"nprn/internal/entity/user/userstorage/userdb"
	"nprn/internal/handler"
	"nprn/internal/service"
	"nprn/pkg/client/mongodb"
	"nprn/pkg/logging"
	"nprn/pkg/server"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("application is started")

	cfg := config.GetConfig()

	router := httprouter.New()
	myServer := server.NewServer()

	myMongo, err := mongodb.NewClient(context.Background(),
		cfg.MongoDB.Host, cfg.MongoDB.Port, cfg.MongoDB.Username,
		cfg.MongoDB.Password, cfg.MongoDB.DBName, cfg.MongoDB.AuthDB)
	if err != nil {
		logger.Fatal(err)
	}

	myUsers := userdb.NewCollection(myMongo, cfg.MongoDB.UserCollection, logger)
	//mySales := salesdb.NewCollection(myMongo, cfg.MongoDB.SaleCollection, logger)

	appService := service.NewService(myUsers, nil, logger)

	handl := handler.NewHandler(appService, logger)

	handl.RegisterRouting(router)

	logger.Fatal(myServer.Run(router, logger, cfg))
}
