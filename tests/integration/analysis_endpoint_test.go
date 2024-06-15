package unit_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-Oyal2/internal/model"
	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-Oyal2/internal/server"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Helper function to read file content
func readFileContent(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

var _ = Describe("Analyze Endpoint", func() {
	var router http.Handler

	beforeEach := func() {
		server := &server.Server{}
		router = server.RegisterRoutes()
	}

	Context("POST /analyze", func() {
		It("should return analysis for valid contracts", func() {
			beforeEach()

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

			body, _ := json.Marshal(contracts)
			req, _ := http.NewRequest("POST", "/analyze", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))

			var analysis model.Analysis
			err := json.Unmarshal(w.Body.Bytes(), &analysis)
			Expect(err).To(BeNil())
			Expect(analysis.RiskRewardGraph).NotTo(BeEmpty())
			Expect(analysis.MaxProfit).To(Equal("+Inf"))
			Expect(analysis.MaxLoss).To(Equal("-1200.00"))
		})

		It("should return error for more than 4 contracts", func() {
			beforeEach()

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
					StrikePrice:    110.0,
					Bid:            5.0,
					Ask:            6.0,
					ExpirationDate: time.Now().AddDate(0, 2, 0),
				},
				{
					Type:           model.Call,
					LongShort:      model.Long,
					StrikePrice:    120.0,
					Bid:            8.0,
					Ask:            9.0,
					ExpirationDate: time.Now().AddDate(0, 3, 0),
				},
				{
					Type:           model.Put,
					LongShort:      model.Short,
					StrikePrice:    130.0,
					Bid:            7.0,
					Ask:            8.0,
					ExpirationDate: time.Now().AddDate(0, 4, 0),
				},
				{
					Type:           model.Call,
					LongShort:      model.Long,
					StrikePrice:    140.0,
					Bid:            6.0,
					Ask:            7.0,
					ExpirationDate: time.Now().AddDate(0, 5, 0),
				},
			}

			body, _ := json.Marshal(contracts)
			req, _ := http.NewRequest("POST", "/analyze", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("only accepting at most 4 options contracts"))
		})

		It("should return error for no contracts", func() {
			beforeEach()

			contracts := []model.OptionsContract{}

			body, _ := json.Marshal(contracts)
			req, _ := http.NewRequest("POST", "/analyze", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("need at least one options contracts"))
		})

		It("should return error for invalid JSON", func() {
			beforeEach()

			invalidJSON := `{ invalid json }`
			req, _ := http.NewRequest("POST", "/analyze", bytes.NewBufferString(invalidJSON))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("invalid character"))
		})

		It("should return error for missing fields in contract", func() {
			beforeEach()

			invalidContract := []map[string]interface{}{
				{
					"type": "Call",
				},
			}

			body, _ := json.Marshal(invalidContract)
			req, _ := http.NewRequest("POST", "/analyze", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("error"))
		})

		It("should return analysis with break-even points", func() {
			beforeEach()

			contracts := []model.OptionsContract{
				{
					Type:           model.Put,
					LongShort:      model.Short,
					StrikePrice:    100.0,
					Bid:            5.0,
					Ask:            6.0,
					ExpirationDate: time.Now().AddDate(0, 1, 0),
				},
			}

			body, _ := json.Marshal(contracts)
			req, _ := http.NewRequest("POST", "/analyze", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))

			var analysis model.Analysis
			err := json.Unmarshal(w.Body.Bytes(), &analysis)
			Expect(err).To(BeNil())
			Expect(analysis.BreakEvenPoints).NotTo(Equal(106))
		})

		It("should return error for invalid option type", func() {
			beforeEach()

			contracts := []model.OptionsContract{
				{
					Type:           "InvalidType",
					LongShort:      model.Long,
					StrikePrice:    100.0,
					Bid:            10.0,
					Ask:            12.0,
					ExpirationDate: time.Now().AddDate(0, 1, 0),
				},
			}

			body, _ := json.Marshal(contracts)
			req, _ := http.NewRequest("POST", "/analyze", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("invalid option type"))
		})

		It("should return error for invalid position type", func() {
			beforeEach()

			contracts := []model.OptionsContract{
				{
					Type:           model.Call,
					LongShort:      "InvalidPosition",
					StrikePrice:    100.0,
					Bid:            10.0,
					Ask:            12.0,
					ExpirationDate: time.Now().AddDate(0, 1, 0),
				},
			}

			body, _ := json.Marshal(contracts)
			req, _ := http.NewRequest("POST", "/analyze", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("invalid position type"))
		})

		It("should return error for non-positive strike price", func() {
			beforeEach()

			contracts := []model.OptionsContract{
				{
					Type:           model.Call,
					LongShort:      model.Long,
					StrikePrice:    -100.0,
					Bid:            10.0,
					Ask:            12.0,
					ExpirationDate: time.Now().AddDate(0, 1, 0),
				},
			}

			body, _ := json.Marshal(contracts)
			req, _ := http.NewRequest("POST", "/analyze", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("strike price must be greater than zero"))
		})

		It("should return error for negative bid", func() {
			beforeEach()

			contracts := []model.OptionsContract{
				{
					Type:           model.Call,
					LongShort:      model.Long,
					StrikePrice:    100.0,
					Bid:            -10.0,
					Ask:            12.0,
					ExpirationDate: time.Now().AddDate(0, 1, 0),
				},
			}

			body, _ := json.Marshal(contracts)
			req, _ := http.NewRequest("POST", "/analyze", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("bid must be non-negative"))
		})

		It("should return error for negative ask", func() {
			beforeEach()

			contracts := []model.OptionsContract{
				{
					Type:           model.Call,
					LongShort:      model.Long,
					StrikePrice:    100.0,
					Bid:            10.0,
					Ask:            -12.0,
					ExpirationDate: time.Now().AddDate(0, 1, 0),
				},
			}

			body, _ := json.Marshal(contracts)
			req, _ := http.NewRequest("POST", "/analyze", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("ask must be non-negative"))
		})

		It("should return error for past expiration date", func() {
			beforeEach()

			contracts := []model.OptionsContract{
				{
					Type:           model.Call,
					LongShort:      model.Long,
					StrikePrice:    100.0,
					Bid:            10.0,
					Ask:            12.0,
					ExpirationDate: time.Now().AddDate(0, -1, 0),
				},
			}

			body, _ := json.Marshal(contracts)
			req, _ := http.NewRequest("POST", "/analyze", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("expiration date must be in the future"))
		})

		It("should return analysis with 2 break even points", func() {
			beforeEach()

			inputData, err := readFileContent("../../testdata/multiple_break_even_95.51_114.5.json")
			Expect(err).To(BeNil())

			req, _ := http.NewRequest("POST", "/analyze", bytes.NewBuffer(inputData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))

			var analysis model.Analysis
			err = json.Unmarshal(w.Body.Bytes(), &analysis)
			Expect(err).To(BeNil())
			Expect(analysis.RiskRewardGraph).NotTo(BeEmpty())

			expectedGraph := []model.RiskRewardGraph{
				{UnderlyingPrice: 75, ProfitLoss: -250},
				{UnderlyingPrice: 77, ProfitLoss: -250},
				{UnderlyingPrice: 79, ProfitLoss: -250},
				{UnderlyingPrice: 81, ProfitLoss: -250},
				{UnderlyingPrice: 83, ProfitLoss: -250},
				{UnderlyingPrice: 85, ProfitLoss: -250},
				{UnderlyingPrice: 87, ProfitLoss: -250},
				{UnderlyingPrice: 89, ProfitLoss: -250},
				{UnderlyingPrice: 91, ProfitLoss: -250},
				{UnderlyingPrice: 93, ProfitLoss: -250},
				{UnderlyingPrice: 95, ProfitLoss: -250},
				{UnderlyingPrice: 97, ProfitLoss: -450},
				{UnderlyingPrice: 99, ProfitLoss: -650},
				{UnderlyingPrice: 101, ProfitLoss: -750},
				{UnderlyingPrice: 103, ProfitLoss: -750},
				{UnderlyingPrice: 105, ProfitLoss: -750},
				{UnderlyingPrice: 107, ProfitLoss: -750},
				{UnderlyingPrice: 109, ProfitLoss: -750},
				{UnderlyingPrice: 111, ProfitLoss: -650},
				{UnderlyingPrice: 113, ProfitLoss: -450},
				{UnderlyingPrice: 115, ProfitLoss: -250},
				{UnderlyingPrice: 117, ProfitLoss: -250},
				{UnderlyingPrice: 119, ProfitLoss: -250},
				{UnderlyingPrice: 121, ProfitLoss: -250},
				{UnderlyingPrice: 123, ProfitLoss: -250},
				{UnderlyingPrice: 125, ProfitLoss: -250},
				{UnderlyingPrice: 127, ProfitLoss: -250},
				{UnderlyingPrice: 129, ProfitLoss: -250},
				{UnderlyingPrice: 131, ProfitLoss: -250},
				{UnderlyingPrice: 133, ProfitLoss: -250},
				{UnderlyingPrice: 135, ProfitLoss: -250},
			}
			Expect(analysis.RiskRewardGraph).To(Equal(expectedGraph))
			Expect(analysis.MaxProfit).To(Equal("50.00"))
			Expect(analysis.MaxLoss).To(Equal("-450.00"))

			expectedBreakEvenPoints := []float64{95.51, 114.5}
			Expect(analysis.BreakEvenPoints).To(Equal(expectedBreakEvenPoints))
		})

		It("should return analysis for 2 leg options", func() {
			beforeEach()

			inputData, err := readFileContent("../../testdata/2leg.json")
			Expect(err).To(BeNil())

			req, _ := http.NewRequest("POST", "/analyze", bytes.NewBuffer(inputData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))

			var analysis model.Analysis
			err = json.Unmarshal(w.Body.Bytes(), &analysis)
			Expect(err).To(BeNil())
			Expect(analysis.RiskRewardGraph).NotTo(BeEmpty())

			expectedGraph := []model.RiskRewardGraph{
				{UnderlyingPrice: 80, ProfitLoss: -2604},
				{UnderlyingPrice: 81.41666666666667, ProfitLoss: -2604},
				{UnderlyingPrice: 82.83333333333334, ProfitLoss: -2604},
				{UnderlyingPrice: 84.25000000000001, ProfitLoss: -2604},
				{UnderlyingPrice: 85.66666666666669, ProfitLoss: -2604},
				{UnderlyingPrice: 87.08333333333336, ProfitLoss: -2604},
				{UnderlyingPrice: 88.50000000000003, ProfitLoss: -2604},
				{UnderlyingPrice: 89.9166666666667, ProfitLoss: -2604},
				{UnderlyingPrice: 91.33333333333337, ProfitLoss: -2604},
				{UnderlyingPrice: 92.75000000000004, ProfitLoss: -2604},
				{UnderlyingPrice: 94.16666666666671, ProfitLoss: -2604},
				{UnderlyingPrice: 95.58333333333339, ProfitLoss: -2604},
				{UnderlyingPrice: 97.00000000000006, ProfitLoss: -2604},
				{UnderlyingPrice: 98.41666666666673, ProfitLoss: -2604},
				{UnderlyingPrice: 99.8333333333334, ProfitLoss: -2604},
				{UnderlyingPrice: 101.25000000000007, ProfitLoss: -2478.99},
				{UnderlyingPrice: 102.66666666666674, ProfitLoss: -2320.66},
				{UnderlyingPrice: 104.08333333333341, ProfitLoss: -2037.33},
				{UnderlyingPrice: 105.50000000000009, ProfitLoss: -1753.99},
				{UnderlyingPrice: 106.91666666666676, ProfitLoss: -1470.66},
				{UnderlyingPrice: 108.33333333333343, ProfitLoss: -1187.33},
				{UnderlyingPrice: 109.7500000000001, ProfitLoss: -903.99},
				{UnderlyingPrice: 111.16666666666677, ProfitLoss: -620.66},
				{UnderlyingPrice: 112.58333333333344, ProfitLoss: -337.33},
				{UnderlyingPrice: 114.00000000000011, ProfitLoss: -53.99},
				{UnderlyingPrice: 115.41666666666679, ProfitLoss: 229.34},
				{UnderlyingPrice: 116.83333333333346, ProfitLoss: 512.67},
				{UnderlyingPrice: 118.25000000000013, ProfitLoss: 796.01},
				{UnderlyingPrice: 119.6666666666668, ProfitLoss: 1079.34},
				{UnderlyingPrice: 121.08333333333347, ProfitLoss: 1362.67},
			}
			Expect(analysis.RiskRewardGraph).To(Equal(expectedGraph))
			Expect(analysis.MaxProfit).To(Equal("+Inf"))
			Expect(analysis.MaxLoss).To(Equal("-2604.00"))

			expectedBreakEvenPoints := []float64{114.27}
			Expect(analysis.BreakEvenPoints).To(Equal(expectedBreakEvenPoints))
		})
	})
})
