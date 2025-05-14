package service

import (
	"context"

	"github.com/mrluzy/gorder-v2/common/metrics"
	"github.com/spf13/viper"

	grpcClient "github.com/mrluzy/gorder-v2/common/client"
	"github.com/mrluzy/gorder-v2/payment/adapters"
	"github.com/mrluzy/gorder-v2/payment/app"
	"github.com/mrluzy/gorder-v2/payment/app/command"
	"github.com/mrluzy/gorder-v2/payment/domain"
	"github.com/mrluzy/gorder-v2/payment/infrastructure/processor"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	orderClient, closeOrderClient, err := grpcClient.NewOrderGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	orderGRPC := adapters.NewOrderGRPC(orderClient)
	stripeProcessor := processor.NewStripeProcessor(viper.GetString("stripe-key"))
	return newApplication(ctx, orderGRPC, stripeProcessor), func() {
		_ = closeOrderClient()
	}
}

func newApplication(_ context.Context, orderGRPC command.OrderService, processor domain.Processor) app.Application {
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreatePayment: command.NewCreatePaymentHandler(processor, orderGRPC, logrus.StandardLogger(), metricClient),
		},
	}
}
