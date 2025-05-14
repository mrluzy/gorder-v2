package broker

import (
	"context"
	"encoding/json"
	"github.com/mrluzy/gorder-v2/common/logging"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

const (
	EventOrderCreated = "order.created"
	EventOrderPaid    = "order.paid"
)

type RoutingType string

const (
	FanOut = "fan-out"
	Direct = "direct"
)

type PublishEventReq struct {
	Channel  *amqp.Channel
	Routing  RoutingType
	Queue    string
	Exchange string
	Body     any
}

func PublishEvent(ctx context.Context, req PublishEventReq) (err error) {
	_, dLog := logging.WhenEventPublish(ctx, req)
	defer dLog(nil, &err)

	if err = checkParam(req); err != nil {
		return err
	}

	switch req.Routing {
	default:
		logrus.WithContext(ctx).Panicf("unknown routing type: %v", req.Routing)
	case FanOut:
		return fanOut(ctx, req)
	case Direct:
		return directQueue(ctx, req)
	}
	return nil
}

func checkParam(req PublishEventReq) error {
	if req.Channel == nil {
		return errors.New("nil channel")
	}
	return nil
}

func directQueue(ctx context.Context, req PublishEventReq) (err error) {
	_, err = req.Channel.QueueDeclare(req.Queue, true, false, false, false, nil)
	if err != nil {
		return err
	}
	jsonBody, err := json.Marshal(req.Body)
	if err != nil {
		return err
	}
	return doPublish(ctx, req.Channel, req.Exchange, req.Queue, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         jsonBody,
		Headers:      InjectRabbitMQHeaders(ctx),
	})
}

func doPublish(ctx context.Context, ch *amqp.Channel, exchange, key string, mandatory bool, immediate bool, msg amqp.Publishing) error {
	// key就是队列的名字
	if err := ch.PublishWithContext(ctx, exchange, key, mandatory, immediate, msg); err != nil {
		logging.Warnf(ctx, nil, "_publishing_event_failed||exchange=%s||key=%s||msg=%v", exchange, key, msg)
		return errors.Wrap(err, "publish event error")
	}
	return nil
}

func fanOut(ctx context.Context, req PublishEventReq) (err error) {
	jsonBody, err := json.Marshal(req.Body)
	if err != nil {
		return err
	}
	return doPublish(ctx, req.Channel, req.Exchange, "", false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         jsonBody,
		Headers:      InjectRabbitMQHeaders(ctx),
	})
}
