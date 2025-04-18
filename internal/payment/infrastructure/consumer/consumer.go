package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel"

	"github.com/mrluzy/gorder-v2/common/broker"
	"github.com/mrluzy/gorder-v2/common/genproto/orderpb"
	"github.com/mrluzy/gorder-v2/payment/app"
	"github.com/mrluzy/gorder-v2/payment/app/command"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	app app.Application
}

func NewConsumer(app app.Application) *Consumer {
	return &Consumer{
		app: app,
	}
}

func (c *Consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
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
			c.handleMessage(msg, q)
		}
	}()
	<-forever
}

func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue) {
	logrus.Infof("Payment recieves a message from %s, mag: %s", q.Name, string(msg.Body))

	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
	t := otel.Tracer("rabbitmq")
	_, span := t.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
	defer span.End()

	o := &orderpb.Order{}
	err := json.Unmarshal(msg.Body, o)
	if err != nil {
		logrus.Infof("Failed to unmarshal msg:%s to order, err:%v", string(msg.Body), err)
		_ = msg.Nack(false, false)
		return
	}

	if _, err := c.app.Commands.CreatePayment.Handle(ctx, command.CreatePayment{Order: o}); err != nil {
		logrus.Infof("Failed to create order, err:%v", err)
		_ = msg.Nack(false, false)
	}

	span.AddEvent("payment.created")

	_ = msg.Ack(false)
	logrus.Infof("consume succcessfully")

}
