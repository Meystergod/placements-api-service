package partner

import (
	"sync"

	partnerContext "github.com/Meystergod/placements-api-service/internal/models/partner/context"
	partnerImp "github.com/Meystergod/placements-api-service/internal/models/partner/imp"
)

type Response struct {
	ID  string                   `json:"id"`
	Imp []partnerImp.ResponseImp `json:"imp"`
}

type Request struct {
	ID      string                  `json:"id"`
	Imp     []partnerImp.RequestImp `json:"imp"`
	Context partnerContext.Context  `json:"context"`
}

func (r *Response) ToMap(wg *sync.WaitGroup, imp map[uint][]partnerImp.ResponseImp) {
	defer wg.Done()

	for _, val := range r.Imp {
		imp[val.ID] = append(imp[val.ID], val)
	}
}
