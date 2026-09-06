package usecase

import (
	"errors"
	"math"

	"github.com/AInativeInc/NOVA/internal/characterownership/domain"
)

var ErrInvalidSplit = errors.New("ownership splits must total 100%")

type Service struct{}

func (Service) ValidateOwnership(c domain.CharacterOwnership) error {
	if len(c.Splits) == 0 {
		return ErrInvalidSplit
	}
	total := 0.0
	for _, s := range c.Splits {
		if s.Percentage <= 0 {
			return ErrInvalidSplit
		}
		total += s.Percentage
	}
	if math.Abs(total-100.0) > 0.0001 {
		return ErrInvalidSplit
	}
	return nil
}
