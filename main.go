package main

import (
	"github.com/Investliftng/ocm-api/merchant/datastore"
	"github.com/Investliftng/ocm-api/merchant/handler"
	_ "github.com/Investliftng/ocm-api/merchant/log"
	micro "github.com/micro/go-micro/v2"
	log "github.com/sirupsen/logrus"

	pbM "github.com/Investliftng/ocm-api/merchant/proto/merchant"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("service.merchant"),
		micro.Version("0.1"),
	)

	// Initialise service
	service.Init()
	connect, err := datastore.NewConnection(log.New())
	if err != nil {
		log.WithError(err).Fatal("database failed to connect")
	}
	defer connect.Close()

	// Register Handler
	pbM.RegisterMerchantServiceHandler(service.Server(), &handler.Service{
		Repository:         connect,
		OutletRepository:   connect,
		TerminalRepository: connect,
	})

	// Run service
	if err := service.Run(); err != nil {
		log.WithError(err).Fatal("unable to run service")
	}
}
