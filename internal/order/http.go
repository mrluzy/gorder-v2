package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mrluzy/gorder-v2/common"
	client "github.com/mrluzy/gorder-v2/common/client/order"
	"github.com/mrluzy/gorder-v2/common/consts"
	"github.com/mrluzy/gorder-v2/common/convertor"
	"github.com/mrluzy/gorder-v2/common/handler/errors"
	"github.com/mrluzy/gorder-v2/order/app"
	"github.com/mrluzy/gorder-v2/order/app/command"
	"github.com/mrluzy/gorder-v2/order/app/dto"
	"github.com/mrluzy/gorder-v2/order/app/query"
)

type HTTPServer struct {
	common.BaseResponse
	app app.Application
}

func (h HTTPServer) PostCustomerCustomerIdOrder(c *gin.Context, _ string) {
	var (
		req  client.CreateOrderRequest
		resp dto.CreateOrderResponse
		err  error
	)

	defer func() {
		h.Response(c, err, &resp)
	}()

	if err = c.ShouldBindJSON(&req); err != nil {
		err = errors.NewWithErr(consts.ErrnoBindRequestError, err)
		return
	}
	if err = h.validate(req); err != nil {
		err = errors.NewWithErr(consts.ErrnoRequestValidateError, err)
		return
	}
	r, err := h.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CreateOrder{
		CustomerID: req.CustomerId,
		Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
	})
	if err != nil {
		return
	}

	resp = dto.CreateOrderResponse{
		CustomerID:  req.CustomerId,
		OrderID:     r.OrderID,
		RedirectURL: fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerId, r.OrderID),
	}

}

func (h HTTPServer) GetCustomerCustomerIdOrderOrderId(c *gin.Context, customerID string, orderID string) {
	var (
		resp struct {
			Order *client.Order `json:"order"`
		}
		err error
	)
	defer func() {
		h.Response(c, err, &resp)
	}()

	o, err := h.app.Queries.GetCustomerOrder.Handle(c.Request.Context(), query.GetCustomerOrder{
		CustomerID: customerID,
		OrderID:    orderID,
	})
	if err != nil {
		return
	}
	resp.Order = &client.Order{
		Id:          o.ID,
		CustomerId:  o.CustomerID,
		Status:      o.Status,
		Items:       convertor.NewItemConvertor().EntitiesToClients(o.Items),
		PaymentLink: o.PaymentLink,
	}

}

func (h HTTPServer) validate(req client.CreateOrderRequest) error {
	for _, v := range req.Items {
		if v.Quantity < 1 {
			return fmt.Errorf("quanti ty must be positive, got %d from %s", v.Quantity, v.Id)
		}
	}
	return nil
}
