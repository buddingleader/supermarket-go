package good

import "container/heap"

// Purchase goods from trafficker
type Purchase struct {
	Barcode    string  `bson:"barcode"`
	Trafficker string  `bson:"trafficker"`
	Price      float64 `bson:"price"`
	Date       string  `bson:"date"`
	Number     uint    `bson:"number"`
}

// Sell goods to customer
type Sell struct {
	Barcode  string  `bson:"barcode"`
	OutPrice float64 `bson:"outprice"`
	Date     string  `bson:"date"`
	Times    uint    `bson:"times"`
}

// SellMaxHeap The most recent sell record is at the top of the heap
type SellMaxHeap []*Sell

// Len the heap length
func (smh SellMaxHeap) Len() int { return len(smh) }

// Less we want Pop to give us the highest, not lowest, priority so we use greater than here.
func (smh SellMaxHeap) Less(i, j int) bool {
	if smh[i].Times == smh[j].Times {
		return smh[i].Date > smh[j].Date
	}
	return smh[i].Times > smh[j].Times
}

// Swap the item location
func (smh SellMaxHeap) Swap(i, j int) {
	smh[i], smh[j] = smh[j], smh[i]
}

// Push a new item to heap
func (smh *SellMaxHeap) Push(x interface{}) {
	item := x.(*Sell)
	*smh = append(*smh, item)
}

// Pop the hightest item
func (smh *SellMaxHeap) Pop() interface{} {
	old := *smh
	n := len(old)
	item := old[n-1]
	*smh = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (smh *SellMaxHeap) update(sell *Sell, count uint, index int) {
	sell.Times += count
	heap.Fix(smh, index)
}
