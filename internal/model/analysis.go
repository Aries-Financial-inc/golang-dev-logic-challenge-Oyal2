package model

// Analysis represents the data structure of the analysis result
type Analysis struct {
	RiskRewardGraph []RiskRewardGraph `json:"risk_reward_graph"`
	MaxProfit       string            `json:"max_profit"`
	MaxLoss         string            `json:"max_loss"`
	BreakEvenPoints []float64         `json:"break_even_points"`
}

// RiskRewardGraph represents a pair of X and Y values
type RiskRewardGraph struct {
	UnderlyingPrice float64 `json:"underlying_price"`
	ProfitLoss      float64 `json:"profit_loss"`
}
