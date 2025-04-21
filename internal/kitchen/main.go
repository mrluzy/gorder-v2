package main

import (
	"context"
	"github.com/mrluzy/gorder-v2/common/broker"
	grpcClient "github.com/mrluzy/gorder-v2/common/client"
	"github.com/mrluzy/gorder-v2/common/tracing"
	"github.com/mrluzy/gorder-v2/kitchen/adapters"
	"github.com/mrluzy/gorder-v2/kitchen/infrastructure/consumer"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mrluzy/gorder-v2/common/config"
	"github.com/mrluzy/gorder-v2/common/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	logging.Init()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceName := viper.GetString("payment.service-name")

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = shutdown(ctx)
	}()

	orderClient, closeFn, err := grpcClient.NewOrderGRPCClient(ctx)
	if err != nil {
		logrus.Fatal(err)
	}
	defer closeFn()
	orderGRPC := adapters.NewOrderGRPC(orderClient)

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

	go consumer.NewConsumer(orderGRPC).Listen(ch)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		logrus.Info("Received shutdown signal, shutting down gracefully")
		os.Exit(0)
	}()
	logrus.Info("to exit, press Ctrl+C")
	select {}
}
