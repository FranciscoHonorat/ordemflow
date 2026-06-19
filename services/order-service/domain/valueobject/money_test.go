package valueobject

import (
	"testing"
	//mo "github.com/FranciscoHonorat/ordemflow/services/order-service/domain/valueobject"
)

var targetMoney Money
var targetPointer *Money

func BenchmarkNewMoney_Stack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m, _ := NewMoney(100, "BRL")
		targetMoney = m
	}
}

func BenchmarkNewMoney_Heap(b *testing.B) {
	newMoneyHeap := func(amount int64, currency string) (*Money, error) {
		m, err := NewMoney(amount, currency)
		return &m, err
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, _ := newMoneyHeap(100, "BRL")
		targetPointer = m
	}
}
