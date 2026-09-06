package domain

import "time"

type OwnershipSplit struct {
	OwnerID    string
	Percentage float64
}

type CharacterOwnership struct {
	CharacterID string
	Splits      []OwnershipSplit
	UpdatedAt   time.Time
}
