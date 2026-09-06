package usecase

import (
	"testing"

	"github.com/AInativeInc/NOVA/internal/characterownership/domain"
)

func TestValidateOwnership(t *testing.T) {
	svc := Service{}
	err := svc.ValidateOwnership(domain.CharacterOwnership{Splits: []domain.OwnershipSplit{{OwnerID: "a", Percentage: 60}, {OwnerID: "b", Percentage: 40}}})
	if err != nil {
		t.Fatalf("expected valid splits, got %v", err)
	}
}

func TestValidateOwnership_Invalid(t *testing.T) {
	svc := Service{}
	err := svc.ValidateOwnership(domain.CharacterOwnership{Splits: []domain.OwnershipSplit{{OwnerID: "a", Percentage: 60}, {OwnerID: "b", Percentage: 30}}})
	if err == nil {
		t.Fatal("expected split validation error")
	}
}
