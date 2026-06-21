package valueobject

import (
	"fmt"
	"time"
)

type OrderItem struct {
	productID ProductID
	unitPrice Money
	quantity  Quantity
	createdAt time.Time
}

func NewOrderItem(productID ProductID, unitPrice Money, quantity Quantity) OrderItem {
	return OrderItem{
		productID: productID,
		unitPrice: unitPrice,
		quantity:  quantity,
		createdAt: time.Now(),
	}
}

func (i OrderItem) SubTotal() (Money, error) {
	subTotal, err := i.unitPrice.Multiply(i.quantity.Value())
	if err != nil {
		return Money{}, fmt.Errorf("subtotal: %w", err)
	}
	return subTotal, nil
}

func (i OrderItem) ProductID() ProductID {
	return i.productID
}

func (i OrderItem) UnitPrice() Money {
	return i.unitPrice
}

func (i OrderItem) Quantity() Quantity {
	return i.quantity
}

func (i OrderItem) CreatedAt() time.Time {
	return i.createdAt
}
