package model

// Analysis represents the data structure of the analysis result
type Analysis struct {
	XYValues        []XYValue `json:"xy_values"`
	MaxProfit       float64   `json:"max_profit"`
	MaxLoss         float64   `json:"max_loss"`
	BreakEvenPoints []float64 `json:"break_even_points"`
}

// XYValue represents a pair of X and Y values
type XYValue struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
