package valueobject

import (
	"encoding/json"

	"github.com/FranciscoHonorat/ordemflow/services/order-service/domain/domainErrors"
)

type Quantity struct {
	value int64
}

func NewQuantity(value int64) (Quantity, error) {
	if value <= 0 {
		return Quantity{}, domainErrors.ErrInvalidQuantity
	}
	return Quantity{value: value}, nil
}

func NewQuantityMust(value int64) Quantity {
	q, err := NewQuantity(value)
	if err != nil {
		panic(err)
	}
	return q
}

func (q Quantity) Value() int64 {
	return q.value
}

func (q Quantity) Equal(o Quantity) bool {
	return q.value == o.value
}

func (q Quantity) MarshalJSON() ([]byte, error) {
	auxQuantity := struct {
		Value int64 `json:"value"`
	}{
		Value: q.value,
	}
	return json.Marshal(auxQuantity)
}

func (q *Quantity) UnmarshalJSON(data []byte) error {
	var quantity struct {
		Value int64 `json:"value"`
	}

	if err := json.Unmarshal(data, &quantity); err != nil {
		return err
	}

	if quantity.Value <= 0 {
		return domainErrors.ErrInvalidQuantity
	}

	q.value = quantity.Value
	return nil
}
