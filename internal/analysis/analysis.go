package analysis

import (
	"math"
	"sort"
	"strconv"

	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-Oyal2/internal/model"
)

// AnalyzeContracts performs the analysis on the given options contracts
func AnalyzeContracts(contracts []model.OptionsContract) model.Analysis {
	// Sort the contracts by strike price
	sort.Slice(contracts, func(i, j int) bool {
		return contracts[i].StrikePrice < contracts[j].StrikePrice
	})
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
	// maxProfit, maxLoss := calculateMaxProfitAndLoss(graph)
	// Calculate the break-even points
	breakEvenPoints := calculateBreakEvenPoints(contracts)
	maxLoss := calculateMaxLoss(contracts)
	var maxLossStr string
	if math.IsInf(maxLoss, 1) {
		maxLossStr = "negative infinity"
	} else {
		maxLossStr = strconv.FormatFloat(maxLoss, 'f', 2, 64)
	}

	maxProfit := calculateMaxProfit(contracts)
	var maxProfitStr string
	if math.IsInf(maxProfit, 1) {
		maxProfitStr = "infinity"
	} else {
		maxProfitStr = strconv.FormatFloat(maxProfit, 'f', 2, 64)
	}

	return model.Analysis{
		XYValues:        graph,
		MaxProfit:       maxProfitStr,
		MaxLoss:         maxLossStr,
		BreakEvenPoints: breakEvenPoints,
	}
}

