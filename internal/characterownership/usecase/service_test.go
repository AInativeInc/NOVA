package usecase

import (
	"testing"

	"github.com/AInativeInc/NOVA/internal/characterownership/domain"
)

func TestValidateOwnership(t *testing.T) {
	svc := Service{}
	err := svc.ValidateOwnership(domain.CharacterOwnership{Splits: []domain.OwnershipSplit{{OwnerID: "a", PercentageBPS: 6000}, {OwnerID: "b", PercentageBPS: 4000}}})
	if err != nil {
		t.Fatalf("expected valid splits, got %v", err)
	}
}

func TestValidateOwnership_Invalid(t *testing.T) {
	svc := Service{}
	err := svc.ValidateOwnership(domain.CharacterOwnership{Splits: []domain.OwnershipSplit{{OwnerID: "a", PercentageBPS: 6000}, {OwnerID: "b", PercentageBPS: 3000}}})
	if err == nil {
		t.Fatal("expected split validation error")
	}
}

func TestValidateOwnership_DuplicateOwner(t *testing.T) {
	svc := Service{}
	err := svc.ValidateOwnership(domain.CharacterOwnership{Splits: []domain.OwnershipSplit{{OwnerID: "a", PercentageBPS: 6000}, {OwnerID: "a", PercentageBPS: 4000}}})
	if err == nil {
		t.Fatal("expected split validation error")
	}
}
