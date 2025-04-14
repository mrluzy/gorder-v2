package ports

import (
	"context"
	"github.com/mrluzy/gorder-v2/common/genproto/orderpb"
	"github.com/mrluzy/gorder-v2/common/genproto/stockpb"
	"github.com/mrluzy/gorder-v2/stock/app"
	"github.com/sirupsen/logrus"
)

type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
	logrus.Infof("rpc_request_in, stock.GetItems")
	defer func() {
		logrus.Infof("rpc_request_out, stock.GetItems")
	}()
	fakeItems := []*orderpb.Item{
		{
			ID: "stockGRPC GetItems id",
		},
	}
	return &stockpb.GetItemsResponse{Items: fakeItems}, nil
}

func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
	logrus.Infof("rpc_request_in, stock.CheckIfItemsInStock")
	defer func() {
		logrus.Infof("rpc_request_out, stock.CheckIfItemsInStock")
	}()
	return nil, nil
}
