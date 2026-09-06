package domain

import "time"

type RoyaltyTransaction struct {
	ID              string
	CharacterID     string
	GrossRevenue    float64
	Currency        string
	CreatedAt       time.Time
	PayoutByOwnerID map[string]float64
}
