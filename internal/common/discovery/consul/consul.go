package consul

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
)

type Registry struct {
	client *api.Client
}

var (
	consulClient *Registry
	once         sync.Once
	initErr      error
)

func New(consulAddr string) (*Registry, error) {

	once.Do(func() {
		config := api.DefaultConfig()
		config.Address = consulAddr
		client, err := api.NewClient(config)
		if err != nil {
			initErr = err
			return
		}
		consulClient = &Registry{client: client}
	})
	if initErr != nil {
		return nil, initErr
	}
	return consulClient, nil
}

func (r Registry) Register(_ context.Context, instanceID, serviceName, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("invalid host:port format")
	}
	host := parts[0]
	port, _ := strconv.Atoi(parts[1])
	return r.client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      instanceID,
		Name:    serviceName,
		Address: host,
		Port:    port,
		Check: &api.AgentServiceCheck{
			CheckID:                        instanceID,
			TLSSkipVerify:                  false,
			TTL:                            "5s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	})
}

func (r Registry) Deregister(_ context.Context, instanceID, serviceName string) error {
	logrus.WithFields(logrus.Fields{
		"instanceID":  instanceID,
		"serviceName": serviceName,
	}).Info("deregister instance from consul")
	return r.client.Agent().ServiceDeregister(instanceID)
}

// 从 Consul 注册中心 获取所有服务的可用实例地址
func (r Registry) Discover(_ context.Context, serviceName string) ([]string, error) {
	// 查询服务健康检查通过的实例
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}
	var ids []string
	// 遍历服务健康检查通过的实例，提取每个entry的address和port
	for _, entry := range entries {
		ids = append(ids, fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port))
	}
	return ids, nil
}

func (r Registry) HealthCheck(instanceID, _ string) error {
	return r.client.Agent().UpdateTTL(instanceID, "online", api.HealthPassing)
}
