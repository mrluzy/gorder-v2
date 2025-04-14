package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mrluzy/gorder-v2/order/app"
	"github.com/mrluzy/gorder-v2/order/app/query"
	"net/http"
)

type HttpServer struct {
	app app.Application
}

func (h HttpServer) PostCustomerCustomerIDOrder(c *gin.Context, customerID string) {
	//TODO implement me
	panic("implement me")
}

func (h HttpServer) GetCustomerCustomerIDOrderOrderID(c *gin.Context, customerID string, orderID string) {
	o, err := h.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{
		CustomerID: "fake-customer-id",
		OrderID:    "fake-id",
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": o})
}
