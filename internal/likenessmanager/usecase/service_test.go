package usecase

import (
	"testing"
	"time"

	"github.com/AInativeInc/NOVA/internal/likenessmanager/domain"
)

func TestRequireActiveConsent(t *testing.T) {
	now := time.Now().UTC()
	svc := Service{}
	consent := domain.ConsentAgreement{
		StartsAt: now.Add(-time.Hour),
		EndsAt:   now.Add(time.Hour),
	}
	if err := svc.RequireActiveConsent(consent, now); err != nil {
		t.Fatalf("expected active consent, got %v", err)
	}
}

func TestRequireActiveConsent_Revoked(t *testing.T) {
	now := time.Now().UTC()
	revoked := now.Add(-time.Minute)
	svc := Service{}
	consent := domain.ConsentAgreement{
		StartsAt:  now.Add(-time.Hour),
		EndsAt:    now.Add(time.Hour),
		RevokedAt: &revoked,
	}
	if err := svc.RequireActiveConsent(consent, now); err == nil {
		t.Fatal("expected inactive consent error")
	}
}