func determinePriceRange(contracts []model.OptionsContract) (float64, float64) {
	minStrike, maxStrike := contracts[0].StrikePrice, contracts[len(contracts)-1].StrikePrice
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
func calculateMaxProfit(contracts []model.OptionsContract) float64 {
	maxProfit := math.Inf(-1)
	entryCredit := calculateEntryPoint(contracts)

	var pieceWiseRange []model.PieceWiseRange
	pieceWiseRange = append(pieceWiseRange, model.PieceWiseRange{A: math.Inf(-1), B: contracts[0].StrikePrice})
	for i := 1; i < len(contracts); i++ {
		pieceWiseRange = append(pieceWiseRange, model.PieceWiseRange{A: contracts[i-1].StrikePrice, B: contracts[i].StrikePrice})
	}
	pieceWiseRange = append(pieceWiseRange, model.PieceWiseRange{A: contracts[len(contracts)-1].StrikePrice, B: math.Inf(1)})

	for _, pieceWiseFunction := range pieceWiseRange {
		var testPoint float64
		if math.IsInf(pieceWiseFunction.A, -1) {
			testPoint = pieceWiseFunction.A
		} else if math.IsInf(pieceWiseFunction.B, 1) {
			testPoint = pieceWiseFunction.B
		} else {
			testPoint = (pieceWiseFunction.A + pieceWiseFunction.B) / 2
		}

		totalRisk := calculateTotalRisk(testPoint, contracts)
		profitLoss := totalRisk - entryCredit

		// Update maxProfit if the current profit is higher
		if profitLoss > maxProfit {
			maxProfit = profitLoss
		}

	}

	return maxProfit
}

// calculateMaxProfitAndLoss calculates the maximum profit and maximum loss from the graph
func calculateMaxLoss(contracts []model.OptionsContract) float64 {
	maxLoss := math.Inf(1)
	entryCredit := calculateEntryPoint(contracts) * 100

	var pieceWiseRange []model.PieceWiseRange
	pieceWiseRange = append(pieceWiseRange, model.PieceWiseRange{A: math.Inf(-1), B: contracts[0].StrikePrice})
	for i := 1; i < len(contracts); i++ {
		pieceWiseRange = append(pieceWiseRange, model.PieceWiseRange{A: contracts[i-1].StrikePrice, B: contracts[i].StrikePrice})
	}
	pieceWiseRange = append(pieceWiseRange, model.PieceWiseRange{A: contracts[len(contracts)-1].StrikePrice, B: math.Inf(1)})

	for _, pieceWiseFunction := range pieceWiseRange {
		var testPoint float64
		if math.IsInf(pieceWiseFunction.A, -1) {
			testPoint = pieceWiseFunction.A
		} else if math.IsInf(pieceWiseFunction.B, 1) {
			testPoint = pieceWiseFunction.B
		} else {
			testPoint = (pieceWiseFunction.A + pieceWiseFunction.B) / 2
		}

		totalRisk := calculateTotalRisk(testPoint, contracts)
		profitLoss := totalRisk - entryCredit

		// Update maxProfit if the current profit is higher
		if profitLoss < maxLoss {
			maxLoss = profitLoss
		}

	}

	return maxLoss
}

func calculateTotalRisk(price float64, contracts []model.OptionsContract) float64 {
	risk := 0.0
	// Sum up the profit for each contract at the given price
	for _, contract := range contracts {
		switch contract.Type {
		case model.Call:
			if contract.LongShort == model.Long {
				risk += callRisk(price, contract.StrikePrice)
			} else {
				risk += (-callRisk(price, contract.StrikePrice))
			}
		case model.Put:
			if contract.LongShort == model.Long {
				risk += putRisk(price, contract.StrikePrice)
			} else {
				risk += (-callRisk(price, contract.StrikePrice))
			}
		}
	}
	// Round off the profit
	return round(risk * 100)
}

func calculateBreakEvenPoints(contracts []model.OptionsContract) []float64 {
	minStrike, maxStrike := contracts[0].StrikePrice, contracts[len(contracts)-1].StrikePrice
	var solutions []float64
	entryCredit := calculateEntryPoint(contracts)
	n := len(contracts)
	if n == 0 {
		return solutions
	}

	// Case 1: x <= min(strikes)
	sumPut := 0.0
	countPut := 0
	for _, contract := range contracts {
		if contract.Type == model.Put && contract.LongShort == model.Long {
			sumPut += contract.StrikePrice
			countPut++
		} else if contract.Type == model.Put && contract.LongShort == model.Short {
			sumPut -= contract.StrikePrice
			countPut--
		}
	}
	if countPut > 0 {
		x := (sumPut - entryCredit) / float64(countPut)
		if x <= minStrike {
			solutions = append(solutions, round(x))
		}
	}

	// Intermediate cases (combinations of call and put)
	for mask := 1; mask < (1 << n); mask++ {
		sum := 0.0
		count := 0
		valid := true
		for i := 0; i < n; i++ {
			if (mask & (1 << i)) != 0 {
				if contracts[i].Type == model.Call && contracts[i].LongShort == model.Long {
					sum += callRisk(sum, contracts[i].StrikePrice)
				} else {
					sum += (-callRisk(sum, contracts[i].StrikePrice))
				}
				count++
			} else {
				if contracts[i].LongShort == model.Long {
					sum += putRisk(sum, contracts[i].StrikePrice)
					count++
				} else {
					sum += (-putRisk(sum, contracts[i].StrikePrice))
					count++
				}
			}
		}
		if valid {
			x := (sum - entryCredit) / float64(count)
			if x > minStrike && x <= maxStrike {
				solutions = append(solutions, round(x))
			}
		}
	}

	// Case 3: x > max(strikes)
	sumCall := 0.0
	countCall := 0
	for _, contract := range contracts {
		if contract.Type == model.Call && contract.LongShort == model.Long {
			sumCall += contract.StrikePrice
			countCall++
		} else if contract.Type == model.Call && contract.LongShort == model.Short {
			sumCall -= contract.StrikePrice
			countCall--
		}
	}
	if countCall > 0 {
		x := (sumCall + entryCredit) / float64(countCall)
		if x > maxStrike {
			solutions = append(solutions, round(x))
		}
	}

	return solutions
}

func calculateEntryPoint(contracts []model.OptionsContract) float64 {
	profit := 0.0
	for _, contract := range contracts {
		switch contract.Type {
		case model.Call:
			if contract.LongShort == model.Long {
				profit += contract.Ask
			} else {
				profit -= contract.Ask
			}
		case model.Put:
			if contract.LongShort == model.Long {
				profit += contract.Bid
			} else {
				profit -= contract.Bid
			}
		}
	}

	return round(profit)
}

func callRisk(price, strike float64) float64 {
	return math.Max(0, price-strike)
}

func putRisk(price, strike float64) float64 {
	return math.Max(0, strike-price)
}

// round rounds a float to two decimal places
func round(x float64) float64 {
	return math.Ceil(x*100) / 100
}
