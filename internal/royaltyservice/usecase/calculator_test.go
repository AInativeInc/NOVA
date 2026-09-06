package usecase

import (
	"testing"

	ownership "github.com/AInativeInc/NOVA/internal/characterownership/domain"
)

func TestDistribute(t *testing.T) {
	calc := Calculator{}
	tx, err := calc.Distribute("char-1", 100.00, "USD", []ownership.OwnershipSplit{{OwnerID: "a", Percentage: 33.33}, {OwnerID: "b", Percentage: 33.33}, {OwnerID: "c", Percentage: 33.34}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tx.PayoutByOwnerID["a"] != 33.33 || tx.PayoutByOwnerID["b"] != 33.33 || tx.PayoutByOwnerID["c"] != 33.34 {
		t.Fatalf("unexpected payouts: %#v", tx.PayoutByOwnerID)
	}
}

func TestDistribute_NegativeRevenue(t *testing.T) {
	calc := Calculator{}
	_, err := calc.Distribute("char-1", -1, "USD", nil)
	if err == nil {
		t.Fatal("expected error for negative revenue")
	}
}
