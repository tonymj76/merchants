package main

import (
	"context"

	"github.com/Investliftng/ocm-api/merchant/proto/merchant"
	"github.com/micro/go-micro/v2"
	"github.com/sirupsen/logrus"
)

func main() {
	service := micro.NewService(
		micro.Name("client.merchant"),
	)
	service.Init()
	client := merchant.NewMerchantService("service.merchant", service.Client())
	getRes, err := client.GetMerchantByID(context.Background(), &merchant.GetIdRequest{Id: 1})
	if err != nil {
		logrus.WithError(err).Fatal("error getting client")
	}
	logrus.WithField("response", getRes).Info("response values from client")

	if err := service.Run(); err != nil {
		logrus.WithError(err).Fatal("unable to run service")
	}
}
