package client

import (
	"context"
	"errors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"net"
	"time"

	"github.com/mrluzy/gorder-v2/common/discovery"
	"github.com/mrluzy/gorder-v2/common/genproto/orderpb"
	"github.com/mrluzy/gorder-v2/common/genproto/stockpb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewOrderGRPCClient(ctx context.Context) (client orderpb.OrderServiceClient, close func() error, err error) {
	if !WaitForOrderGRPCClient(viper.GetDuration("dial-grpc-timeout") * time.Second) {
		return nil, nil, errors.New("order grpc not available")
	}
	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("order.service-name"))
	if err != nil {
		return nil, func() error { return nil }, err
	}
	if grpcAddr == "" {
		logrus.Warnf("no order grpc address found")
	}

	opts := grpcDialOpts(grpcAddr)

	conn, err := grpc.NewClient(grpcAddr, opts...)
	if err != nil {
		return nil, func() error { return nil }, err
	}
	return orderpb.NewOrderServiceClient(conn), conn.Close, nil
}

func NewStockGRPCClient(ctx context.Context) (client stockpb.StockServiceClient, close func() error, err error) {
	// 检查stockGRPC服务是否可用
	if !WaitForStockGRPCClient(viper.GetDuration("dial-grpc-timeout") * time.Second) {
		return nil, nil, errors.New("stock grpc not available")
	}

	// GetServiceAddr 从 Consul 注册中心 获取一个服务地址
	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("stock.service-name"))
	if err != nil {
		return nil, func() error { return nil }, err
	}
	if grpcAddr == "" {
		logrus.Warnf("no stock grpc address found")
	}

	opts := grpcDialOpts(grpcAddr)

	conn, err := grpc.NewClient(grpcAddr, opts...)
	if err != nil {
		return nil, func() error { return nil }, err
	}
	return stockpb.NewStockServiceClient(conn), conn.Close, nil
}

// 为 gRPC 客户端构造连接选项
func grpcDialOpts(_ string) []grpc.DialOption {
	return []grpc.DialOption{
		// 默认情况下，gRPC 使用 TLS 加密传输
		// 这个配置明确表示使用不安全的明文连接，一般只在开发环境或测试阶段使用
		grpc.WithTransportCredentials(insecure.NewCredentials()),

		// 这个选项的目的是让你的 gRPC 客户端能自动记录 trace 信息（配合 Jaeger、Zipkin 等观察系统）：
		// 发送请求时创建 trace span
		// 和服务端 trace 信息链路打通
		// 自动传播上下文（context）
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}
}

func WaitForOrderGRPCClient(timeout time.Duration) bool {
	logrus.Infof("waiting for order grpc client, timeout: %v seconds", timeout.Seconds())
	return waitFor(viper.GetString("order.grpc-addr"), timeout)
}

func WaitForStockGRPCClient(timeout time.Duration) bool {
	logrus.Infof("waiting for stock grpc client, timeout: %v seconds", timeout.Seconds())
	return waitFor(viper.GetString("stock.grpc-addr"), timeout)
}

func waitFor(addr string, timeout time.Duration) bool {
	portAvailable := make(chan struct{})
	timeoutCh := time.After(timeout)

	go func() {
		for {
			select {
			case <-timeoutCh:
				return
			default:
				// continue
			}
			_, err := net.Dial("tcp", addr)
			if err == nil {
				close(portAvailable)
				return
			}
			time.Sleep(200 * time.Millisecond)
		}
	}()

	select {
	case <-portAvailable:
		return true
	case <-timeoutCh:
		return false
	}
}
