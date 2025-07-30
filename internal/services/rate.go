package services

import (
	"errors"
	"sermorpheus-engine-test/internal/models"
	"time"

	"gorm.io/gorm"
)

type RateService struct {
	db *gorm.DB
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

	newRate := &models.USDTRate{
		IDRToUSDTRate: 15420.50,
	}

	if err := rs.db.Create(newRate).Error; err != nil {
		return nil, err
	}

	return newRate, nil
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
