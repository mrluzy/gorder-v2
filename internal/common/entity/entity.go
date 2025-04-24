package entity

import (
	"errors"
	"fmt"
	"strings"
)

type Item struct {
	ID       string
	Name     string
	Quantity int32
	PriceID  string
}

func (i Item) validate() error {
	var invalidFields []string
	if i.ID == "" {
		invalidFields = append(invalidFields, "ID")
	}
	if i.Name == "" {
		invalidFields = append(invalidFields, "Name")
	}
	if i.PriceID == "" {
		invalidFields = append(invalidFields, "PriceID")
	}
	return fmt.Errorf("item %v invalid, empty fields: [%s]", i, strings.Join(invalidFields, ", "))
}

func NewItem(ID string, name string, quantity int32, priceID string) *Item {
	return &Item{ID: ID, Name: name, Quantity: quantity, PriceID: priceID}
}

func NewValidateItem(ID string, name string, quantity int32, priceID string) (*Item, error) {
	item := NewItem(ID, name, quantity, priceID)
	if err := item.validate(); err != nil {
		return nil, err
	}
	return item, nil
}

type ItemWithQuantity struct {
	ID       string
	Quantity int32
}

func (i ItemWithQuantity) validate() error {
	var invalidFields []string
	if i.ID == "" {
		invalidFields = append(invalidFields, "ID")
	}
	return errors.New(strings.Join(invalidFields, ", "))
}

func NewItemWithQuantity(ID string, quantity int32) *ItemWithQuantity {
	return &ItemWithQuantity{Quantity: quantity, ID: ID}
}

func NewValidateItemWithQuantity(ID string, quantity int32) (*ItemWithQuantity, error) {
	item := NewItemWithQuantity(ID, quantity)
	if err := item.validate(); err != nil {
		return nil, err
	}
	return item, nil
}

type Order struct {
	ID          string
	CustomerID  string
	Status      string
	PaymentLink string
	Items       []*Item
}

func NewOrder(ID string, customerID string, status string, paymentLink string, items []*Item) *Order {
	return &Order{ID: ID, CustomerID: customerID, Status: status, PaymentLink: paymentLink, Items: items}
}

func NewValidateOrder(ID string, customerID string, status string, paymentLink string, items []*Item) (*Order, error) {
	for _, item := range items {
		if err := item.validate(); err != nil {
			return nil, err
		}
	}
	return NewOrder(ID, customerID, status, paymentLink, items), nil
}
