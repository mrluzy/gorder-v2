package order

import (
	"fmt"
	"github.com/mrluzy/gorder-v2/order/entity"

	"github.com/pkg/errors"
	"github.com/stripe/stripe-go/v82"
)

// aggregate
type Order struct {
	ID          string
	CustomerID  string
	Status      string
	PaymentLink string
	Items       []*entity.Item
}

func NewOrder(ID string, customerID string, status string, paymentLink string, items []*entity.Item) (*Order, error) {
	if ID == "" {
		return nil, errors.New("empty ID")
	}
	if customerID == "" {
		return nil, errors.New("empty customerID")
	}
	if status == "" {
		return nil, errors.New("empty status")
	}
	if items == nil {
		return nil, errors.New("empty items")
	}
	return &Order{ID: ID, CustomerID: customerID, Status: status, PaymentLink: paymentLink, Items: items}, nil
}

func NewPendingOrder(customerID string, items []*entity.Item) (*Order, error) {
	if customerID == "" {
		return nil, errors.New("empty customerID")
	}

	if items == nil {
		return nil, errors.New("empty items")
	}
	return &Order{CustomerID: customerID, Status: "pending", Items: items}, nil
}
func (o *Order) IsPaid() error {
	if o.Status == string(stripe.CheckoutSessionPaymentStatusPaid) {
		return nil
	}
	return fmt.Errorf("order status not paid, order id = %s, status = %s", o.ID, o.Status)
}
