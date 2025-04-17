package command

import (
	"context"

	"github.com/mrluzy/gorder-v2/common/genproto/orderpb"
)

type OrderService interface {
	UpdateOrder(ctx context.Context, order *orderpb.Order) error
}
