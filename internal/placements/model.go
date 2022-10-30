package placements

type PlacementRequest struct {
	ID      string           `json:"id" validate:"required"`
	Tiles   []PlacementTiles `json:"tiles" validate:"required,min=1,dive"`
	Context PlacementContext `json:"context" validate:"required,dive"`
}

type PlacementTiles struct {
	ID    uint    `json:"id" validate:"required"`
	Width uint    `json:"width" validate:"required,min=0"`
	Ratio float64 `json:"ratio" validate:"required,min=0"`
}

type PlacementContext struct {
	IP        string `json:"ip" validate:"required,ip4_addr"`
	UserAgent string `json:"user_agent" validate:"required"`
}

type PlacementResponse struct {
	ID  string                 `json:"id"`
	Imp []PlacementResponseImp `json:"imp"`
}

type PlacementResponseImp struct {
	ID     uint   `json:"id"`
	Width  uint   `json:"width"`
	Height uint   `json:"height"`
	Title  string `json:"title"`
	Url    string `json:"url"`
}

type PartnerRequest struct {
	ID      string              `json:"id"`
	Imp     []PartnerRequestImp `json:"imp"`
	Context PartnerContext      `json:"context"`
}

type PartnerRequestImp struct {
	ID        uint `json:"id"`
	MinWidth  uint `json:"minwidth"`
	MinHeight uint `json:"minheight"`
}

type PartnerContext struct {
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
}

type PartnerResponse struct {
	ID  string               `json:"id"`
	Imp []PartnerResponseImp `json:"imp"`
}

type PartnerResponseImp struct {
	ID     uint    `json:"id"`
	Width  uint    `json:"width"`
	Height uint    `json:"height"`
	Title  string  `json:"title"`
	Url    string  `json:"url"`
	Price  float64 `json:"price"`
}
