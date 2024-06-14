package analysis

import (
	"math"

	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-Oyal2/internal/model"
)

// AnalyzeContracts performs the analysis on the given options contracts
func AnalyzeContracts(contracts []model.OptionsContract) model.Analysis {
	// Get the Price Range
	minPrice, maxPrice := determinePriceRange(contracts)
	priceStep := (maxPrice - minPrice) / 100
	var graph []model.XYValue
	// Go through every price and calculate the total profit at that price
	for price := minPrice; price <= maxPrice; price += priceStep {
		profit := calculateTotalProfit(contracts, price)
		graph = append(graph, model.XYValue{X: price, Y: profit})
	}

	// Calculate the maximum profit and maximum loss from the graph
	maxProfit, maxLoss := calculateMaxProfitAndLoss(graph)
	// Calculate the break-even points
	breakEvenPoints := calculateBreakEvenPoints(graph)

	return model.Analysis{
		XYValues:        graph,
		MaxProfit:       maxProfit,
		MaxLoss:         maxLoss,
		BreakEvenPoints: breakEvenPoints,
	}
}

func determinePriceRange(contracts []model.OptionsContract) (float64, float64) {
	if len(contracts) == 0 {
		return 0.0, 0.0
	}
	minStrike, maxStrike := contracts[0].StrikePrice, contracts[0].StrikePrice
	// Find the minimum and maximum strike prices among the contracts
	for _, contract := range contracts {
		if contract.StrikePrice < minStrike {
			minStrike = contract.StrikePrice
		}
		if contract.StrikePrice > maxStrike {
			maxStrike = contract.StrikePrice
		}
	}

	// Calculate a buffer for the price range
	buffer := (maxStrike + minStrike) / 2
	// Ensure the minimum price is not negative
	return math.Max(0, minStrike-buffer), maxStrike + buffer
}

func calculateTotalProfit(contracts []model.OptionsContract, price float64) float64 {
	profit := 0.0
	// Sum up the profit for each contract at the given price
	for _, contract := range contracts {
		switch contract.Type {
		case model.Call:
			profit += (calculateCallProfit(contract, price) * 100)
		case model.Put:
			profit += (calculatePutProfit(contract, price) * 100)
		}
	}
	// Round off the profit
	return round(profit)
}

// calculateCallProfit calculates the profit for a call option at a given price
func calculateCallProfit(contract model.OptionsContract, price float64) float64 {
	// For a long position, profit is the difference between the price and strike price minus the ask
	if contract.LongShort == model.Long {
		return max(0, price-contract.StrikePrice) - contract.Ask
	}
	// For a short position, profit is the bid minus the difference between the price and strike price
	return contract.Bid - max(0, price-contract.StrikePrice)
}

// calculatePutProfit calculates the profit for a put option at a given price
func calculatePutProfit(contract model.OptionsContract, price float64) float64 {
	// For a long position, profit is the difference between the strike price and price minus the ask
	if contract.LongShort == model.Long {
		return max(0, contract.StrikePrice-price) - contract.Ask
	}
	// For a short position, profit is the bid minus the difference between the strike price and price
	return contract.Bid - max(0, contract.StrikePrice-price)
}

// calculateMaxProfitAndLoss calculates the maximum profit and maximum loss from the graph
func calculateMaxProfitAndLoss(graph []model.XYValue) (float64, float64) {
	maxProfit, maxLoss := -math.MaxFloat64, math.MaxFloat64
	// Iterate through the graph to find the maximum and minimum profit
	for _, point := range graph {
		if point.Y > maxProfit {
			maxProfit = point.Y
		}
		if point.Y < maxLoss {
			maxLoss = point.Y
		}
	}
	return maxProfit, maxLoss
}

// calculateBreakEvenPoints calculates the break-even points from the graph
func calculateBreakEvenPoints(graph []model.XYValue) []float64 {
	var breakEvenPoints []float64
	// Iterate through the graph to find the points where the profit crosses zero
	for i := 1; i < len(graph); i++ {
		if (graph[i-1].Y <= 0 && graph[i].Y > 0) || (graph[i-1].Y >= 0 && graph[i].Y < 0) {
			breakEvenPoints = append(breakEvenPoints, graph[i].X)
		}
	}
	return breakEvenPoints
}

// round rounds a float to two decimal places
func round(x float64) float64 {
	return math.Ceil(x*100) / 100
}
