package grpc

import (
	"context"
	"github.com/mrluzy/gorder-v2/common/genproto/orderpb"
	"github.com/mrluzy/gorder-v2/common/genproto/stockpb"
	"github.com/sirupsen/logrus"
)

type StockGRPC struct {
	client stockpb.StockServiceClient
}

func NewStockGRPC(client stockpb.StockServiceClient) *StockGRPC {
	return &StockGRPC{client: client}
}

func (s StockGRPC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error) {
	res, err := s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{Items: items})
	logrus.Infof("CheckIfItemsInStock response: %v", res)
	return res, err
}

func (s StockGRPC) GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error) {
	resp, err := s.client.GetItems(ctx, &stockpb.GetItemsRequest{ItemIDs: itemIDs})
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}
