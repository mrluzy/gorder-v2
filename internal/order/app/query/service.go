package query

import (
	"context"
	"github.com/mrluzy/gorder-v2/common/genproto/orderpb"
)

type StockService interface {
	CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) error
	GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error)
}
