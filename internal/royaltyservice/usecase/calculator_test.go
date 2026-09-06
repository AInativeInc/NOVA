package usecase

import (
	"testing"

	ownership "github.com/AInativeInc/NOVA/internal/characterownership/domain"
)

func TestDistribute(t *testing.T) {
	calc := Calculator{}
	tx, err := calc.Distribute("char-1", 10000, "USD", []ownership.OwnershipSplit{{OwnerID: "a", PercentageBPS: 3333}, {OwnerID: "b", PercentageBPS: 3333}, {OwnerID: "c", PercentageBPS: 3334}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tx.PayoutByOwnerID["a"] != 3333 || tx.PayoutByOwnerID["b"] != 3333 || tx.PayoutByOwnerID["c"] != 3334 {
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

func TestDistribute_InvalidSplitTotal(t *testing.T) {
	calc := Calculator{}
	_, err := calc.Distribute("char-1", 100, "USD", []ownership.OwnershipSplit{{OwnerID: "a", PercentageBPS: 6000}, {OwnerID: "b", PercentageBPS: 3000}})
	if err == nil {
		t.Fatal("expected split validation error")
	}
}

func TestDistribute_DuplicateOwner(t *testing.T) {
	calc := Calculator{}
	_, err := calc.Distribute("char-1", 100, "USD", []ownership.OwnershipSplit{{OwnerID: "a", PercentageBPS: 6000}, {OwnerID: "a", PercentageBPS: 4000}})
	if err == nil {
		t.Fatal("expected duplicate owner validation error")
	}
}

func TestDistribute_EmptySplits(t *testing.T) {
	calc := Calculator{}
	_, err := calc.Distribute("char-1", 0, "USD", nil)
	if err == nil {
		t.Fatal("expected split validation error")
	}
}

func TestDistribute_DeterministicAcrossInputOrder(t *testing.T) {
	calc := Calculator{}
	first, err := calc.Distribute("char-1", 1, "USD", []ownership.OwnershipSplit{
		{OwnerID: "b", PercentageBPS: 5000},
		{OwnerID: "a", PercentageBPS: 5000},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := calc.Distribute("char-1", 1, "USD", []ownership.OwnershipSplit{
		{OwnerID: "a", PercentageBPS: 5000},
		{OwnerID: "b", PercentageBPS: 5000},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if first.PayoutByOwnerID["a"] != second.PayoutByOwnerID["a"] || first.PayoutByOwnerID["b"] != second.PayoutByOwnerID["b"] {
		t.Fatalf("expected deterministic payouts, got first=%v second=%v", first.PayoutByOwnerID, second.PayoutByOwnerID)
	}
}
