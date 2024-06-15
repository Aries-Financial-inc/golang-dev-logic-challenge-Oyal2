package analysis

import (
	"math"

	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-Oyal2/internal/model"
)

// CalculateBreakEvenPoints calculates the break-even points for a given set of options contracts
func CalculateBreakEvenPoints(contracts []model.OptionsContract) (breakEvenPoints []float64) {
	// Calculate the entry price
	entryPrice := CalculateEntryPoint(contracts)

	// Case 1: x <= min(strikes) Left Extremity
	sumPut := 0.0
	countPut := 0
	// Iterate over all contracts to calculate sum and count for long and short puts
	for _, contract := range contracts {
		if contract.Type == model.Put && contract.LongShort == model.Long {
			sumPut += contract.StrikePrice
			countPut++
		} else if contract.Type == model.Put && contract.LongShort == model.Short {
			sumPut -= contract.StrikePrice
			countPut--
		}
	}

	// Calculate break-even point for the put contracts
	if countPut > 0 {
		x := (sumPut - entryPrice) / float64(countPut)
		if x <= contracts[0].StrikePrice {
			breakEvenPoints = append(breakEvenPoints, round(x))
		}
	}

	// Iterate over all contracts to find intermediate break-even points
	for i := 1; i < len(contracts); i++ {
		var x1, x2 float64
		x1 = contracts[i-1].StrikePrice
		x2 = contracts[i].StrikePrice

		minProfitLoss := CalculateProfitLoss(x1, entryPrice, contracts)
		maxProfitLoss := CalculateProfitLoss(x2, entryPrice, contracts)

		// Check if either strike price is a break-even point
		if minProfitLoss == 0 {
			breakEvenPoints = append(breakEvenPoints, x1)
		}
		if maxProfitLoss == 0 {
			breakEvenPoints = append(breakEvenPoints, x2)
		}

		// Check if there is a sign change between x1 and x2 indicating a break-even point
		if minProfitLoss < 0 && 0 < maxProfitLoss || maxProfitLoss < 0 && 0 < minProfitLoss {
			breakEvenPoints = append(breakEvenPoints, round(BisectionMethod(x1, x2, entryPrice, contracts)))
		}
	}

	// Case 3: x > max(strikes) Right Extremity
	sumCall := 0.0
	countCall := 0

	// Iterate over all contracts to calculate sum and count for long and short calls
	for _, contract := range contracts {
		if contract.Type == model.Call && contract.LongShort == model.Long {
			sumCall += contract.StrikePrice
			countCall++
		} else if contract.Type == model.Call && contract.LongShort == model.Short {
			sumCall -= contract.StrikePrice
			countCall--
		}
	}

	// Calculate break-even point for the call contracts
	if countCall > 0 {
		x := (sumCall + entryPrice) / float64(countCall)
		if x > contracts[len(contracts)-1].StrikePrice {
			breakEvenPoints = append(breakEvenPoints, round(x))
		}
	}

	// Return the list of calculated break-even points
	return breakEvenPoints
}

// BisectionMethod uses the bisection method to find a root of the profit/loss function within an interval [a, b]
func BisectionMethod(a, b, entryPrice float64, contracts []model.OptionsContract) float64 {
	const tolerance = 1e-3 // Define the tolerance for stopping the iteration (.001)

	// Iterate until the interval is sufficiently small
	for (b-a)/2 > tolerance {
		midPoint := (a + b) / 2
		profitLossMid := CalculateProfitLoss(midPoint, entryPrice, contracts)

		// If the profit/loss at midPoint is within the tolerance, return midPoint
		if math.Abs(profitLossMid) <= tolerance {
			return midPoint
		}

		profitLossA := CalculateProfitLoss(a, entryPrice, contracts)

		// Determine which sub-interval contains the root
		if (profitLossMid > 0 && profitLossA < 0) || (profitLossMid < 0 && profitLossA > 0) {
			b = midPoint
		} else {
			a = midPoint
		}
	}

	// Return the midpoint of the final interval as the root approximation
	return (a + b) / 2
}
