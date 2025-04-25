package command

import (
	"context"
	"github.com/mrluzy/gorder-v2/common/logging"
	"github.com/pkg/errors"

	"fmt"
	"github.com/mrluzy/gorder-v2/common/convertor"
	"github.com/mrluzy/gorder-v2/common/entity"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/status"

	"github.com/mrluzy/gorder-v2/common/broker"
	"github.com/mrluzy/gorder-v2/common/decorator"
	"github.com/mrluzy/gorder-v2/order/app/query"
	domain "github.com/mrluzy/gorder-v2/order/domain/order"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type CreateOrder struct {
	CustomerID string
	Items      []*entity.ItemWithQuantity
}

type CreateOrderResult struct {
	OrderID string
}

type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]

type createOrderHandler struct {
	orderRepo domain.Repository
	stockGRPC query.StockService
	channel   *amqp.Channel
}

func NewCreateOrderHandler(orderRepo domain.Repository, stockGRPC query.StockService, logger *logrus.Logger, channel *amqp.Channel, metricClient decorator.MetricsClient) CreateOrderHandler {
	if orderRepo == nil {
		panic("orderRepo is nil")
	}
	if stockGRPC == nil {
		panic("stockGRPC is nil")
	}
	if channel == nil {
		panic("channel is nil")
	}
	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
		createOrderHandler{
			orderRepo: orderRepo,
			stockGRPC: stockGRPC,
			channel:   channel,
		},
		logger,
		metricClient,
	)
}

// 该处理器的核心逻辑包括验证库存、创建订单并将其发布到 RabbitMQ 队列
func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
	var err error
	defer logging.WhenCommandExecute(ctx, "CreateOrderHandler", cmd, err)

	//

	// 创建 OpenTelemetry 的 Span 用于追踪
	t := otel.Tracer("rabbitmq")
	ctx, span := t.Start(ctx, fmt.Sprintf("rabbitmq.%s.publish", broker.EventOrderCreated))
	defer span.End()

	// 验证库存
	validItems, err := c.validate(ctx, cmd.Items)
	if err != nil {
		return nil, err
	}

	// 创建待处理订单
	pendingOrder, err := domain.NewPendingOrder(cmd.CustomerID, validItems)
	if err != nil {
		return nil, err
	}

	// 存储订单
	o, err := c.orderRepo.Create(ctx, pendingOrder)
	if err != nil {
		return nil, err
	}

	// 声明 RabbitMQ 队列, 并将订单信息发布到 RabbitMQ
	err = broker.PublishEvent(ctx, broker.PublishEventReq{
		Channel:  c.channel,
		Routing:  broker.Direct,
		Queue:    broker.EventOrderCreated,
		Exchange: "",
		Body:     o,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "publish order error, queue name: %s", broker.EventOrderCreated)
	}

	return &CreateOrderResult{OrderID: o.ID}, nil
}

func (c createOrderHandler) validate(ctx context.Context, items []*entity.ItemWithQuantity) ([]*entity.Item, error) {
	if len(items) == 0 {
		return nil, errors.New("no items, must have at least one item")
	}
	items = packItems(items)
	res, err := c.stockGRPC.CheckIfItemsInStock(ctx, convertor.NewItemWithQuantityConvertor().EntitiesToProtos(items))
	if err != nil {
		return nil, status.Convert(err).Err()
	}
	return convertor.NewItemConvertor().ProtosToEntities(res.Items), nil
}

func packItems(items []*entity.ItemWithQuantity) []*entity.ItemWithQuantity {
	mergedItems := make(map[string]int32)
	for _, item := range items {
		mergedItems[item.ID] += item.Quantity
	}
	var res []*entity.ItemWithQuantity
	for id, quantity := range mergedItems {
		res = append(res, entity.NewItemWithQuantity(id, quantity))
	}
	return res
}
