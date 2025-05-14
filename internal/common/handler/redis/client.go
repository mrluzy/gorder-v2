package redis

import (
	"context"
	"errors"
	"github.com/mrluzy/gorder-v2/common/logging"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"time"
)

// SetNX 执行 Redis 的 SETNX 操作。如果 key 不存在，则设置该 key 的值为 value，并为该 key 设置一个 TTL（过期时间）
func SetNX(ctx context.Context, client *redis.Client, key, value string, ttl time.Duration) (err error) {
	now := time.Now()
	// 打印日志
	defer func() {
		l := logrus.WithContext(ctx).WithFields(logrus.Fields{
			"start":       now,
			"key":         key,
			"value":       value,
			logging.Error: err,
			logging.Cost:  time.Since(now).Milliseconds(),
		})
		if err == nil {
			l.Info("_redis_setnx_success")
		} else {
			l.Warn("_redis_setnx_error")
		}
	}()

	if client == nil {
		return errors.New("redis client is nil")
	}
	_, err = client.SetNX(ctx, key, value, ttl).Result()
	return err
}

func Del(ctx context.Context, client *redis.Client, key string) (err error) {
	now := time.Now()
	defer func() {
		l := logrus.WithContext(ctx).WithFields(logrus.Fields{
			"start":       now,
			"key":         key,
			logging.Error: err,
			logging.Cost:  time.Since(now).Milliseconds(),
		})
		if err == nil {
			l.Info("_redis_del_success")
		} else {
			l.Warn("_redis_del_error")
		}
	}()
	if client == nil {
		return errors.New("redis client is nil")
	}
	_, err = client.Del(ctx, key).Result()
	return err
}
