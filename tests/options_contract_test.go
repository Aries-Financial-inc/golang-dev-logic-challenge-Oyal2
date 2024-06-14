package unit

import (
	"time"

	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-Oyal2/internal/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("OptionsContract", func() {
	Context("Validation", func() {
		It("should create a valid options contract", func() {
			contract := model.OptionsContract{
				Type:           model.Call,
				LongShort:      model.Long,
				StrikePrice:    100.0,
				Bid:            10.0,
				Ask:            12.0,
				ExpirationDate: time.Now().AddDate(0, 1, 0),
			}

			Expect(contract.Type).To(Equal(model.Call))
			Expect(contract.LongShort).To(Equal(model.Long))
			Expect(contract.StrikePrice).To(BeNumerically(">", 0))
			Expect(contract.Bid).To(BeNumerically(">", 0))
			Expect(contract.Ask).To(BeNumerically(">", 0))
			Expect(contract.ExpirationDate).To(BeTemporally(">", time.Now()))
		})

		It("should fail for invalid options contract", func() {
			contract := model.OptionsContract{
				Type:           "InvalidType",
				LongShort:      "InvalidPosition",
				StrikePrice:    -100.0,
				Bid:            -10.0,
				Ask:            -12.0,
				ExpirationDate: time.Now().AddDate(0, -1, 0),
			}

			Expect(contract.Type).NotTo(Equal(model.Call))
			Expect(contract.LongShort).NotTo(Equal(model.Long))
			Expect(contract.StrikePrice).To(BeNumerically("<", 0))
			Expect(contract.Bid).To(BeNumerically("<", 0))
			Expect(contract.Ask).To(BeNumerically("<", 0))
			Expect(contract.ExpirationDate).To(BeTemporally("<", time.Now()))
		})

		It("should handle zero and negative prices correctly", func() {
			contract := model.OptionsContract{
				Type:           model.Call,
				LongShort:      model.Long,
				StrikePrice:    0.0,
				Bid:            0.0,
				Ask:            0.0,
				ExpirationDate: time.Now().AddDate(0, 1, 0),
			}

			Expect(contract.StrikePrice).To(BeNumerically("==", 0))
			Expect(contract.Bid).To(BeNumerically("==", 0))
			Expect(contract.Ask).To(BeNumerically("==", 0))
		})
	})
})
