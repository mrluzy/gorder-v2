package adapters

import (
	"context"
	"github.com/mrluzy/gorder-v2/common/logging"
	"strconv"
	"sync"
	"time"

	domain "github.com/mrluzy/gorder-v2/order/domain/order"
)

type MemoryOrderRepository struct {
	lock  *sync.RWMutex
	store []*domain.Order
}

func NewMemoryOrderRepository() *MemoryOrderRepository {
	s := make([]*domain.Order, 0)
	s = []*domain.Order{
		{
			ID:          "fake-id",
			CustomerID:  "fake-customer-id",
			Status:      "fake-order-status",
			PaymentLink: "fake-payment-link",
			Items:       nil,
		},
	}
	return &MemoryOrderRepository{
		lock:  &sync.RWMutex{},
		store: s,
	}
}

func (m *MemoryOrderRepository) Create(ctx context.Context, order *domain.Order) (created *domain.Order, err error) {
	_, dLog := logging.WhenRequest(ctx, "MemoryOrderRepository.Create", map[string]interface{}{"order": order})
	defer dLog(created, &err)

	m.lock.Lock()
	defer m.lock.Unlock()
	newOrder := &domain.Order{
		ID:          strconv.FormatInt(time.Now().Unix(), 10),
		CustomerID:  order.CustomerID,
		Status:      order.Status,
		PaymentLink: order.PaymentLink,
		Items:       order.Items,
	}
	return newOrder, nil
}

func (m *MemoryOrderRepository) Get(ctx context.Context, id, customerID string) (got *domain.Order, err error) {
	_, dLog := logging.WhenRequest(ctx, "MemoryOrderRepository.Get", map[string]interface{}{
		"id":         id,
		"customerID": customerID})
	defer dLog(got, &err)

	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, order := range m.store {
		if order.ID == id && order.CustomerID == customerID {
			return order, nil
		}
	}
	return nil, domain.NotFoundError{OrderID: id}
}

func (m *MemoryOrderRepository) Update(ctx context.Context, order *domain.Order, UpdateFn func(context.Context, *domain.Order) (*domain.Order, error)) (err error) {
	_, dLog := logging.WhenRequest(ctx, "MemoryOrderRepository.Update ", map[string]interface{}{"order": order})
	defer dLog(nil, &err)

	m.lock.Lock()
	defer m.lock.Unlock()
	found := false
	for i, o := range m.store {
		if o.ID == order.ID && o.CustomerID == order.CustomerID {
			found = true
			updatedOrder, err := UpdateFn(ctx, order)
			if err != nil {
				return err
			}
			m.store[i] = updatedOrder
		}
	}
	if !found {
		return domain.NotFoundError{OrderID: order.ID}
	}
	return nil
}
