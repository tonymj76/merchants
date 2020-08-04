package main

import (
	"github.com/Investliftng/ocm-api/merchant/api/handler"
	"github.com/Investliftng/ocm-api/merchant/proto/merchant"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/web"
	"github.com/sirupsen/logrus"
)

func main() {
	srv := web.NewService(
		web.Name("api.merchant"),
		web.Version("0.1"),
	)
	srv.Init()
	handler.ClientMerchant = merchant.NewMerchantService("service.merchant", client.DefaultClient)
	srv.Handle("/", handler.Router())

	if err := srv.Run(); err != nil {
		logrus.WithError(err).Fatal("unable to run service")
	}
}
