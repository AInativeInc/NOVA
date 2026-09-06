package domain

import "time"

type RealPersonModel struct {
	ID             string
	PersonID       string
	DisplayName    string
	VerificationID string
	Metadata       map[string]string
	CreatedAt      time.Time
	RevokedAt      *time.Time
}

type ConsentAgreement struct {
	ID              string
	ModelID         string
	AllowedUses     []string
	RestrictedUses  []string
	Territories     []string
	StartsAt        time.Time
	EndsAt          time.Time
	RevokedAt       *time.Time
	VerificationRef string
}

func (c ConsentAgreement) IsActive(at time.Time) bool {
	if c.RevokedAt != nil {
		return false
	}
	return (at.Equal(c.StartsAt) || at.After(c.StartsAt)) && (at.Equal(c.EndsAt) || at.Before(c.EndsAt))
}
