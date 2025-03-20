package service

import (
	"context"
	"exchanger/internal/models"
)

type currencyService struct {
	currencyRepo currencyRepository
}

func NewCurrencyService(currencyRepo currencyRepository) *currencyService {
	return &currencyService{
		currencyRepo: currencyRepo,
	}
}

type currencyRepository interface {
	GetAllCurrencies(ctx context.Context) ([]models.Currency, error)
	GetCurrencyByCode(ctx context.Context, code string) (models.Currency, error)
	AddCurrency(ctx context.Context, currency models.Currency) (models.Currency, error)
}

func (s *currencyService) GetAllCurrencies(ctx context.Context) ([]models.Currency, error) {
	return s.currencyRepo.GetAllCurrencies(ctx)
}

func (s *currencyService) GetCurrencyByCode(ctx context.Context, code string) (models.Currency, error) {
	return s.currencyRepo.GetCurrencyByCode(ctx, code)
}

func (s *currencyService) AddCurrency(ctx context.Context, currency models.Currency) (models.Currency, error) {
	return s.currencyRepo.AddCurrency(ctx, currency)
}
