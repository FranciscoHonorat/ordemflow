package valueobject_test

import (
	"testing"

	"github.com/FranciscoHonorat/ordemflow/services/order-service/domain/domainErrors"
	orderid "github.com/FranciscoHonorat/ordemflow/services/order-service/domain/valueobject"
	"github.com/google/uuid"
)

func TestOrderID(t *testing.T) {
	t.Run("Test for NewOrderID", func(t *testing.T) {
		tests := []struct {
			name    string
			id      uuid.UUID
			wantErr error
		}{
			{"Valid id", uuid.New(), nil},
			{"Invalid id", uuid.Nil, domainErrors.ErrInvalidOrderID},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				p, err := orderid.NewOrderID(tt.id)
				if err != tt.wantErr {
					t.Errorf("Expected error: %v, got: %v", tt.wantErr, err)
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
		orderIDObj, _ := orderid.NewOrderID(id)

		if orderIDObj.String() != id.String() {
			t.Errorf("Expected String: %v, got: %v", id.String(), orderIDObj.String())
		}
	})

	t.Run("Test for Equal method", func(t *testing.T) {
		id1 := uuid.New()
		id2 := uuid.New()

		orderID1, _ := orderid.NewOrderID(id1)
		orderID2, _ := orderid.NewOrderID(id1)
		orderID3, _ := orderid.NewOrderID(id2)

		if !orderID1.Equal(orderID2) {
			t.Errorf("Expected orderID1 to be equal to orderID2")
		}

		if orderID1.Equal(orderID3) {
			t.Errorf("Expected orderID1 to not be equal to orderID3")
		}
	})

	t.Run("Test for MarshalJSON and UnmarshalJSON", func(t *testing.T) {
		id := uuid.New()
		orderIDObj, _ := orderid.NewOrderID(id)

		jsonData, err := orderIDObj.MarshalJSON()
		if err != nil {
			t.Errorf("Error marshalling OrderID: %v", err)
		}

		var unmarshalledOrderID orderid.OrderID
		err = unmarshalledOrderID.UnmarshalJSON(jsonData)
		if err != nil {
			t.Errorf("Error unmarshalling OrderID: %v", err)
		}

		if !orderIDObj.Equal(unmarshalledOrderID) {
			t.Errorf("Expected unmarshalled OrderID to be equal to original")
		}
	})

	t.Run("Test for UnmarshalJSON with invalid data", func(t *testing.T) {
		invalidJSON := []byte(`{"id": "00000000-0000-0000-0000-000000000000"}`)

		var unmarshalledOrderID orderid.OrderID
		err := unmarshalledOrderID.UnmarshalJSON(invalidJSON)
		if err != domainErrors.ErrInvalidOrderID {
			t.Errorf("Expected error: %v, got: %v", domainErrors.ErrInvalidOrderID, err)
		}
	})
}
