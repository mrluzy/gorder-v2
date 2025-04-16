package processor

import (
	"context"
	"github.com/mrluzy/gorder-v2/common/genproto/orderpb"
)

type InmemProcessor struct{}

func NewInmemProcessor() *InmemProcessor {
	return &InmemProcessor{}
	
}

func (i InmemProcessor) CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error) {
	return "inmem_payment_link", nil
}
