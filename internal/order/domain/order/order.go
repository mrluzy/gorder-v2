package order

import (
	"fmt"
	"github.com/mrluzy/gorder-v2/common/consts"
	"github.com/mrluzy/gorder-v2/common/entity"
	"slices"

	"github.com/pkg/errors"
)

// aggregate
type Order struct {
	ID          string
	CustomerID  string
	Status      string
	PaymentLink string
	Items       []*entity.Item
}

func (o *Order) isValidStatus(to string) bool {
	switch o.Status {
	case consts.OrderStatusPending:
		return slices.Contains([]string{consts.OrderStatusWaitingForPayment}, to)
	case consts.OrderStatusWaitingForPayment:
		return slices.Contains([]string{consts.OrderStatusPaid}, to)
	case consts.OrderStatusPaid:
		return slices.Contains([]string{consts.OrderStatusReady}, to)
	default:
		return false
	}
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
	return &Order{CustomerID: customerID, Status: consts.OrderStatusPending, Items: items}, nil
}

func (o *Order) UpdateStatus(to string) error {
	if !o.isValidStatus(to) {
		return fmt.Errorf("cannot transmit to invalid status: %s", to)
	}
	o.Status = to
	return nil
}

func (o *Order) UpdatePaymentLink(link string) error {
	//if link == "" {
	//	return errors.New("empty payment link")
	//}
	o.PaymentLink = link
	return nil
}

func (o *Order) UpdateItems(res []*entity.Item) error {
	o.Items = res
	return nil
}
