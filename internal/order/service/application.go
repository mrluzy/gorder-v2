package service

import (
	"context"
	"github.com/mrluzy/gorder-v2/order/adapters"
	"github.com/mrluzy/gorder-v2/order/app"
)

func NewApplication(ctx context.Context) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	return app.Application{}
}
