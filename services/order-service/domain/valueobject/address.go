package valueobject

import (
	"encoding/json"
	"regexp"

	"github.com/FranciscoHonorat/ordemflow/services/order-service/domain/domainErrors"
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

func NewAddress(CEP, Street, Neighborhood string, Number int64, ReferencePoint, Complement string) (Address, error) {
	if !cepRegex.MatchString(CEP) {
		return Address{}, domainErrors.ErrInvalidCEP
	}

	if Street == "" {
		return Address{}, domainErrors.ErrFieldEmpty
	}

	if Neighborhood == "" {
		return Address{}, domainErrors.ErrFieldEmpty
	}

	if Number <= 0 {
		return Address{}, domainErrors.ErrInvalidNumber
	}

	if ReferencePoint == "" {
		return Address{}, domainErrors.ErrFieldEmpty
	}

	if Complement == "" {
		return Address{}, domainErrors.ErrFieldEmpty
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

func NewAddressMust(CEP, Street, Neighborhood string, Number int64, ReferencePoint, Complement string) Address {
	a, err := NewAddress(CEP, Street, Neighborhood, Number, ReferencePoint, Complement)
	if err != nil {
		panic(err)
	}
	return a
}

func (a Address) Cep() string {
	return a.cep
}

func (a Address) Street() string {
	return a.street
}

func (a Address) Neighborhood() string {
	return a.neighborhood
}

func (a Address) Number() int64 {
	return a.number
}

func (a Address) ReferencePoint() string {
	return a.referencePoint
}

func (a Address) Complement() string {
	return a.complement
}

func (a Address) MarshalJSON() ([]byte, error) {
	auxAddress := struct {
		CEP            string `json:"cep"`
		Street         string `json:"street"`
		Neighborhood   string `json:"neighborhood"`
		Number         int64  `json:"number"`
		ReferencePoint string `json:"referencePoint"`
		Complement     string `json:"complement"`
	}{
		CEP:            a.cep,
		Street:         a.street,
		Neighborhood:   a.neighborhood,
		Number:         a.number,
		ReferencePoint: a.referencePoint,
		Complement:     a.complement,
	}
	return json.Marshal(auxAddress)
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
		return domainErrors.ErrInvalidCEP
	}

	if address.Street == "" {
		return domainErrors.ErrFieldEmpty
	}

	if address.Neighborhood == "" {
		return domainErrors.ErrFieldEmpty
	}

	if address.Number <= 0 {
		return domainErrors.ErrInvalidNumber
	}

	a.cep = address.CEP
	a.street = address.Street
	a.neighborhood = address.Neighborhood
	a.number = address.Number
	a.referencePoint = address.ReferencePoint
	a.complement = address.Complement
	return nil
}
