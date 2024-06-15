package unit

import (
	"math"

	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-Oyal2/internal/analysis"
	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-Oyal2/internal/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CalculateMaxLossAndProfit", func() {
	It("should calculate the maximum profit and minimum loss correctly", func() {
		contracts := []model.OptionsContract{
			{StrikePrice: 100, Type: model.Call, Ask: 12.04, Bid: 10.05, LongShort: model.Short},
			{StrikePrice: 102.5, Type: model.Call, Ask: 14, Bid: 12.1, LongShort: model.Short},
			{StrikePrice: 103, Type: model.Put, Ask: 15.5, Bid: 14, LongShort: model.Long},
			{StrikePrice: 105, Type: model.Put, Ask: 18, Bid: 16, LongShort: model.Short},
		}

		maxProfit, minLoss := analysis.CalculateMaxLossAndProfit(contracts)
		expectedMaxProfit := 2604.0
		expectedMinLoss := math.Inf(-1)

		Expect(maxProfit).To(Equal(expectedMaxProfit))
		Expect(minLoss).To(Equal(expectedMinLoss))
	})

	It("should calculate the correct maximum profit and minimum loss for the new set of contracts", func() {
		contracts := []model.OptionsContract{
			{StrikePrice: 110, Type: model.Call, Ask: 9.5, Bid: 8.0, LongShort: model.Long},
			{StrikePrice: 115, Type: model.Call, Ask: 6.5, Bid: 5.0, LongShort: model.Short},
			{StrikePrice: 100, Type: model.Put, Ask: 5.5, Bid: 4.0, LongShort: model.Long},
			{StrikePrice: 95, Type: model.Put, Ask: 3.5, Bid: 2.5, LongShort: model.Short},
		}

		maxProfit, minLoss := analysis.CalculateMaxLossAndProfit(contracts)
		expectedMaxProfit := 50.0 // Replace with the expected value
		expectedMinLoss := -450.0 // Replace with the expected value

		Expect(maxProfit).To(BeNumerically("~", expectedMaxProfit, 1e-3))
		Expect(minLoss).To(BeNumerically("~", expectedMinLoss, 1e-3))
	})

	It("should calculate both infinite maximum profit and minimum loss", func() {
		contracts := []model.OptionsContract{
			{StrikePrice: 100, Type: model.Call, Ask: 12.0, Bid: 10.0, LongShort: model.Long},
			{StrikePrice: 95, Type: model.Put, Ask: 9.0, Bid: 7.0, LongShort: model.Short},
		}

		maxProfit, minLoss := analysis.CalculateMaxLossAndProfit(contracts)

		Expect(math.IsInf(maxProfit, 1)).To(BeTrue())
		Expect(math.IsInf(minLoss, -1)).To(BeTrue())
	})
})
