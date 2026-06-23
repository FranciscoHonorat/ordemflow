package valueobject_test

import (
	"testing"

	status "github.com/FranciscoHonorat/ordemflow/services/order-service/domain/valueobject"
)

func TestOrderStatus(t *testing.T) {
	t.Run("Test for IsValid method", func(t *testing.T) {
		tests := []struct {
			name   string
			status status.OrderStatus
			want   bool
		}{
			{"Valid status PENDING", status.OrderStatusPending, true},
			{"Valid status CONFIRMED", status.OrderStatusConfirmed, true},
			{"Valid status SHIPPED", status.OrderStatusShipped, true},
			{"Valid status DELIVERED", status.OrderStatusDelivered, true},
			{"Valid status CANCELLED", status.OrderStatusCancelled, true},
			{"Invalid status", "INVALID_STATUS", false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.status.IsValid(); got != tt.want {
					t.Errorf("IsValid() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("Test for String method", func(t *testing.T) {
		tests := []struct {
			name   string
			status status.OrderStatus
			want   string
		}{
			{"Status PENDING", status.OrderStatusPending, "PENDING"},
			{"Status CONFIRMED", status.OrderStatusConfirmed, "CONFIRMED"},
			{"Status SHIPPED", status.OrderStatusShipped, "SHIPPED"},
			{"Status DELIVERED", status.OrderStatusDelivered, "DELIVERED"},
			{"Status CANCELLED", status.OrderStatusCancelled, "CANCELLED"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.status.String(); got != tt.want {
					t.Errorf("String() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("Test for MarshalJSON and UnmarshalJSON methods", func(t *testing.T) {
		tests := []struct {
			name   string
			status status.OrderStatus
		}{
			{"Status PENDING", status.OrderStatusPending},
			{"Status CONFIRMED", status.OrderStatusConfirmed},
			{"Status SHIPPED", status.OrderStatusShipped},
			{"Status DELIVERED", status.OrderStatusDelivered},
			{"Status CANCELLED", status.OrderStatusCancelled},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				jsonData, err := tt.status.MarshalJSON()
				if err != nil {
					t.Errorf("MarshalJSON() error = %v", err)
				}

				var unmarshaledStatus status.OrderStatus
				if err := unmarshaledStatus.UnmarshalJSON(jsonData); err != nil {
					t.Errorf("UnmarshalJSON() error = %v", err)
				}

				if unmarshaledStatus != tt.status {
					t.Errorf("Unmarshaled status = %v, want %v", unmarshaledStatus, tt.status)
				}
			})
		}
	})
}
