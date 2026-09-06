package usecase

import (
	"slices"

	"github.com/AInativeInc/NOVA/internal/licensingengine/domain"
	likeness "github.com/AInativeInc/NOVA/internal/likenessmanager/domain"
)

type Validator struct{}

func (Validator) Evaluate(req domain.UsageRequest, consent likeness.ConsentAgreement) domain.LicensingDecision {
	if !consent.IsActive(nowUTC()) {
		return domain.LicensingDecision{Allowed: false, Reason: "consent inactive"}
	}
	if !req.HasValidKYC {
		return domain.LicensingDecision{Allowed: false, Reason: "requester not verified"}
	}
	if len(consent.Territories) > 0 && !slices.Contains(consent.Territories, req.Territory) {
		return domain.LicensingDecision{Allowed: false, Reason: "territory not allowed"}
	}
	if slices.Contains(consent.RestrictedUses, req.IntendedUse) {
		return domain.LicensingDecision{Allowed: false, Reason: "use restricted by consent"}
	}
	if len(consent.AllowedUses) > 0 && !slices.Contains(consent.AllowedUses, req.IntendedUse) {
		return domain.LicensingDecision{Allowed: false, Reason: "use not in consent allowlist"}
	}
	return domain.LicensingDecision{Allowed: true, Reason: "approved"}
}
