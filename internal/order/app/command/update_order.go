package command

import (
	"context"
	"github.com/mrluzy/gorder-v2/common/logging"

	"github.com/mrluzy/gorder-v2/common/decorator"
	domain "github.com/mrluzy/gorder-v2/order/domain/order"
	"github.com/sirupsen/logrus"
)

type UpdateOrder struct {
	Order    *domain.Order
	UpdateFn func(context.Context, *domain.Order) (*domain.Order, error)
}

type UpdateOrderHandler decorator.CommandHandler[UpdateOrder, interface{}]

type updateOrderHandler struct {
	orderRepo domain.Repository
	//stockGRPC
}

func NewUpdateOrderHandler(orderRepo domain.Repository, logger *logrus.Logger, metricClient decorator.MetricsClient) UpdateOrderHandler {
	if orderRepo == nil {
		panic("nil orderRepo")
	}
	return decorator.ApplyCommandDecorators[UpdateOrder, interface{}](
		updateOrderHandler{orderRepo: orderRepo},
		logger,
		metricClient,
	)
}

func (c updateOrderHandler) Handle(ctx context.Context, cmd UpdateOrder) (interface{}, error) {
	var err error
	defer logging.WhenCommandExecute(ctx, "UpdateOrderHandler", cmd, err)

	if cmd.UpdateFn == nil {
		logrus.Panicf("UpdateFn is nil, cmd=%+v", cmd)
	}
	if err := c.orderRepo.Update(ctx, cmd.Order, cmd.UpdateFn); err != nil {
		return nil, err
	}
	return nil, nil
}
