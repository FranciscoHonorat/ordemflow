package valueobject

import (
	"encoding/json"
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

func (i OrderItem) MarshalJSON() ([]byte, error) {
	auxOrderItem := struct {
		ProductID ProductID `json:"productID"`
		UnitPrice Money     `json:"unitPrice"`
		Quantity  Quantity  `json:"quantity"`
		CreatedAt string    `json:"createdAt"`
	}{
		ProductID: i.productID,
		UnitPrice: i.unitPrice,
		Quantity:  i.quantity,
		CreatedAt: i.createdAt.Format(time.RFC3339),
	}
	return json.Marshal(auxOrderItem)
}

func (i *OrderItem) UnmarshalJSON(data []byte) error {
	var orderItem struct {
		ProductID ProductID `json:"productID"`
		UnitPrice Money     `json:"unitPrice"`
		Quantity  Quantity  `json:"quantity"`
		CreatedAt string    `json:"createdAt"`
	}

	if err := json.Unmarshal(data, &orderItem); err != nil {
		return err
	}

	i.productID = orderItem.ProductID
	i.unitPrice = orderItem.UnitPrice
	i.quantity = orderItem.Quantity

	var err error
	i.createdAt, err = time.Parse(time.RFC3339, orderItem.CreatedAt)
	if err != nil {
		return fmt.Errorf("unmarshal createdAt: %w", err)
	}

	return nil
}
