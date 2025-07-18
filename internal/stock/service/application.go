package service

import (
	"context"
	"github.com/mrluzy/gorder-v2/stock/infrastructure/integration"
	"github.com/mrluzy/gorder-v2/stock/infrastructure/persistent"

	"github.com/mrluzy/gorder-v2/common/metrics"
	"github.com/mrluzy/gorder-v2/stock/adapters"
	"github.com/mrluzy/gorder-v2/stock/app"
	"github.com/mrluzy/gorder-v2/stock/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(_ context.Context) app.Application {
	db := persistent.NewMySQL()
	stockRepo := adapters.NewMySQLStockRepository(db)
	stripeAPI := integration.NewStripeAPI()
	metricsClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, stripeAPI, logrus.StandardLogger(), metricsClient),
			GetItems:            query.NewGetItemsHandler(stockRepo, logrus.StandardLogger(), metricsClient),
		},
	}

}
