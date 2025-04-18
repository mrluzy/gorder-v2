package main

import (
	"context"
	"github.com/mrluzy/gorder-v2/common/tracing"

	"github.com/mrluzy/gorder-v2/common/broker"
	"github.com/mrluzy/gorder-v2/common/config"
	"github.com/mrluzy/gorder-v2/common/logging"
	"github.com/mrluzy/gorder-v2/common/server"
	"github.com/mrluzy/gorder-v2/payment/infrastructure/consumer"
	"github.com/mrluzy/gorder-v2/payment/service"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceName := viper.GetString("payment.service-name")
	serverType := viper.GetString("payment.server-to-run")
	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = shutdown(ctx)
	}()

	application, cleanup := service.NewApplication(ctx)
	defer cleanup()

	ch, closeCh := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)
	defer func() {
		_ = ch.Close()
		_ = closeCh()
	}()

	go consumer.NewConsumer(application).Listen(ch)

	paymentHandler := NewPaymentHandler(ch)
	switch serverType {
	case "http":
		server.RunHTTPServer(serviceName, paymentHandler.RegisterRoutes)
	case "grpc":
		logrus.Panic("unsupported server type: grpc")
	default:
		logrus.Panic("unreachable code")
	}
}
