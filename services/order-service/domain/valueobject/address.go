package valueobject

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/FranciscoHonorat/ordemflow/services/order-service/domain/errors"
)

var cepRegex = regexp.MustCompile(`^\d{8}$`)

type Address struct {
	cep            string
	street         string
	neighborhood   string
	number         int64
	referencePoint string
	complement     string
}

func NewAddress(CEP string, Number int64, Street, Neighborhood, ReferencePoint, Complement string) (Address, error) {
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
		cep:            CEP,
		street:         Street,
		neighborhood:   Neighborhood,
		number:         Number,
		referencePoint: ReferencePoint,
		complement:     Complement,
	}, nil
}

func (a Address) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"cep": "%s", "street": "%s", "neighborhood": "%s", "number": %d, "referencePoint": "%s", "complement": "%s"}`, a.cep, a.street, a.neighborhood, a.number, a.referencePoint, a.complement)), nil
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

	a.cep = address.CEP
	a.street = address.Street
	a.neighborhood = address.Neighborhood
	a.number = address.Number
	a.referencePoint = address.ReferencePoint
	a.complement = address.Complement
	return nil
}
