package discovery

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/mrluzy/gorder-v2/common/discovery/consul"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func RegisterToConsul(ctx context.Context, serviceName string) (func() error, error) {
	// 初始化 Consul 注册客户端：
	// 通过 consul.New(viper.GetString("consul.addr")) 获取一个 Registry 实例，
	// 连接到 Consul 服务
	registry, err := consul.New(viper.GetString("consul.addr"))
	if err != nil {
		return func() error { return nil }, err
	}

	instanceID := GenerateInstanceID(serviceName)
	grpcAddr := viper.Sub(serviceName).GetString("grpc-addr")

	// 注册服务到 Consul：调用 registry.Register 完成服务注册，
	// 同时为服务设置了一个 5 秒的 TTL 健康检查
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		return func() error { return nil }, err
	}
	// 心跳检查：使用 goroutine 每秒向 Consul 发送一次心跳更新（HealthCheck），
	// 以保持服务处于健康状态。
	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				logrus.Panicf("no heartbeat from %s to registry, err=%v", serviceName, err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	logrus.WithFields(logrus.Fields{
		"serviceName": serviceName,
		"grpcAddr":    grpcAddr,
	}).Info("registered to consul")
	return func() error {
		return registry.Deregister(ctx, instanceID, serviceName)
	}, nil
}

// 从 Consul 注册中心 获取一个服务地址
func GetServiceAddr(ctx context.Context, serviceName string) (string, error) {
	registry, err := consul.New(viper.GetString("consul.addr"))
	if err != nil {
		return "", err
	}

	// Discover 从 Consul 注册中心 获取所有服务的可用实例地址
	addrs, err := registry.Discover(ctx, serviceName)
	if err != nil {
		return "", err
	}
	if len(addrs) == 0 {
		return "", fmt.Errorf("get empty %s address from consul", serviceName)
	}

	// 随机返回一个可用的addr
	i := rand.Intn(len(addrs))
	logrus.Infof("Discovered %d instance of %s address from consul, addrs:%v", len(addrs), serviceName, addrs)
	return addrs[i], nil
}
