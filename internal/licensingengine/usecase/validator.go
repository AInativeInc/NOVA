package usecase

import (
	"slices"
	"strings"
	"time"

	"github.com/AInativeInc/NOVA/internal/licensingengine/domain"
	likeness "github.com/AInativeInc/NOVA/internal/likenessmanager/domain"
)

type Validator struct {
	Now func() time.Time
}

func (v Validator) Evaluate(req domain.UsageRequest, consent likeness.ConsentAgreement) domain.LicensingDecision {
	if req.ModelID != consent.ModelID {
		return domain.LicensingDecision{Allowed: false, Reason: "consent does not belong to model"}
	}
	now := v.Now
	if now == nil {
		now = defaultNowUTC
	}
	if !consent.IsActive(now()) {
		return domain.LicensingDecision{Allowed: false, Reason: "consent inactive"}
	}
	if !req.HasValidKYC {
		return domain.LicensingDecision{Allowed: false, Reason: "requester not verified"}
	}
	reqTerritory := strings.ToUpper(req.Territory)
	reqUse := strings.ToLower(req.IntendedUse)
	territories := normalizeUpper(consent.Territories)
	restrictedUses := normalizeLower(consent.RestrictedUses)
	allowedUses := normalizeLower(consent.AllowedUses)

	if len(territories) > 0 && !slices.Contains(territories, reqTerritory) {
		return domain.LicensingDecision{Allowed: false, Reason: "territory not allowed"}
	}
	if slices.Contains(restrictedUses, reqUse) {
		return domain.LicensingDecision{Allowed: false, Reason: "use restricted by consent"}
	}
	if len(allowedUses) > 0 && !slices.Contains(allowedUses, reqUse) {
		return domain.LicensingDecision{Allowed: false, Reason: "use not in consent allowlist"}
	}
	return domain.LicensingDecision{Allowed: true, Reason: "approved"}
}

func normalizeUpper(values []string) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		out = append(out, strings.ToUpper(v))
	}
	return out
}

func normalizeLower(values []string) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		out = append(out, strings.ToLower(v))
	}
	return out
}
