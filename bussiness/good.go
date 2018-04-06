package bussiness

// good 商品
type good struct {
	Barcode       int64   `json:"barCode"`
	Name          string  `json:"Name"`
	InPrice       float32 `json:"InPrice"`
	OutPrice      float32 `json:"OutPrice"`
	Specification string  `json:"Specification"`
	quantity      int32   `json:"quantity"`
}
