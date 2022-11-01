package placement

import (
	"encoding/json"
	"sync"

	"github.com/Meystergod/placements-api-service/internal/models/partner"
	partnerImp "github.com/Meystergod/placements-api-service/internal/models/partner/imp"
	placementContext "github.com/Meystergod/placements-api-service/internal/models/placement/context"
	placementImp "github.com/Meystergod/placements-api-service/internal/models/placement/imp"
	placementTile "github.com/Meystergod/placements-api-service/internal/models/placement/tile"
	"github.com/Meystergod/placements-api-service/pkg/logging"
)

type Request struct {
	ID      string                   `json:"id" validate:"required"`
	Tiles   []placementTile.Tile     `json:"tiles" validate:"required,min=1,dive"`
	Context placementContext.Context `json:"context" validate:"required,dive"`
}

type Response struct {
	ID  string             `json:"id"`
	Imp []placementImp.Imp `json:"imp"`
}

func NewResponse(logger *logging.Logger, p Request, responses []string) Response {
	partnerResponsesList := make([]partner.Response, len(responses))
	partnerResponsesSet := make(map[uint][]partnerImp.ResponseImp)

	var wg sync.WaitGroup

	logger.Info("start converting string slice of responses to map")
	for i := range responses {
		logger.Infof("try unmarshalling %d response", i+1)
		if err := json.Unmarshal([]byte(responses[i]), &partnerResponsesList[i]); err != nil {
			logger.WithError(err).Error("unmarshalling failed")
			continue
		}
		wg.Add(1)
		logger.Info("converting to map")
		go partnerResponsesList[i].ToMap(&wg, partnerResponsesSet)
	}

	wg.Wait()

	var placementResponseImpList []placementImp.Imp

	logger.Info("start find max price and creating placement response")
	for j := range p.Tiles {
		if len(partnerResponsesSet[p.Tiles[j].ID]) != 0 {
			partnerMaxValue := findMaxPrice(partnerResponsesSet[p.Tiles[j].ID])
			imp := placementImp.Imp{
				ID:     partnerMaxValue.ID,
				Width:  partnerMaxValue.Width,
				Height: partnerMaxValue.Height,
				Title:  partnerMaxValue.Title,
				Url:    partnerMaxValue.Url,
			}
			placementResponseImpList = append(placementResponseImpList, imp)
		}
	}

	return Response{
		ID:  p.ID,
		Imp: placementResponseImpList,
	}
}

func (r *Request) ToPartnerRequest() partner.Request {
	var partnerRequestImpList []partnerImp.RequestImp
	for _, tile := range r.Tiles {
		partnerRequestImpList = append(partnerRequestImpList, *tile.ConvertToPartnerRequestImp())
	}

	return partner.Request{
		ID:      r.ID,
		Imp:     partnerRequestImpList,
		Context: r.Context.ToPartnerContext(),
	}
}

func findMaxPrice(data []partnerImp.ResponseImp) partnerImp.ResponseImp {
	var resultIndex = 0
	maxVal := data[0].Price

	for i := range data {
		if data[i].Price > maxVal {
			maxVal = data[i].Price
			resultIndex = i
		}
	}

	return data[resultIndex]
}
