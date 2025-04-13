package main

import (
	"github.com/mrluzy/gorder-v2/common/genproto/stockpb"
	"github.com/mrluzy/gorder-v2/common/server"
	"github.com/mrluzy/gorder-v2/stock/ports"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to_run")
	switch serverType {
	case "grpc":
		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
			svc := ports.NewGRPCServer()
			stockpb.RegisterStockServiceServer(server, svc)
		})
	case "http":
		// TODO
	default:
		panic("invalid stock server type")
	}

}
