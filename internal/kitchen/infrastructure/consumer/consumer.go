package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mrluzy/gorder-v2/common/broker"
	"github.com/mrluzy/gorder-v2/common/genproto/orderpb"
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

type Order struct {
	ID          string
	CustomerID  string
	Status      string
	PaymentLink string
	Items       []*orderpb.Item
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
	var err error
	logrus.Infof("Kicken recieves a message from %s, mag: %s", q.Name, string(msg.Body))

	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
	t := otel.Tracer("rabbitmq")
	mgCtx, span := t.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))

	defer func() {
		span.End()
		if err != nil {
			_ = msg.Nack(false, false)
		} else {
			_ = msg.Ack(false)
		}
	}()

	o := &Order{}
	err = json.Unmarshal(msg.Body, o)
	if err != nil {
		logrus.Infof("Failed to unmarshal msg:%s to order, err:%v", string(msg.Body), err)
		return
	}

	if o.Status == "paid" {
		err = errors.New("order is not paid, cannot cook")
	}
	cook(o)
	span.AddEvent(fmt.Sprintf("order_cook:%v", o))
	if err = c.orderGRPC.UpdateOrder(mgCtx, &orderpb.Order{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      "ready",
		Items:       o.Items,
		PaymentLink: o.PaymentLink,
	}); err != nil {
		if err = broker.HandleRetry(mgCtx, ch, &msg); err != nil {
			logrus.Warnf("kichen: failed to handle retry:%s, err:%v", string(msg.Body), err)
		}
		return
	}

	span.AddEvent("kitchen.order.finished.updated")
	logrus.Infof("consume succcessfully")

}

func cook(o *Order) {
	logrus.Printf("cooking order:%s", o.ID)
	time.Sleep(time.Second * 5)
	logrus.Printf("order:%s done!", o.ID)
}
