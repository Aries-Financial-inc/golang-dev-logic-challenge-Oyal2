package server

import (
	"net/http"

	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-Oyal2/internal/analysis"
	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-Oyal2/internal/model"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	r.POST("/analyze", s.AnaylzeHandler)

	return r
}

func (s *Server) AnaylzeHandler(c *gin.Context) {
	var contracts []model.OptionsContract

	// Extract the incoming json POST request data
	if err := c.ShouldBindJSON(&contracts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Make sure that we cant have more than 4 contracts
	if len(contracts) > 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only accepting at most 4 options contracts"})
		return
	}

	// Make sure that we cant have 0 contracts
	if len(contracts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "need at least one options contracts"})
		return
	}

	for _, contract := range contracts {
		if err := model.IsOptionsContractValid(contract); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// Analyze Contracts. I am also assuming that the contracts are holding 100 share since the option size isnt mentioned.
	analysis := analysis.AnalyzeContracts(contracts)

	c.JSON(http.StatusOK, analysis)
}
