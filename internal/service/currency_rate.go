package service

import (
	"context"
	"exchanger/internal/models"
)

type exchangeRateService struct {
	exchangeRateRepo exchangeRateRepository
}

func NewExchangeRateService(exchangeRateRepo exchangeRateRepository) *exchangeRateService {
	return &exchangeRateService{
		exchangeRateRepo: exchangeRateRepo,
	}
}

type exchangeRateRepository interface {
	GetAllExchangeRates(ctx context.Context) ([]models.ExchangeRate, error)
	GetExchangeRate(ctx context.Context, baseCode, targetCode string) (models.ExchangeRate, error)
	AddExchangeRate(ctx context.Context, baseCode, targetCode string, rate float64) (models.ExchangeRate, error)
	UpdateExchangeRate(ctx context.Context, baseCode, targetCode string, rate float64) (models.ExchangeRate, error)
}

func (s *exchangeRateService) GetAllExchangeRates(ctx context.Context) ([]models.ExchangeRate, error) {
	return s.exchangeRateRepo.GetAllExchangeRates(ctx)
}

func (s *exchangeRateService) GetExchangeRate(ctx context.Context, baseCode, targetCode string) (models.ExchangeRate, error) {
	return s.exchangeRateRepo.GetExchangeRate(ctx, baseCode, targetCode)
}

func (s *exchangeRateService) AddExchangeRate(ctx context.Context, baseCode, targetCode string, rate float64) (models.ExchangeRate, error) {
	return s.exchangeRateRepo.AddExchangeRate(ctx, baseCode, targetCode, rate)
}

func (s *exchangeRateService) UpdateExchangeRate(ctx context.Context, baseCode, targetCode string, rate float64) (models.ExchangeRate, error) {
	return s.exchangeRateRepo.UpdateExchangeRate(ctx, baseCode, targetCode, rate)
}
