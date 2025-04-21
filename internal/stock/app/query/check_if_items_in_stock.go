package query

import (
	"context"
	"github.com/mrluzy/gorder-v2/common/decorator"
	"github.com/mrluzy/gorder-v2/stock/entity"

	domain "github.com/mrluzy/gorder-v2/stock/domain/stock"
	"github.com/mrluzy/gorder-v2/stock/infrastructure/integration"
	"github.com/sirupsen/logrus"
)

type CheckIfItemsInStock struct {
	Items []*entity.ItemWithQuantity
}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*entity.Item]

type checkIfItemsInStockHandler struct {
	stockRepo domain.Repository
	stripeAPI *integration.StripeAPI
}

func NewCheckIfItemsInStockHandler(
	stockRepo domain.Repository,
	stripeAPI *integration.StripeAPI,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) CheckIfItemsInStockHandler {
	if stockRepo == nil {
		panic("nil stockRepo")
	}
	if stripeAPI == nil {
		panic("nil stripeAPI")
	}
	return decorator.ApplyQueryDecorators[CheckIfItemsInStock, []*entity.Item](
		checkIfItemsInStockHandler{
			stockRepo: stockRepo,
			stripeAPI: stripeAPI,
		},
		logger,
		metricClient,
	)
}

// todo:删除
var stub = map[string]string{
	"1": "price_1RE3xECQIkU5HEs5mCwUNKQ5",
	"2": "price_1RE3wRCQIkU5HEs5mtvBQE7U",
}

func (c checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*entity.Item, error) {
	var items []*entity.Item
	for _, item := range query.Items {
		// TODO：改成从数据库或stripe获取
		priceID, err := c.stripeAPI.GetPriceByProductID(ctx, item.ID)
		if err != nil || priceID == "" {
			return nil, err
		}
		items = append(items, &entity.Item{
			ID:       item.ID,
			Quantity: item.Quantity,
			PriceID:  priceID,
		})
	}
	return items, nil
}

func getStubPriceID(id string) string {
	priceID, ok := stub[id]
	if !ok {
		priceID = stub["1"]
	}
	return priceID
}
