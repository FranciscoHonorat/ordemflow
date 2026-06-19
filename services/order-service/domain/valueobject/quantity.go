package valueobject

import "fmt"

type Quantity struct {
	value int64
}

func NewQuantity(value int64) Quantity {
	return Quantity{value: value}
}

func (q Quantity) Value() int64 {
	return q.value
}

func (q Quantity) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", q.value)), nil
}

func (q *Quantity) UnmarshalJSON(data []byte) error {
	var value int64
	_, err := fmt.Sscanf(string(data), "%d", &value)
	if err != nil {
		return err
	}

	q.value = value
	return nil
}
