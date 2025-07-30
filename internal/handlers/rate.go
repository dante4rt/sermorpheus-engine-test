package handlers

import (
	"net/http"
	"sermorpheus-engine-test/internal/services"

	"github.com/gin-gonic/gin"
)

type RateHandler struct {
	rateService *services.RateService
}

func NewRateHandler(rateService *services.RateService) *RateHandler {
	return &RateHandler{
		rateService: rateService,
	}
}

func (rh *RateHandler) GetCurrentRate(c *gin.Context) {
	rate, err := rh.rateService.GetCurrentRate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get current rate",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rate_id":          rate.ID,
		"idr_to_usdt_rate": rate.IDRToUSDTRate,
		"created_at":       rate.CreatedAt,
		"source":           "live_api",
	})
}
