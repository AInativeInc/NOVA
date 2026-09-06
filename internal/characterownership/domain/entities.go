package domain

import "time"

type OwnershipSplit struct {
	OwnerID       string
	PercentageBPS int32
}

type CharacterOwnership struct {
	CharacterID string
	Splits      []OwnershipSplit
	UpdatedAt   time.Time
}
