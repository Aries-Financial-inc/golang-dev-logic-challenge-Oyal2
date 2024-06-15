package analysis

import (
	"math"

	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-Oyal2/internal/model"
)

// round rounds a float to two decimal places
func round(x float64) float64 {
	return math.Ceil(x*100) / 100
}

// CallRisk calculates the risk for a call option given the current price and strike price
func CallRisk(price, strike float64) float64 {
	return math.Max(0, price-strike)
}

// PutRisk calculates the risk for a put option given the current price and strike price
func PutRisk(price, strike float64) float64 {
	return math.Max(0, strike-price)
}

// CalculateEntryPoint calculates the total entry point for a set of options contracts
func CalculateEntryPoint(contracts []model.OptionsContract) float64 {
	profit := 0.0

	// Iterate over each contract to sum up the entry costs
	for _, contract := range contracts {
		switch contract.Type {
		case model.Call:
			if contract.LongShort == model.Long {
				profit += contract.Ask // For long call, add the ask price
			} else {
				profit -= contract.Ask // For short call, subtract the ask price
			}
		case model.Put:
			if contract.LongShort == model.Long {
				profit += contract.Bid // For long put, add the bid price
			} else {
				profit -= contract.Bid // For short put, subtract the bid price
			}
		}
	}

	// Round the total entry point to two decimal places
	return round(profit)
}

// MultiplyBySharesAmount multiplies a given float64 value by a specified shares count
func MultiplyBySharesAmount(value, factor float64) float64 {
	return round(value * factor)
}
