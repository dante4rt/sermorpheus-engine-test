package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"sermorpheus-engine-test/internal/models"

	"gorm.io/gorm"
)

type RateService struct {
	db *gorm.DB
}

type ExchangeRateResponse struct {
	Rates map[string]float64 `json:"rates"`
}

func NewRateService(db *gorm.DB) *RateService {
	return &RateService{db: db}
}

func (rs *RateService) GetCurrentRate() (*models.USDTRate, error) {
	var rate models.USDTRate

	fiveMinutesAgo := time.Now().Add(-5 * time.Minute)
	err := rs.db.Where("created_at > ?", fiveMinutesAgo).
		Order("created_at DESC").
		First(&rate).Error

	if err == nil {
		return &rate, nil
	}

	usdToIdr, err := rs.fetchUSDToIDR()
	if err != nil {
		usdToIdr = 15420.50
	}

	newRate := &models.USDTRate{
		IDRToUSDTRate: usdToIdr,
	}

	if err := rs.db.Create(newRate).Error; err != nil {
		return nil, err
	}

	return newRate, nil
}

func (rs *RateService) fetchUSDToIDR() (float64, error) {
	resp, err := http.Get("https://api.exchangerate-api.com/v4/latest/USD")
	if err != nil {
		return 0, fmt.Errorf("failed to fetch exchange rate: %w", err)
	}
	defer resp.Body.Close()

	var exchangeResp ExchangeRateResponse
	if err := json.NewDecoder(resp.Body).Decode(&exchangeResp); err != nil {
		return 0, fmt.Errorf("failed to decode exchange rate response: %w", err)
	}

	idrRate, exists := exchangeResp.Rates["IDR"]
	if !exists {
		return 0, errors.New("IDR rate not found in response")
	}

	return idrRate, nil
}

func (rs *RateService) CreateRate(idrToUSDT float64) (*models.USDTRate, error) {
	if idrToUSDT <= 0 {
		return nil, errors.New("invalid exchange rate")
	}

	rate := &models.USDTRate{
		IDRToUSDTRate: idrToUSDT,
	}

	if err := rs.db.Create(rate).Error; err != nil {
		return nil, err
	}

	return rate, nil
}
