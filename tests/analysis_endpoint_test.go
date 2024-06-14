package unit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-Oyal2/internal/model"
	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-Oyal2/internal/server"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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
			Expect(analysis.XYValues).NotTo(BeEmpty())
			Expect(analysis.MaxProfit).To(BeNumerically(">", 0))
			Expect(analysis.MaxLoss).To(BeNumerically("<", 0))
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
			Expect(analysis.BreakEvenPoints).NotTo(BeEmpty())
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
	})
})
