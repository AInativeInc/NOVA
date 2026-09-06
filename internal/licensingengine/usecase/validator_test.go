package usecase

import (
	"testing"
	"time"

	"github.com/AInativeInc/NOVA/internal/licensingengine/domain"
	likeness "github.com/AInativeInc/NOVA/internal/likenessmanager/domain"
)

func TestEvaluate_ModelMismatch(t *testing.T) {
	fixedNow := time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC)
	v := Validator{Now: func() time.Time { return fixedNow }}
	consent := likeness.ConsentAgreement{
		ModelID:  "model-a",
		StartsAt: fixedNow.Add(-time.Hour),
		EndsAt:   fixedNow.Add(time.Hour),
	}
	d := v.Evaluate(domain.UsageRequest{ModelID: "model-b", HasValidKYC: true}, consent)
	if d.Allowed {
		t.Fatal("expected denied decision for mismatched model")
	}
}

func TestEvaluate_ConsentInactive(t *testing.T) {
	fixedNow := time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC)
	v := Validator{Now: func() time.Time { return fixedNow }}
	consent := likeness.ConsentAgreement{
		ModelID:  "model-a",
		StartsAt: fixedNow.Add(-2 * time.Hour),
		EndsAt:   fixedNow.Add(-time.Hour),
	}
	d := v.Evaluate(domain.UsageRequest{ModelID: "model-a", HasValidKYC: true}, consent)
	if d.Allowed {
		t.Fatal("expected denied decision for inactive consent")
	}
}

func TestEvaluate_MissingKYC(t *testing.T) {
	fixedNow := time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC)
	v := Validator{Now: func() time.Time { return fixedNow }}
	consent := likeness.ConsentAgreement{
		ModelID:  "model-a",
		StartsAt: fixedNow.Add(-time.Hour),
		EndsAt:   fixedNow.Add(time.Hour),
	}
	d := v.Evaluate(domain.UsageRequest{ModelID: "model-a", HasValidKYC: false}, consent)
	if d.Allowed {
		t.Fatal("expected denied decision for missing KYC")
	}
}

func TestEvaluate_TerritoryAndRestrictions(t *testing.T) {
	fixedNow := time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC)
	v := Validator{Now: func() time.Time { return fixedNow }}
	consent := likeness.ConsentAgreement{
		ModelID:        "model-a",
		StartsAt:       fixedNow.Add(-time.Hour),
		EndsAt:         fixedNow.Add(time.Hour),
		Territories:    []string{" US "},
		AllowedUses:    []string{" marketing "},
		RestrictedUses: []string{"political"},
	}

	deniedTerritory := v.Evaluate(domain.UsageRequest{
		ModelID:     "model-a",
		IntendedUse: "marketing",
		Territory:   "eu",
		HasValidKYC: true,
	}, consent)
	if deniedTerritory.Allowed {
		t.Fatal("expected denied decision for disallowed territory")
	}

	deniedRestricted := v.Evaluate(domain.UsageRequest{
		ModelID:     "model-a",
		IntendedUse: "Political",
		Territory:   "US",
		HasValidKYC: true,
	}, consent)
	if deniedRestricted.Allowed {
		t.Fatal("expected denied decision for restricted use")
	}

	deniedAllowlist := v.Evaluate(domain.UsageRequest{
		ModelID:     "model-a",
		IntendedUse: "Editorial",
		Territory:   "US",
		HasValidKYC: true,
	}, consent)
	if deniedAllowlist.Allowed {
		t.Fatal("expected denied decision for use outside allowlist")
	}

	allowed := v.Evaluate(domain.UsageRequest{
		ModelID:     "model-a",
		IntendedUse: " Marketing ",
		Territory:   " us ",
		HasValidKYC: true,
	}, consent)
	if !allowed.Allowed {
		t.Fatal("expected approved decision")
	}
}

func TestEvaluate_DefaultClockPath(t *testing.T) {
	v := Validator{}
	now := time.Now().UTC()
	consent := likeness.ConsentAgreement{
		ModelID:        "model-a",
		StartsAt:       now.Add(-24 * time.Hour),
		EndsAt:         now.Add(24 * time.Hour),
		Territories:    []string{"US"},
		AllowedUses:    []string{"marketing"},
		RestrictedUses: []string{"political"},
	}
	d := v.Evaluate(domain.UsageRequest{
		ModelID:     "model-a",
		IntendedUse: "marketing",
		Territory:   "US",
		HasValidKYC: true,
	}, consent)
	if !d.Allowed {
		t.Fatalf("expected approved decision, got %s", d.Reason)
	}
}
