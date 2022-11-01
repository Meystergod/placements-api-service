package imp

type RequestImp struct {
	ID        uint `json:"id"`
	MinWidth  uint `json:"minwidth"`
	MinHeight uint `json:"minheight"`
}

type ResponseImp struct {
	ID     uint    `json:"id"`
	Width  uint    `json:"width"`
	Height uint    `json:"height"`
	Title  string  `json:"title"`
	Url    string  `json:"url"`
	Price  float64 `json:"price"`
}
