package unit

import (
	"time"

	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-Oyal2/internal/analysis"
	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-Oyal2/internal/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("analysis.CalculateBreakEvenPoints", func() {
	Context("with various options contracts", func() {
		It("should calculate break-even points for single call option", func() {
			contracts := []model.OptionsContract{
				{
					Type:           model.Call,
					LongShort:      model.Long,
					StrikePrice:    100.0,
					Bid:            10.0,
					Ask:            12.0,
					ExpirationDate: time.Now().AddDate(0, 1, 0),
				},
			}

			breakEvenPoints := analysis.CalculateBreakEvenPoints(contracts)
			Expect(breakEvenPoints).To(Equal([]float64{112.0})) // The break-even point for a long call option is Strike Price + Ask
		})

		It("should calculate break-even points for single put option", func() {
			contracts := []model.OptionsContract{
				{
					Type:           model.Put,
					LongShort:      model.Long,
					StrikePrice:    90.0,
					Bid:            8.0,
					Ask:            10.0,
					ExpirationDate: time.Now().AddDate(0, 1, 0),
				},
			}

			breakEvenPoints := analysis.CalculateBreakEvenPoints(contracts)
			Expect(breakEvenPoints).To(Equal([]float64{82.0}))
		})

		It("should calculate break-even points for multiple options contracts", func() {
			contracts := []model.OptionsContract{
				{
					Type:           model.Call,
					LongShort:      model.Long,
					StrikePrice:    100.0,
					Bid:            10.0,
					Ask:            12.0,
					ExpirationDate: time.Now().AddDate(0, 1, 0),
				},
				{
					Type:           model.Put,
					LongShort:      model.Short,
					StrikePrice:    95.0,
					Bid:            7.0,
					Ask:            9.0,
					ExpirationDate: time.Now().AddDate(0, 1, 0),
				},
			}

			breakEvenPoints := analysis.CalculateBreakEvenPoints(contracts)
			Expect(breakEvenPoints).To(ConsistOf(105.0))
		})

		It("should calculate break-even points for edge case with zero bid and ask", func() {
			contracts := []model.OptionsContract{
				{
					Type:           model.Call,
					LongShort:      model.Long,
					StrikePrice:    100.0,
					Bid:            0.0,
					Ask:            0.0,
					ExpirationDate: time.Now().AddDate(0, 1, 0),
				},
			}

			breakEvenPoints := analysis.CalculateBreakEvenPoints(contracts)
			Expect(breakEvenPoints).To(Equal([]float64{100.0})) // The break-even point for zero bid/ask is the strike price itself
		})
	})
})
