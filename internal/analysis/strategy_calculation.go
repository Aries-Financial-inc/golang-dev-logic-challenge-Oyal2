package analysis

import (
	"math"

	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-Oyal2/internal/model"
)

// CalculateMaxLossAndProfit calculates the maximum profit and minimum loss for a set of options contracts
func CalculateMaxLossAndProfit(contracts []model.OptionsContract) (float64, float64) {
	// Calculate the entry price using a helper function
	entryPrice := CalculateEntryPoint(contracts)

	// Initialize maxProfit to negative infinity and minLoss to positive infinity
	maxProfit := math.Inf(-1)
	minLoss := math.Inf(1)
	infinitePos := false
	infiniteNeg := false

	// Iterate over all contracts to find max profit and min loss at strike prices
	for _, contract := range contracts {
		profitLoss := CalculateProfitLoss(contract.StrikePrice, entryPrice, contracts)
		if profitLoss > maxProfit {
			maxProfit = profitLoss
		}
		if profitLoss < minLoss {
			minLoss = profitLoss
		}
	}

	// Check for intermediate max profit and min loss between strike prices
	for i := 0; i < len(contracts)-2; i++ {
		price1 := contracts[i].StrikePrice
		price2 := contracts[i+1].StrikePrice
		midPoint := (price1 + price2) / 2
		profitLoss := CalculateProfitLoss(midPoint, entryPrice, contracts)
		if profitLoss > maxProfit {
			maxProfit = profitLoss
		}
		if profitLoss < minLoss {
			minLoss = profitLoss
		}
	}

	// Check for potential infinite profit or loss
	largeStrike := contracts[len(contracts)-1].StrikePrice + 1000
	profitLoss := CalculateProfitLoss(largeStrike, entryPrice, contracts)
	if profitLoss > maxProfit {
		infinitePos = true
	}
	if profitLoss < minLoss {
		infiniteNeg = true
	}

	smallStrike := contracts[len(contracts)-1].StrikePrice - 1000
	profitLoss = CalculateProfitLoss(smallStrike, entryPrice, contracts)
	if profitLoss > maxProfit {
		infinitePos = true
	}
	if profitLoss < minLoss {
		infiniteNeg = true
	}

	// Set maxProfit and minLoss to infinity if indicated
	if infinitePos {
		maxProfit = math.Inf(1)
	}
	if infiniteNeg {
		minLoss = math.Inf(-1)
	}

	// Return the maximum profit and minimum loss
	return maxProfit, minLoss
}

// CalculateTotalProfit calculates the total profit for a set of options contracts at a given price
func CalculateTotalProfit(contracts []model.OptionsContract, price float64) float64 {
	profit := 0.0
	// Sum up the profit for each contract at the given price
	for _, contract := range contracts {
		switch contract.Type {
		case model.Call:
			profit += (CalculateCallProfit(contract, price))
		case model.Put:
			profit += (CalculatePutProfit(contract, price))
		}
	}

	// Round off the profit to the nearest tenth
	return profit
}

// CalculateCallProfit calculates the profit for a call option at a given price
func CalculateCallProfit(contract model.OptionsContract, price float64) float64 {
	// For a long position, profit is the difference between the price and strike price minus the ask
	if contract.LongShort == model.Long {
		return max(0, price-contract.StrikePrice) - contract.Ask
	}
	// For a short position, profit is the bid minus the difference between the price and strike price
	return contract.Bid - max(0, price-contract.StrikePrice)
}

// CalculatePutProfit calculates the profit for a put option at a given price
func CalculatePutProfit(contract model.OptionsContract, price float64) float64 {
	// For a long position, profit is the difference between the strike price and price minus the ask
	if contract.LongShort == model.Long {
		return max(0, contract.StrikePrice-price) - contract.Ask
	}
	// For a short position, profit is the bid minus the difference between the strike price and price
	return contract.Bid - max(0, contract.StrikePrice-price)
}

// CalculateProfitLoss calculates the profit or loss for a set of options contracts at a given price
func CalculateProfitLoss(price float64, entryPrice float64, contracts []model.OptionsContract) float64 {
	profit_loss := -entryPrice // Initialize profit/loss with negative entry price
	for _, contract := range contracts {
		switch contract.Type {
		case model.Call:
			if contract.LongShort == model.Long {
				profit_loss += CallRisk(price, contract.StrikePrice)
			} else {
				profit_loss += (-CallRisk(price, contract.StrikePrice))
			}
		case model.Put:
			if contract.LongShort == model.Long {
				profit_loss += PutRisk(price, contract.StrikePrice)
			} else {
				profit_loss += (-PutRisk(price, contract.StrikePrice))
			}
		}
	}

	// Return the calculated profit or loss
	return profit_loss
}
