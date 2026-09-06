package usecase

import (
	"errors"
	"slices"
	"time"

	ownership "github.com/AInativeInc/NOVA/internal/characterownership/domain"
	"github.com/AInativeInc/NOVA/internal/royaltyservice/domain"
	"github.com/google/uuid"
)

var ErrNegativeRevenue = errors.New("gross revenue must be non-negative")
var ErrInvalidSplits = errors.New("ownership splits must total 100%")
var ErrDuplicateOwner = errors.New("ownership splits contain duplicate owner")

type Calculator struct{}

func (Calculator) Distribute(characterID string, grossMinor int64, currency string, splits []ownership.OwnershipSplit) (domain.RoyaltyTransaction, error) {
	if grossMinor < 0 {
		return domain.RoyaltyTransaction{}, ErrNegativeRevenue
	}
	if len(splits) == 0 {
		return domain.RoyaltyTransaction{}, ErrInvalidSplits
	}
	total := int32(0)
	seenOwners := make(map[string]struct{}, len(splits))
	for _, split := range splits {
		if split.PercentageBPS <= 0 {
			return domain.RoyaltyTransaction{}, ErrInvalidSplits
		}
		if _, exists := seenOwners[split.OwnerID]; exists {
			return domain.RoyaltyTransaction{}, ErrDuplicateOwner
		}
		seenOwners[split.OwnerID] = struct{}{}
		total += split.PercentageBPS
	}
	if total != 10000 {
		return domain.RoyaltyTransaction{}, ErrInvalidSplits
	}
	sortedSplits := slices.Clone(splits)
	slices.SortFunc(sortedSplits, func(a, b ownership.OwnershipSplit) int {
		if a.OwnerID < b.OwnerID {
			return -1
		}
		if a.OwnerID > b.OwnerID {
			return 1
		}
		return 0
	})
	payouts := make(map[string]int64, len(splits))
	var running int64
	type remainderRank struct {
		ownerID   string
		remainder int64
	}
	ranks := make([]remainderRank, 0, len(sortedSplits))
	for _, split := range sortedSplits {
		numerator := grossMinor * int64(split.PercentageBPS)
		value := numerator / 10000
		payouts[split.OwnerID] = value
		running += value
		ranks = append(ranks, remainderRank{
			ownerID:   split.OwnerID,
			remainder: numerator % 10000,
		})
	}
	leftover := grossMinor - running
	slices.SortFunc(ranks, func(a, b remainderRank) int {
		if a.remainder > b.remainder {
			return -1
		}
		if a.remainder < b.remainder {
			return 1
		}
		if a.ownerID < b.ownerID {
			return -1
		}
		if a.ownerID > b.ownerID {
			return 1
		}
		return 0
	})
	for i := int64(0); i < leftover; i++ {
		payouts[ranks[i%int64(len(ranks))].ownerID]++
	}
	return domain.RoyaltyTransaction{
		ID:                uuid.NewString(),
		CharacterID:       characterID,
		GrossRevenueMinor: grossMinor,
		Currency:          currency,
		CreatedAt:         time.Now().UTC(),
		PayoutByOwnerID:   payouts,
	}, nil
}
