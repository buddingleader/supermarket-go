package bussiness

import "strconv"

// checkBarcode check standard barcode
func checkBarcode(str string) (int64, bool) {
	if len(str) <= 8 {
		return 0, false
	}
	v, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

// checkPrice check price
func checkPrice(str string) (float64, bool) {
	if len(str) > 8 {
		return 0, false
	}
	v, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0.00, false
	}
	return v, true
}
