package decorator

import (
	"context"
	"fmt"
	"time"
)

type MetricsClient interface {
	Inc(key string, value int)
}

type queryMetricsDecorator[C, R any] struct {
	base   QueryHandler[C, R]
	client MetricsClient
}

func (q queryMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	start := time.Now()
	actionName := generateActionName(cmd)
	defer func() {
		end := time.Since(start)
		q.client.Inc(fmt.Sprintf("querys.%s.duration", actionName), int(end.Seconds()))
		if err == nil {
			q.client.Inc(fmt.Sprintf("querys.%s.success", actionName), 1)
		} else {
			q.client.Inc(fmt.Sprintf("querys.%s.failure", actionName), 1)
		}
	}()
	return q.base.Handle(ctx, cmd)
}

type commandMetricsDecorator[C, R any] struct {
	base   CommandHandler[C, R]
	client MetricsClient
}

func (q commandMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	start := time.Now()
	actionName := generateActionName(cmd)
	defer func() {
		end := time.Since(start)
		q.client.Inc(fmt.Sprintf("commands.%s.duration", actionName), int(end.Seconds()))
		if err == nil {
			q.client.Inc(fmt.Sprintf("commands.%s.success", actionName), 1)
		} else {
			q.client.Inc(fmt.Sprintf("commands.%s.failure", actionName), 1)
		}
	}()
	return q.base.Handle(ctx, cmd)
}
