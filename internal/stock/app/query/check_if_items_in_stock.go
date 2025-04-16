package query

import (
	"context"
	"github.com/mrluzy/gorder-v2/common/decorator"
	"github.com/mrluzy/gorder-v2/common/genproto/orderpb"
	domain "github.com/mrluzy/gorder-v2/stock/domain/stock"
	"github.com/sirupsen/logrus"
)

type CheckIfItemsInStock struct {
	Items []*orderpb.ItemWithQuantity
}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*orderpb.Item]

type checkIfItemsInStockHandler struct {
	stockRepo domain.Repository
}

func NewCheckIfItemsInStockHandler(
	stockRepo domain.Repository,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) CheckIfItemsInStockHandler {
	if stockRepo == nil {
		panic("nil stockRepo")
	}
	return decorator.ApplyQueryDecorators[CheckIfItemsInStock, []*orderpb.Item](
		checkIfItemsInStockHandler{stockRepo: stockRepo},
		logger,
		metricClient,
	)
}

// todo:删除
var stub = map[string]string{
	"1": "price_1RE3xECQIkU5HEs5mCwUNKQ5",
	"2": "price_1RE3wRCQIkU5HEs5mtvBQE7U",
}

func (c checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*orderpb.Item, error) {
	var items []*orderpb.Item
	for _, item := range query.Items {
		// TODO：改成从数据库或stripe获取
		priceID, ok := stub[item.ID]
		if !ok {
			priceID = stub["1"]
		}
		items = append(items, &orderpb.Item{
			ID:       item.ID,
			Quantity: item.Quantity,
			PriceID:  priceID,
		})
	}
	return items, nil
}
