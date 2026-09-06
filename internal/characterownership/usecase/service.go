package usecase

import (
	"errors"

	"github.com/AInativeInc/NOVA/internal/characterownership/domain"
)

var ErrInvalidSplit = errors.New("ownership splits must total 100%")

type Service struct{}

func (Service) ValidateOwnership(c domain.CharacterOwnership) error {
	if len(c.Splits) == 0 {
		return ErrInvalidSplit
	}
	total := int32(0)
	owners := make(map[string]struct{}, len(c.Splits))
	for _, s := range c.Splits {
		if s.PercentageBPS <= 0 {
			return ErrInvalidSplit
		}
		if _, exists := owners[s.OwnerID]; exists {
			return ErrInvalidSplit
		}
		owners[s.OwnerID] = struct{}{}
		total += s.PercentageBPS
	}
	if total != 10000 {
		return ErrInvalidSplit
	}
	return nil
}
