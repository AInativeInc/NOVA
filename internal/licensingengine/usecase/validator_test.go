package usecase

import (
	"testing"
	"time"

	"github.com/AInativeInc/NOVA/internal/licensingengine/domain"
	likeness "github.com/AInativeInc/NOVA/internal/likenessmanager/domain"
)

func TestEvaluate_ModelMismatch(t *testing.T) {
	v := Validator{}
	consent := likeness.ConsentAgreement{
		ModelID:  "model-a",
		StartsAt: time.Now().UTC().Add(-time.Hour),
		EndsAt:   time.Now().UTC().Add(time.Hour),
	}
	d := v.Evaluate(domain.UsageRequest{ModelID: "model-b", HasValidKYC: true}, consent)
	if d.Allowed {
		t.Fatal("expected denied decision for mismatched model")
	}
}

func TestEvaluate_ConsentInactive(t *testing.T) {
	v := Validator{}
	consent := likeness.ConsentAgreement{
		ModelID:  "model-a",
		StartsAt: time.Now().UTC().Add(-2 * time.Hour),
		EndsAt:   time.Now().UTC().Add(-time.Hour),
	}
	d := v.Evaluate(domain.UsageRequest{ModelID: "model-a", HasValidKYC: true}, consent)
	if d.Allowed {
		t.Fatal("expected denied decision for inactive consent")
	}
}

func TestEvaluate_MissingKYC(t *testing.T) {
	v := Validator{}
	consent := likeness.ConsentAgreement{
		ModelID:  "model-a",
		StartsAt: time.Now().UTC().Add(-time.Hour),
		EndsAt:   time.Now().UTC().Add(time.Hour),
	}
	d := v.Evaluate(domain.UsageRequest{ModelID: "model-a", HasValidKYC: false}, consent)
	if d.Allowed {
		t.Fatal("expected denied decision for missing KYC")
	}
}

func TestEvaluate_TerritoryAndRestrictions(t *testing.T) {
	v := Validator{}
	consent := likeness.ConsentAgreement{
		ModelID:        "model-a",
		StartsAt:       time.Now().UTC().Add(-time.Hour),
		EndsAt:         time.Now().UTC().Add(time.Hour),
		Territories:    []string{"US"},
		AllowedUses:    []string{"marketing"},
		RestrictedUses: []string{"political"},
	}

	deniedTerritory := v.Evaluate(domain.UsageRequest{
		ModelID:     "model-a",
		IntendedUse: "marketing",
		Territory:   "EU",
		HasValidKYC: true,
	}, consent)
	if deniedTerritory.Allowed {
		t.Fatal("expected denied decision for disallowed territory")
	}

	deniedRestricted := v.Evaluate(domain.UsageRequest{
		ModelID:     "model-a",
		IntendedUse: "political",
		Territory:   "US",
		HasValidKYC: true,
	}, consent)
	if deniedRestricted.Allowed {
		t.Fatal("expected denied decision for restricted use")
	}

	deniedAllowlist := v.Evaluate(domain.UsageRequest{
		ModelID:     "model-a",
		IntendedUse: "editorial",
		Territory:   "US",
		HasValidKYC: true,
	}, consent)
	if deniedAllowlist.Allowed {
		t.Fatal("expected denied decision for use outside allowlist")
	}

	allowed := v.Evaluate(domain.UsageRequest{
		ModelID:     "model-a",
		IntendedUse: "marketing",
		Territory:   "US",
		HasValidKYC: true,
	}, consent)
	if !allowed.Allowed {
		t.Fatal("expected approved decision")
	}
}
