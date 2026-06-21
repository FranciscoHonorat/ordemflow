package valueobject

import (
	"encoding/json"
	"fmt"

	"github.com/FranciscoHonorat/ordemflow/services/order-service/domain/errors"
	"github.com/google/uuid"
)

type OrderID struct {
	id uuid.UUID
}

func NewOrderID(id uuid.UUID) (OrderID, error) {
	if id == uuid.Nil {
		return OrderID{}, errors.ErrInvalidOrderID
	}

	return OrderID{id: id}, nil
}

func (o OrderID) ID() uuid.UUID {
	return o.id
}

func (o OrderID) String() string {
	return o.id.String()
}

func (o OrderID) Equal(other OrderID) bool {
	return o.id == other.id
}

func (o OrderID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"id": "%s"}`, o.id.String())), nil
}

func (o *OrderID) UnmarshalJSON(data []byte) error {
	var order struct {
		ID uuid.UUID
	}

	if err := json.Unmarshal(data, &order); err != nil {
		return err
	}

	if order.ID == uuid.Nil {
		return errors.ErrInvalidOrderID
	}

	o.id = order.ID
	return nil
}
