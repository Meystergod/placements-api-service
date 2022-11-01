package tile

import (
	"math"

	partnerImp "github.com/Meystergod/placements-api-service/internal/models/partner/imp"
)

type Tile struct {
	ID    uint    `json:"id" validate:"required"`
	Width uint    `json:"width" validate:"required,min=0"`
	Ratio float64 `json:"ratio" validate:"required,min=0"`
}

func (t *Tile) ConvertToPartnerRequestImp() *partnerImp.RequestImp {
	id := t.ID
	minWidth := t.Width
	minHeight := uint(math.Floor(float64(t.Width) * t.Ratio))

	return &partnerImp.RequestImp{
		ID:        id,
		MinWidth:  minWidth,
		MinHeight: minHeight,
	}
}
