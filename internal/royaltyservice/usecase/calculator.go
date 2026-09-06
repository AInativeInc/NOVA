package usecase

import (
	"errors"
	"math"
	"time"

	ownership "github.com/AInativeInc/NOVA/internal/characterownership/domain"
	"github.com/AInativeInc/NOVA/internal/royaltyservice/domain"
	"github.com/google/uuid"
)

var ErrNegativeRevenue = errors.New("gross revenue must be non-negative")

type Calculator struct{}

func (Calculator) Distribute(characterID string, gross float64, currency string, splits []ownership.OwnershipSplit) (domain.RoyaltyTransaction, error) {
	if gross < 0 {
		return domain.RoyaltyTransaction{}, ErrNegativeRevenue
	}
	payouts := make(map[string]float64, len(splits))
	running := 0.0
	for i, split := range splits {
		value := round2(gross * split.Percentage / 100.0)
		if i == len(splits)-1 {
			value = round2(gross - running)
		}
		payouts[split.OwnerID] = value
		running = round2(running + value)
	}
	return domain.RoyaltyTransaction{
		ID:              uuid.NewString(),
		CharacterID:     characterID,
		GrossRevenue:    gross,
		Currency:        currency,
		CreatedAt:       time.Now().UTC(),
		PayoutByOwnerID: payouts,
	}, nil
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
