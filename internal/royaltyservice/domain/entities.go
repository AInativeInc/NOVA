package domain

import "time"

type RoyaltyTransaction struct {
	ID                string
	CharacterID       string
	GrossRevenueMinor int64
	Currency          string
	CreatedAt         time.Time
	PayoutByOwnerID   map[string]int64
}
