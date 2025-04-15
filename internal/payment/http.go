package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PaymentHandler struct{}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{}
}

func (p *PaymentHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("./api/webhook", p.handleWebhook)
}

func (p *PaymentHandler) handleWebhook(context *gin.Context) {
	logrus.Infof("Got Webhook from stripe ")
}
