package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mrluzy/gorder-v2/common/broker"
	"github.com/mrluzy/gorder-v2/common/consts"
	"github.com/mrluzy/gorder-v2/common/convertor"
	"github.com/mrluzy/gorder-v2/common/entity"
	"github.com/mrluzy/gorder-v2/common/genproto/orderpb"
	"github.com/mrluzy/gorder-v2/common/logging"
	"github.com/pkg/errors"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

type OrderService interface {
	UpdateOrder(ctx context.Context, request *orderpb.Order) error
}

type Consumer struct {
	orderGRPC OrderService
}

func NewConsumer(orderGRPC OrderService) *Consumer {
	return &Consumer{orderGRPC: orderGRPC}
}

func (c *Consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare("", true, false, true, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	err = ch.QueueBind(q.Name, "", broker.EventOrderPaid, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Warnf("Failed to consume q:%s, err:%v", q.Name, err)
	}
	var forever chan struct{}
	go func() {
		for msg := range msgs {
			c.handleMessage(ch, msg, q)
		}
	}()
	<-forever
}

func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
	t := otel.Tracer("rabbitmq")
	ctx, span := t.Start(broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers), fmt.Sprintf("rabbitmq.%s.consume", q.Name))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			logging.Warnf(ctx, nil, "consume failed||from=%s||message=%s||error=%+v", q.Name, msg, err)
			_ = msg.Nack(false, false)
		} else {
			logging.Infof(ctx, nil, "%s", "consume success")
			_ = msg.Ack(false)
		}
	}()

	o := &entity.Order{}
	err = json.Unmarshal(msg.Body, o)
	if err != nil {
		err = errors.Wrap(err, "failed unmarshal msg.body to order")
		return
	}

	if o.Status == "paid" {
		err = errors.New("order is not paid, cannot cook")
	}
	cook(ctx, o)
	span.AddEvent(fmt.Sprintf("order_cook:%v", o))
	if err = c.orderGRPC.UpdateOrder(ctx, &orderpb.Order{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      consts.OrderStatusReady,
		Items:       convertor.NewItemConvertor().EntitiesToProtos(o.Items),
		PaymentLink: o.PaymentLink,
	}); err != nil {
		logging.Errorf(ctx, nil, "error updating order||rderID = %s||err = %v", o.ID, err)
		if err = broker.HandleRetry(ctx, ch, &msg); err != nil {
			err = errors.Wrapf(err, "retry_error||error handle retry, messageID = %s||err = %v", msg.MessageId, err)
		}
		return
	}

	span.AddEvent("kitchen.order.finished.updated")
	logrus.Infof("consume succcessfully")

}

func cook(ctx context.Context, o *entity.Order) {
	logrus.WithContext(ctx).Printf("cooking order:%s", o.ID)
	time.Sleep(time.Second * 5)
	logrus.WithContext(ctx).Printf("order:%s done!", o.ID)
}
