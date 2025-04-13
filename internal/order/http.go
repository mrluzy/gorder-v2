package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mrluzy/gorder-v2/order/app"
)

type HttpServer struct {
	app app.Application
}

func (h HttpServer) PostCustomerCustomerIDOrder(c *gin.Context, customerID string) {
	//TODO implement me
	panic("implement me")
}

func (h HttpServer) GetCustomerCustomerIDOrderOrderID(c *gin.Context, customerID string, orderID string) {
	//TODO implement me
	panic("implement me")
}
