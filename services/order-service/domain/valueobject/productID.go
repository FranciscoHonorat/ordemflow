package valueobject

import (
	"encoding/json"
	"fmt"

	"github.com/FranciscoHonorat/ordemflow/services/order-service/domain/errors"
	"github.com/google/uuid"
)

type ProductID struct {
	id uuid.UUID
}

func NewProductID(id uuid.UUID) (ProductID, error) {
	if id == uuid.Nil {
		return ProductID{}, errors.ErrInvalidProductID
	}
	return ProductID{id: id}, nil
}

func (p ProductID) ID() uuid.UUID {
	return p.id
}

func (p ProductID) String() string {
	return p.id.String()
}

func (p ProductID) Equal(o ProductID) bool {
	return p.id == o.id
}

func (p ProductID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"id": "%s"}`, p.id.String())), nil
}

func (p *ProductID) UnmarshalJSON(data []byte) error {
	var product struct {
		ID uuid.UUID
	}

	if err := json.Unmarshal(data, &product); err != nil {
		return err
	}

	if product.ID == uuid.Nil {
		return errors.ErrInvalidProductID
	}

	p.id = product.ID
	return nil
}
