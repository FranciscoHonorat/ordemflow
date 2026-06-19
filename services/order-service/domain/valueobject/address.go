package valueobject

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/FranciscoHonorat/ordemflow/services/order-service/domain/errors"
)

var cepRegex = regexp.MustCompile(`^\d{8}$`)

type Address struct {
	CEP            string
	Street         string
	Neighborhood   string
	Number         int64
	ReferencePoint string
	Complement     string
}

func NewAddrees(CEP string, Number int64, Street, Neighborhood, ReferencePoint, Complement string) (Address, error) {
	if !cepRegex.MatchString(CEP) {
		return Address{}, errors.ErrInvalidCEP
	}

	if Street == "" {
		return Address{}, errors.ErrFieldEmpty
	}

	if Neighborhood == "" {
		return Address{}, errors.ErrFieldEmpty
	}

	if Number <= 0 {
		return Address{}, errors.ErrInvalidNumber
	}

	return Address{
		CEP:            CEP,
		Street:         Street,
		Neighborhood:   Neighborhood,
		Number:         Number,
		ReferencePoint: ReferencePoint,
		Complement:     Complement,
	}, nil
}

func (a Address) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"cep": "%s", "street": "%s", "neighborhood": "%s", "number": %d, "referencePoint": "%s", "complement": "%s"}`, a.CEP, a.Street, a.Neighborhood, a.Number, a.ReferencePoint, a.Complement)), nil
}

func (a *Address) UnmarshalJSON(data []byte) error {
	var address struct {
		CEP            string `json:"cep"`
		Street         string `json:"street"`
		Neighborhood   string `json:"neighborhood"`
		Number         int64  `json:"number"`
		ReferencePoint string `json:"referencePoint"`
		Complement     string `json:"complement"`
	}

	if err := json.Unmarshal(data, &address); err != nil {
		return err
	}

	if !cepRegex.MatchString(address.CEP) {
		return errors.ErrInvalidCEP
	}

	if address.Street == "" {
		return errors.ErrFieldEmpty
	}

	if address.Neighborhood == "" {
		return errors.ErrFieldEmpty
	}

	if address.Number <= 0 {
		return errors.ErrInvalidNumber
	}

	a.CEP = address.CEP
	a.Street = address.Street
	a.Neighborhood = address.Neighborhood
	a.Number = address.Number
	a.ReferencePoint = address.ReferencePoint
	a.Complement = address.Complement
	return nil
}
