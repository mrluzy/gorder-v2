package main

import (
	"github.com/mrluzy/gorder-v2/common/config"
	"github.com/mrluzy/gorder-v2/common/logging"
	"github.com/mrluzy/gorder-v2/common/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	logging.Init()
	if err := config.NewViperConfig(); err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	serverName := viper.GetString("payment.service-name")
	serverType := viper.GetString("payment.server-to-run")

	paymentHandler := NewPaymentHandler()

	switch serverType {
	case "http":
		server.RunHttpServer(serverName, paymentHandler.RegisterRoutes)
	case "grpc":
		logrus.Println("grpc: unsupported server type")
	default:
		logrus.Println("unsupported server type")
	}

}
