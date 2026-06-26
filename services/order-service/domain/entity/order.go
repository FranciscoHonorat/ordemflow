package entity

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/FranciscoHonorat/ordemflow/services/order-service/domain/domainErrors"
	"github.com/FranciscoHonorat/ordemflow/services/order-service/domain/event"
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
	event      []event.DomainEvent
}

func NewOrder(id valueobject.OrderID, customerID valueobject.CustomerID) (*Order, error) {
	if id == (valueobject.OrderID{}) {
		return nil, domainErrors.ErrInvalidOrderID
	}
	if customerID == (valueobject.CustomerID{}) {
		return nil, domainErrors.ErrInvalidCustomerID
	}

	return &Order{
		id:         id,
		customerID: customerID,
		items:      []valueobject.OrderItem{},
		status:     valueobject.OrderStatusPending,
		createdAt:  time.Now(),
		updatedAt:  time.Now(),
	}, nil
}

func NewOrderMust(id valueobject.OrderID, customerID valueobject.CustomerID) *Order {
	order, err := NewOrder(id, customerID)
	if err != nil {
		panic(err)
	}
	return order
}

func (o *Order) recalculateTotal() error {
	if len(o.items) == 0 {
		o.totalPrice = valueobject.Money{} // Reset total price to zero if there are no items
		return nil
	}

	total, err := o.items[0].SubTotal()
	if err != nil {
		return fmt.Errorf("Recalculate total: %w", err)
	}

	for _, item := range o.items[1:] {
		sub, err := item.SubTotal()
		if err != nil {
			return fmt.Errorf("Recalculate total: %w", err)
		}
		total, err = total.Add(sub)
		if err != nil {
			return fmt.Errorf("Recalculate total: %w", err)
		}
	}
	o.totalPrice = total
	return nil
}

func (o *Order) AddItem(item valueobject.OrderItem) error {
	previous := make([]valueobject.OrderItem, len(o.items))
	copy(previous, o.items)
	o.items = append(o.items, item)

	if err := o.recalculateTotal(); err != nil {
		o.items = previous
		return err
	}
	return nil
}

func (o *Order) UpdateStatus(newStatus valueobject.OrderStatus) error {
	if !newStatus.IsValid() {
		return domainErrors.ErrInvalidOrderStatus
	}
	o.status = newStatus
	o.updatedAt = time.Now()
	return nil
}

func (o *Order) Place() {
	evt := event.NewOrderPlaced(
		o.id.String(),
		o.customerID.String(),
		o.totalPrice,
		len(o.items),
	)
	o.event = append(o.event, evt)
}

func (o *Order) DomainEvent() []event.DomainEvent {
	cp := make([]event.DomainEvent, len(o.event))
	copy(cp, o.event)
	return cp
}

func (o *Order) ClearEvent() {
	o.event = nil
}

func (o *Order) ID() valueobject.OrderID {
	return o.id
}

func (o *Order) CustomerID() valueobject.CustomerID {
	return o.customerID
}

func (o *Order) TotalPrice() valueobject.Money {
	return o.totalPrice
}

func (o *Order) Items() []valueobject.OrderItem {
	result := make([]valueobject.OrderItem, len(o.items))
	copy(result, o.items)
	return result
}

func (o *Order) Address() valueobject.Address {
	return o.address
}

func (o *Order) Status() valueobject.OrderStatus {
	return o.status
}

func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

func (o *Order) UpdatedAt() time.Time {
	return o.updatedAt
}

func (o *Order) MarshalJSON() ([]byte, error) {
	order := struct {
		ID         valueobject.OrderID     `json:"id"`
		CustomerID valueobject.CustomerID  `json:"customerID"`
		TotalPrice valueobject.Money       `json:"totalPrice"`
		Items      []valueobject.OrderItem `json:"items"`
		Address    valueobject.Address     `json:"address"`
		Status     valueobject.OrderStatus `json:"status"`
		CreatedAt  string                  `json:"createdAt"`
		UpdatedAt  string                  `json:"updatedAt"`
	}{
		ID:         o.id,
		CustomerID: o.customerID,
		TotalPrice: o.totalPrice,
		Items:      o.Items(),
		Address:    o.address,
		Status:     o.status,
		CreatedAt:  o.createdAt.Format(time.RFC3339),
		UpdatedAt:  o.updatedAt.Format(time.RFC3339),
	}
	return json.Marshal(order)
}

func (o *Order) UnmarshalJSON(data []byte) error {
	var order struct {
		ID         valueobject.OrderID     `json:"id"`
		CustomerID valueobject.CustomerID  `json:"customerID"`
		TotalPrice valueobject.Money       `json:"totalPrice"`
		Items      []valueobject.OrderItem `json:"items"`
		Address    valueobject.Address     `json:"address"`
		Status     valueobject.OrderStatus `json:"status"`
		CreatedAt  string                  `json:"createdAt"`
		UpdatedAt  string                  `json:"updatedAt"`
	}

	if err := json.Unmarshal(data, &order); err != nil {
		return err
	}

	o.id = order.ID
	o.customerID = order.CustomerID
	o.totalPrice = order.TotalPrice
	o.items = order.Items
	o.address = order.Address
	o.status = valueobject.OrderStatus(order.Status)

	var err error
	o.createdAt, err = time.Parse(time.RFC3339, order.CreatedAt)
	if err != nil {
		return fmt.Errorf("unmarshal createdAt: %w", err)
	}

	o.updatedAt, err = time.Parse(time.RFC3339, order.UpdatedAt)
	if err != nil {
		return fmt.Errorf("unmarshal updatedAt: %w", err)
	}

	return nil
}
