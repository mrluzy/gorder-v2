package main

import (
	"fmt"
	"github.com/mrluzy/gorder-v2/common/tracing"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrluzy/gorder-v2/common/genproto/orderpb"
	"github.com/mrluzy/gorder-v2/order/app"
	"github.com/mrluzy/gorder-v2/order/app/command"
	"github.com/mrluzy/gorder-v2/order/app/query"
)

type HTTPServer struct {
	app app.Application
}

func (h HTTPServer) PostCustomerCustomerIDOrder(c *gin.Context, customerID string) {
	ctx, span := tracing.Start(c, "PostCustomerCustomerIDOrder")
	defer span.End()
	var req orderpb.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r, err := h.app.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
		CustomerID: req.CustomerID,
		Items:      req.Items,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":      "success",
		"customer_id":  req.CustomerID,
		"trace_id":     tracing.TraceID(ctx),
		"order_id":     r.OrderID,
		"redirect_url": fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerID, r.OrderID),
	})
}

func (h HTTPServer) GetCustomerCustomerIDOrderOrderID(c *gin.Context, customerID string, orderID string) {
	ctx, span := tracing.Start(c, "PostCustomerCustomerIDOrder")
	defer span.End()
	o, err := h.app.Queries.GetCustomerOrder.Handle(ctx, query.GetCustomerOrder{
		CustomerID: customerID,
		OrderID:    orderID,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  "success",
		"trace_id": tracing.TraceID(ctx),
		"data": gin.H{
			"Order": o,
		},
	})
}
