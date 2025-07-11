package ports

import (
	"context"
	"github.com/mrluzy/gorder-v2/common/convertor"

	"github.com/mrluzy/gorder-v2/common/genproto/orderpb"
	"github.com/mrluzy/gorder-v2/order/app"
	"github.com/mrluzy/gorder-v2/order/app/command"
	"github.com/mrluzy/gorder-v2/order/app/query"
	domain "github.com/mrluzy/gorder-v2/order/domain/order"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
	_, err := G.app.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
		CustomerID: request.CustomerID,
		Items:      convertor.NewItemWithQuantityConvertor().ProtosToEntities(request.Items),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
	o, err := G.app.Queries.GetCustomerOrder.Handle(ctx, query.GetCustomerOrder{
		OrderID:    request.OrderID,
		CustomerID: request.CustomerID,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &orderpb.Order{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		Items:       convertor.NewItemConvertor().EntitiesToProtos(o.Items),
		PaymentLink: o.PaymentLink,
	}, nil
}

func (G GRPCServer) UpdateOrder(ctx context.Context, request *orderpb.Order) (_ *emptypb.Empty, err error) {
	order, err := domain.NewOrder(
		request.ID,
		request.CustomerID,
		request.Status,
		request.PaymentLink,
		convertor.NewItemConvertor().ProtosToEntities(request.Items),
	)
	if err != nil {
		err = status.Error(codes.Internal, err.Error())
		return nil, err
	}
	_, err = G.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
		Order: order,
		UpdateFn: func(ctx context.Context, oldOrder *domain.Order) (*domain.Order, error) {
			if err := oldOrder.UpdateStatus(request.Status); err != nil {
				return nil, err
			}
			if err := oldOrder.UpdatePaymentLink(request.PaymentLink); err != nil {
				return nil, err
			}
			if err := oldOrder.UpdateItems(convertor.NewItemConvertor().ProtosToEntities(request.Items)); err != nil {
				return nil, err
			}
			return oldOrder, nil
		},
	})
	return nil, err
}
