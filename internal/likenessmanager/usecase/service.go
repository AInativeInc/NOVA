package usecase

import (
	"errors"
	"time"

	"github.com/AInativeInc/NOVA/internal/likenessmanager/domain"
	"github.com/google/uuid"
)

var ErrInactiveConsent = errors.New("consent agreement is not active")

type Service struct{}

func (Service) NewModel(personID, displayName, verificationID string, metadata map[string]string) domain.RealPersonModel {
	return domain.RealPersonModel{
		ID:             uuid.NewString(),
		PersonID:       personID,
		DisplayName:    displayName,
		VerificationID: verificationID,
		Metadata:       metadata,
		CreatedAt:      time.Now().UTC(),
	}
}

func (Service) RequireActiveConsent(consent domain.ConsentAgreement, at time.Time) error {
	if !consent.IsActive(at) {
		return ErrInactiveConsent
	}
	return nil
}

func (Service) RevokeConsent(consent *domain.ConsentAgreement, revokedAt time.Time) {
	consent.RevokedAt = &revokedAt
}
