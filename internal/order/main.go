package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mrluzy/gorder-v2/common/broker"
	_ "github.com/mrluzy/gorder-v2/common/config"
	"github.com/mrluzy/gorder-v2/common/discovery"
	"github.com/mrluzy/gorder-v2/common/genproto/orderpb"
	"github.com/mrluzy/gorder-v2/common/logging"
	"github.com/mrluzy/gorder-v2/common/server"
	"github.com/mrluzy/gorder-v2/common/tracing"
	"github.com/mrluzy/gorder-v2/order/infrastructure/consumer"
	"github.com/mrluzy/gorder-v2/order/ports"
	"github.com/mrluzy/gorder-v2/order/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	logging.Init()
}

func main() {
	serviceName := viper.GetString("order.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		// 在 main 函数退出前调用 shutdown(ctx)，用于优雅关闭 Jaeger tracer
		_ = shutdown(ctx)
	}()

	application, cleanup := service.NewApplication(ctx)
	defer cleanup()

	deregisterFn, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = deregisterFn()
	}()

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

	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		svc := ports.NewGRPCServer(application)
		orderpb.RegisterOrderServiceServer(server, svc)
	})

	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		router.StaticFile("/success", "../../public/success.html")
		ports.RegisterHandlersWithOptions(router, HTTPServer{app: application}, ports.GinServerOptions{
			BaseURL:      "/api",
			Middlewares:  nil,
			ErrorHandler: nil,
		})
	})
}
