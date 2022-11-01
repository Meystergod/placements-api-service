package imp

import (
	partnerImp "github.com/Meystergod/placements-api-service/internal/models/partner/imp"
)

type Imp struct {
	ID     uint   `json:"id"`
	Width  uint   `json:"width"`
	Height uint   `json:"height"`
	Title  string `json:"title"`
	Url    string `json:"url"`
}

func (i *Imp) ConvertToPartnerImp(price float64) partnerImp.ResponseImp {
	return partnerImp.ResponseImp{
		ID:     i.ID,
		Width:  i.Width,
		Height: i.Height,
		Title:  i.Title,
		Url:    i.Url,
		Price:  price,
	}
}
