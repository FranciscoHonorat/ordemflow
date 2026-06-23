package valueobject_test

import (
	"testing"

	quantity "github.com/FranciscoHonorat/ordemflow/services/order-service/domain/valueobject"
)

func TestQuantity(t *testing.T) {
	t.Run("Test NewQuantity method", func(t *testing.T) {
		tests := []struct {
			name      string
			value     int64
			expectErr bool
		}{
			{
				name:      "Valid quantity",
				value:     10,
				expectErr: false,
			},
			{
				name:      "Invalid quantity (zero)",
				value:     0,
				expectErr: true,
			},
			{
				name:      "Invalid quantity (negative)",
				value:     -5,
				expectErr: true,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := quantity.NewQuantity(tt.value)
				if (err != nil) != tt.expectErr {
					t.Errorf("Expected error: %v, got: %v", tt.expectErr, err)
				}
			})
		}
	})
	t.Run("Test Equal method", func(t *testing.T) {
		q1, _ := quantity.NewQuantity(10)
		q2, _ := quantity.NewQuantity(10)
		q3, _ := quantity.NewQuantity(5)

		if !q1.Equal(q2) {
			t.Errorf("Expected quantities to be equal")
		}
		if q1.Equal(q3) {
			t.Errorf("Expected quantities to be not equal")
		}
	})

	t.Run("Test MarshalJSON and UnmarshalJSON methods", func(t *testing.T) {
		q, _ := quantity.NewQuantity(10)
		data, err := q.MarshalJSON()
		if err != nil {
			t.Errorf("Error marshaling JSON: %v", err)
		}

		var unmarshaledQuantity quantity.Quantity
		err = unmarshaledQuantity.UnmarshalJSON(data)
		if err != nil {
			t.Errorf("Error unmarshaling JSON: %v", err)
		}

		if !q.Equal(unmarshaledQuantity) {
			t.Errorf("Expected quantities to be equal after unmarshaling")
		}
	})
}
