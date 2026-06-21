package entity

import (
	"fmt"
	"time"

	"github.com/FranciscoHonorat/ordemflow/services/order-service/domain/errors"
	"github.com/FranciscoHonorat/ordemflow/services/order-service/domain/valueobject"
)

type Order struct {
	id         valueobject.OrderID
	customerID valueobject.CustomerID
	totalPrice valueobject.Money
	items      []valueobject.OrderItem
	address    valueobject.Address
	status     valueobject.OrderStatus
	createdAt  time.Time
	updatedAt  time.Time
}

func NewOrder(id valueobject.OrderID, customerID valueobject.CustomerID) *Order {
	return &Order{
		id:         id,
		customerID: customerID,
		items:      []valueobject.OrderItem{},
		status:     valueobject.OrderStatusPending,
		createdAt:  time.Now(),
		updatedAt:  time.Now(),
	}
}

func (o *Order) recalculateTotal() error {
	var totalCents int64
	currency := "BRL"

	for _, item := range o.items {
		sub, err := item.SubTotal()
		if err != nil {
			return fmt.Errorf("recalculate total: %w", err)
		}
		totalCents += sub.Amount()
		currency = sub.Currency()
	}

	newTotal, err := valueobject.NewMoney(totalCents, currency)
	if err != nil {
		return fmt.Errorf("recalculate total: %w", err)
	}

	o.totalPrice = newTotal
	return nil
}

func (o *Order) AddItem(item valueobject.OrderItem) error {
	if o.status != valueobject.OrderStatusPending {
		return errors.ErrInvalidStatus
	}

	previousItem := o.items

	o.items = append(o.items, item)
	if err := o.recalculateTotal(); err != nil {
		o.items = previousItem
		return fmt.Errorf("add item: %w", err)
	}
	o.updatedAt = time.Now()
	return nil
}
