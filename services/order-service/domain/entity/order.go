package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID
	Value     int64
	Itens     int64
	Address   string
	Status    string
	CreatedAt time.Time
	UpdateAt  time.Time
}

func NewOrder(ID uuid.UUID, Value, Itens int64, Address, Status string) (Order, error) {
	return Order{
		ID:        ID,
		Value:     Value,
		Itens:     Itens,
		Address:   Address,
		Status:    Status,
		CreatedAt: time.Time{},
		UpdateAt:  time.Time{},
	}, nil
}
