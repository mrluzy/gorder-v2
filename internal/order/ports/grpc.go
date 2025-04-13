package ports

import (
	"context"
	"github.com/mrluzy/gorder-v2/common/genproto/orderpb"
	"github.com/mrluzy/gorder-v2/order/app"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (GRPCServer) UpdateOrder(ctx context.Context, order *orderpb.Order) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
