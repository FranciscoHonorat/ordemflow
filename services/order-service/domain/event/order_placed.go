package event

import "github.com/FranciscoHonorat/ordemflow/services/order-service/domain/valueobject"

type OrderPlaced struct {
	BaseEvent
	CustomerID  string
	TotalAmount valueobject.Money
	ItemCount   int
}

func NewOrderPlaced(orderID, customerID string, total valueobject.Money, item int) OrderPlaced {
	return OrderPlaced{
		BaseEvent:   NewBaseEvent("order.placed", orderID),
		CustomerID:  customerID,
		TotalAmount: total,
		ItemCount:   item,
	}
}
