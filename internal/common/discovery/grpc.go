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
	registry, err := consul.New(viper.GetString("consul.addr"))
	if err != nil {
		return func() error { return nil }, err
	}
	instanceID := GenerateInstanceID(serviceName)
	grpcAddr := viper.Sub(serviceName).GetString("grpc-addr")
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		return func() error { return nil }, err
	}
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

// GetServiceAddr 从 Consul 注册中心 获取一个服务地址
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
