package entity_test

import (
	"testing"

	"github.com/FranciscoHonorat/ordemflow/services/order-service/domain/domainErrors"
	orderEntity "github.com/FranciscoHonorat/ordemflow/services/order-service/domain/entity"
	"github.com/FranciscoHonorat/ordemflow/services/order-service/domain/valueobject"
	"github.com/google/uuid"
)

func TestOrder(t *testing.T) {
	t.Run("Test for NewOrder", func(t *testing.T) {
		tests := []struct {
			name        string
			id          uuid.UUID
			customerID  uuid.UUID
			expectedErr error
		}{
			{"Valid Order", uuid.New(), uuid.New(), nil},
			{"Invalid Order ID", uuid.Nil, uuid.New(), domainErrors.ErrInvalidOrderID},
			{"Invalid Customer ID", uuid.New(), uuid.Nil, domainErrors.ErrInvalidCustomerID},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				orderID, err := valueobject.NewOrderID(tt.id)
				if err != nil {
					t.Errorf("Unexpected error creating order ID: %v", err)
				}
				customerID, err := valueobject.NewCustomerID(tt.customerID)
				if err != nil {
					t.Errorf("Unexpected error creating customer ID: %v", err)
				}
				_, err = orderEntity.NewOrder(orderID, customerID)
				if err != tt.expectedErr {
					t.Errorf("Expected error: %v, got: %v", tt.expectedErr, err)
				}
			})
		}
	})

	t.Run("Test for NewOrderMust", func(t *testing.T) {
		tests := []struct {
			name        string
			id          uuid.UUID
			customerID  uuid.UUID
			expectPanic bool
		}{
			{"Valid Order", uuid.New(), uuid.New(), false},
			{"Invalid Order ID", uuid.Nil, uuid.New(), true},
			{"Invalid Customer ID", uuid.New(), uuid.Nil, true},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				orderID, err := valueobject.NewOrderID(tt.id)
				if err != nil && !tt.expectPanic {
					t.Errorf("Unexpected error creating order ID: %v", err)
				}
				customerID, err := valueobject.NewCustomerID(tt.customerID)
				if err != nil && !tt.expectPanic {
					t.Errorf("Unexpected error creating customer ID: %v", err)
				}

				defer func() {
					if r := recover(); r != nil {
						if !tt.expectPanic {
							t.Errorf("Unexpected panic: %v", r)
						}
					} else {
						if tt.expectPanic {
							t.Errorf("Expected panic but did not get one")
						}
					}
				}()

				orderEntity.NewOrderMust(orderID, customerID)
			})
		}
	})

	t.Run("Test for recalculateTotal", func(t *testing.T) {

	})

	t.Run("Test for AddItem", func(t *testing.T) {

	})

	t.Run("Test for UpdateStatus", func(t *testing.T) {

	})

	t.Run("Test for ID and CustomerID", func(t *testing.T) {

	})

	t.Run("Test for TotalPrice", func(t *testing.T) {

	})

	t.Run("Test for Items", func(t *testing.T) {

	})
}
