package convertor

import (
	"github.com/mrluzy/gorder-v2/common/genproto/orderpb"
	"github.com/mrluzy/gorder-v2/stock/entity"
)

type OrderConvertor struct{}
type ItemConvertor struct{}
type ItemWithQuantityConvertor struct{}

func (c *OrderConvertor) EntityToProto(o *entity.Order) *orderpb.Order {
	c.check(o)
	return &orderpb.Order{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		Items:       NewItemConvertor().EntitiesToProtos(o.Items),
		PaymentLink: o.PaymentLink,
	}
}

func (c *OrderConvertor) ProtoToEntity(o *orderpb.Order) *entity.Order {
	c.check(o)
	return &entity.Order{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		PaymentLink: o.PaymentLink,
		Items:       NewItemConvertor().ProtosToEntities(o.Items),
	}
}

func (c *OrderConvertor) check(o interface{}) {
	if o == nil {
		panic("cannot convert nil order")
	}
}

func (c *ItemConvertor) ProtosToEntities(items []*orderpb.Item) (res []*entity.Item) {
	for _, item := range items {
		res = append(res, c.ProtoToEntity(item))
	}
	return
}

func (c *ItemConvertor) EntitiesToProtos(items []*entity.Item) (res []*orderpb.Item) {
	for _, item := range items {
		res = append(res, c.EntityToProeto(item))
	}
	return
}

func (c *ItemConvertor) ProtoToEntity(item *orderpb.Item) *entity.Item {
	return &entity.Item{
		ID:       item.ID,
		Name:     item.Name,
		Quantity: item.Quantity,
		PriceID:  item.PriceID,
	}
}

func (c *ItemConvertor) EntityToProeto(item *entity.Item) *orderpb.Item {
	return &orderpb.Item{
		ID:       item.ID,
		Name:     item.Name,
		Quantity: item.Quantity,
		PriceID:  item.PriceID,
	}
}

func (c ItemWithQuantityConvertor) EntitiesToProtos(items []*entity.ItemWithQuantity) (res []*orderpb.ItemWithQuantity) {
	for _, item := range items {
		res = append(res, c.EntityToProto(item))
	}
	return
}

func (c ItemWithQuantityConvertor) ProtosToEntities(items []*orderpb.ItemWithQuantity) (res []*entity.ItemWithQuantity) {
	for _, item := range items {
		res = append(res, c.ProtoToEntity(item))
	}
	return
}

func (c ItemWithQuantityConvertor) EntityToProto(item *entity.ItemWithQuantity) *orderpb.ItemWithQuantity {
	return &orderpb.ItemWithQuantity{
		ID:       item.ID,
		Quantity: item.Quantity,
	}
}

func (c ItemWithQuantityConvertor) ProtoToEntity(item *orderpb.ItemWithQuantity) *entity.ItemWithQuantity {
	return &entity.ItemWithQuantity{
		ID:       item.ID,
		Quantity: item.Quantity,
	}
}
