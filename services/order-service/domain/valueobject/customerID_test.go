package valueobject_test

import (
	"errors"
	"testing"

	"github.com/FranciscoHonorat/ordemflow/services/order-service/domain/domainErrors"
	customer "github.com/FranciscoHonorat/ordemflow/services/order-service/domain/valueobject"
	"github.com/google/uuid"
)

func TestCustomerID(t *testing.T) {
	t.Run("Test for NewCustomerID", func(t *testing.T) {
		tests := []struct {
			name          string
			id            uuid.UUID
			expectedError error
		}{
			{"Valid id", uuid.New(), nil},
			{"Invalid id", uuid.Nil, domainErrors.ErrInvalidID},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				p, err := customer.NewCustomerID(tt.id)
				if !errors.Is(err, tt.expectedError) {
					t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
				}

				if err == nil {
					if p.ID() != tt.id {
						t.Errorf("Expected ID: %v, got: %v", tt.id, p.ID())
					}
				}
			})
		}
	})

	t.Run("Test for String method", func(t *testing.T) {
		id := uuid.New()
		c, _ := customer.NewCustomerID(id)

		if c.String() != id.String() {
			t.Errorf("Expected: %s, got: %s", id.String(), c.String())
		}
	})

	t.Run("Test for Equal method", func(t *testing.T) {
		id1 := uuid.New()
		id2 := uuid.New()

		tests := []struct {
			name     string
			p1       customer.CustomerID
			p2       customer.CustomerID
			expected bool
		}{
			{"Equal ID", customer.NewCustomerIDMust(id1), customer.NewCustomerIDMust(id1), true},
			{"Diferente ID", customer.NewCustomerIDMust(id1), customer.NewCustomerIDMust(id2), true},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.p1.Equal(tt.p2) != tt.expected {
					t.Errorf("Expected equality: %v, got: %v", tt.expected, tt.p1.Equal(tt.p2))
				}
			})
		}
	})

	t.Run("Test for MarhsalJSON", func(t *testing.T) {
		id1 := uuid.New()
		tests := []struct {
			name     string
			c        customer.CustomerID
			expected string
		}{
			{"ID", customer.NewCustomerIDMust(id1), `{"id":"` + id1.String() + `"}`},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				jsonData, err := tt.c.MarshalJSON()
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if string(jsonData) != tt.expected {
					t.Errorf("Expected JSON: %s, got: %s", tt.expected, string(jsonData))
				}
			})
		}
	})

	t.Run("Test for UnmarhsalJSON", func(t *testing.T) {
		id1 := uuid.New()
		tests := []struct {
			name     string
			c        customer.CustomerID
			expected string
		}{
			{"Valid JSON", customer.NewCustomerIDMust(id1), `{"id":"` + id1.String() + `"}`},
			{"Invalid JSON", customer.CustomerID{}, `{"id":"invalid-uuid"}`},
		}
		for _, tt := range tests {
			var c customer.CustomerID
			t.Run(tt.name, func(t *testing.T) {
				err := c.UnmarshalJSON([]byte(tt.expected))
				if tt.name == "Invalid JSON" {
					if !errors.Is(err, domainErrors.ErrInvalidID) {
						t.Errorf("Expected error: %v, got: %v", domainErrors.ErrInvalidID, err)
					}
					return
				}
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				} else {
					if c.ID() != tt.c.ID() {
						t.Errorf("Expected ID: %v, got: %v", tt.c.ID(), c.ID())
					}
				}
			})
		}
	})

}
