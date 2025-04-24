package convertor

import "sync"

var (
	orderConvertor *OrderConvertor
	orderOnce      sync.Once
)

func NewOrderConvertor() *OrderConvertor {
	orderOnce.Do(func() {
		orderConvertor = new(OrderConvertor)
	})
	return orderConvertor
}

var (
	itemConvertor *ItemConvertor
	itemOnce      sync.Once
)

func NewItemConvertor() *ItemConvertor {
	itemOnce.Do(func() {
		itemConvertor = new(ItemConvertor)
	})
	return itemConvertor
}

var (
	itemWithQuantityConvertor *ItemWithQuantityConvertor
	itemWithQuantityOnce      sync.Once
)

func NewItemWithQuantityConvertor() *ItemWithQuantityConvertor {
	itemWithQuantityOnce.Do(func() {
		itemWithQuantityConvertor = new(ItemWithQuantityConvertor)
	})
	return itemWithQuantityConvertor
}
