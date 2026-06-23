package valueobject_test

import (
	"testing"

	product "github.com/FranciscoHonorat/ordemflow/services/order-service/domain/valueobject"
	"github.com/google/uuid"
)

func TestProductID(t *testing.T) {
	t.Run("Test for NewProductID", func(t *testing.T) {
		tests := []struct {
			name    string
			id      uuid.UUID
			wantErr bool
		}{
			{
				name:    "Valid ProductID",
				id:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				wantErr: false,
			},
			{
				name:    "Invalid ProductID (Nil UUID)",
				id:      uuid.Nil,
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := product.NewProductID(tt.id)
				if (err != nil) != tt.wantErr {
					t.Errorf("NewProductID() error = %v, wantErr %v", err, tt.wantErr)
				}
				if !tt.wantErr {
					pid, _ := product.NewProductID(tt.id)
					if pid.ID() != tt.id {
						t.Errorf("NewProductID() ID = %v, want %v", pid.ID(), tt.id)
					}
				}
			})
		}
	})

	t.Run("Test for NewProductIDMust", func(t *testing.T) {
		validID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
		pid := product.NewProductIDMust(validID)
		if pid.ID() != validID {
			t.Errorf("NewProductIDMust() ID = %v, want %v", pid.ID(), validID)
		}
	})

	t.Run("Test for String method", func(t *testing.T) {
		id := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
		pid, _ := product.NewProductID(id)
		if pid.String() != id.String() {
			t.Errorf("String() = %v, want %v", pid.String(), id.String())
		}
	})

	t.Run("Test for Equal method", func(t *testing.T) {
		id1 := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
		id2 := uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")

		pid1, _ := product.NewProductID(id1)
		pid2, _ := product.NewProductID(id1)
		pid3, _ := product.NewProductID(id2)

		if !pid1.Equal(pid2) {
			t.Errorf("Equal() = false, want true")
		}
		if pid1.Equal(pid3) {
			t.Errorf("Equal() = true, want false")
		}
	})

	t.Run("Test for MarshalJSON and UnmarshalJSON", func(t *testing.T) {
		id := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
		pid, _ := product.NewProductID(id)

		data, err := pid.MarshalJSON()
		if err != nil {
			t.Errorf("MarshalJSON() error = %v", err)
		}

		var newPID product.ProductID
		err = newPID.UnmarshalJSON(data)
		if err != nil {
			t.Errorf("UnmarshalJSON() error = %v", err)
		}

		if !pid.Equal(newPID) {
			t.Errorf("Unmarshaled ProductID does not match original")
		}
	})
}
