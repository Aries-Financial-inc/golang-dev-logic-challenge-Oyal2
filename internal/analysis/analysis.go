package analysis

import (
	"math"
	"sort"
	"strconv"

	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-Oyal2/internal/model"
)

const SHARES_PER_CONTRACT = 100

// AnalyzeContracts performs the analysis on the given options contracts
func AnalyzeContracts(contracts []model.OptionsContract) model.Analysis {
	// Sort the contracts by strike price
	sort.Slice(contracts, func(i, j int) bool {
		return contracts[i].StrikePrice < contracts[j].StrikePrice
	})
	// Get the Price Range
	minPrice, maxPrice := DeterminePriceRange(contracts)
	priceStep := (maxPrice - minPrice) / 30
	var riskRewardGraph []model.RiskRewardGraph
	// Go through every price and calculate the total profit at that price
	for price := minPrice; price <= maxPrice; price += priceStep {
		profit := MultiplyBySharesAmount(CalculateTotalProfit(contracts, price), SHARES_PER_CONTRACT)
		riskRewardGraph = append(riskRewardGraph, model.RiskRewardGraph{UnderlyingPrice: price, ProfitLoss: profit})
	}

	// Calculate the break-even points
	breakEvenPoints := CalculateBreakEvenPoints(contracts)

	// Calculate the maximum profit and maximum loss from the graph
	maxProfit, maxLoss := CalculateMaxLossAndProfit(contracts)
	maxLossStr := strconv.FormatFloat(maxLoss, 'f', 2, 64)
	maxProfitStr := strconv.FormatFloat(maxProfit, 'f', 2, 64)

	return model.Analysis{
		RiskRewardGraph: riskRewardGraph,
		MaxProfit:       maxProfitStr,
		MaxLoss:         maxLossStr,
		BreakEvenPoints: breakEvenPoints,
	}
}

// DeterminePriceRange calculates the price range for a set of options contracts
func DeterminePriceRange(contracts []model.OptionsContract) (float64, float64) {
	const delta = 20  // A constant delta value for the buffer
	const beta = 0.20 // A constant beta value for the buffer
	// Initialize minStrike and maxStrike using the first and last contract's strike price
	minStrike, maxStrike := contracts[0].StrikePrice, contracts[len(contracts)-1].StrikePrice
	spread := maxStrike - minStrike // Calculate the spread between the max and min strike prices

	// Calculate a buffer based on the maximum of delta and spread*beta
	buffer := math.Max(delta, spread*beta)

	// Ensure the minimum price is not negative and return the price range with buffer
	return math.Max(0, minStrike-buffer), maxStrike + buffer
}
