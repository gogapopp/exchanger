package service

import (
	"context"
	"exchanger/internal/models"
	"exchanger/internal/repository"
	"fmt"
)

type convertService struct {
	currencyRepo     currencyRepository
	exchangeRateRepo exchangeRateRepository
}

func NewConvertService(currencyRepo currencyRepository, exchangeRateRepo exchangeRateRepository) *convertService {
	return &convertService{
		currencyRepo:     currencyRepo,
		exchangeRateRepo: exchangeRateRepo,
	}
}

const usdCode = "USD"

func (s *convertService) ConvertCurrency(ctx context.Context, fromCode, toCode string, amount float64) (models.CurrencyConversion, error) {
	const op = "internal.service.service.ConvertCurrency"

	baseCurrency, err := s.currencyRepo.GetCurrencyByCode(ctx, fromCode)
	if err != nil {
		return models.CurrencyConversion{}, fmt.Errorf("%s: %w", op, err)
	}

	targetCurrency, err := s.currencyRepo.GetCurrencyByCode(ctx, toCode)
	if err != nil {
		return models.CurrencyConversion{}, fmt.Errorf("%s: %w", op, err)
	}

	// in ExchangeRates we have currency pair AB
	exchangeRate, err := s.exchangeRateRepo.GetExchangeRate(ctx, fromCode, toCode)
	if err == nil {
		convertedAmount := amount * exchangeRate.Rate

		return models.CurrencyConversion{
			BaseCurrency:    baseCurrency,
			TargetCurrency:  targetCurrency,
			Rate:            exchangeRate.Rate,
			Amount:          amount,
			ConvertedAmount: convertedAmount,
		}, nil
	}

	// in ExchangeRates we have currency pair BA
	reverseRate, err := s.exchangeRateRepo.GetExchangeRate(ctx, toCode, fromCode)
	if err == nil {
		rate := 1 / reverseRate.Rate
		convertedAmount := amount * rate

		return models.CurrencyConversion{
			BaseCurrency:    baseCurrency,
			TargetCurrency:  targetCurrency,
			Rate:            rate,
			Amount:          amount,
			ConvertedAmount: convertedAmount,
		}, nil
	}

	// in ExchangeRates we have currency pairs USD-A and USD-B
	usdToBase, errBase := s.exchangeRateRepo.GetExchangeRate(ctx, usdCode, fromCode)
	usdToTarget, errTarget := s.exchangeRateRepo.GetExchangeRate(ctx, usdCode, toCode)

	if errBase == nil && errTarget == nil {
		rate := usdToTarget.Rate / usdToBase.Rate
		convertedAmount := amount * rate

		return models.CurrencyConversion{
			BaseCurrency:    baseCurrency,
			TargetCurrency:  targetCurrency,
			Rate:            rate,
			Amount:          amount,
			ConvertedAmount: convertedAmount,
		}, nil
	}

	return models.CurrencyConversion{}, fmt.Errorf("%s: %w", op, repository.ErrExchangeRateNotFound)
}
