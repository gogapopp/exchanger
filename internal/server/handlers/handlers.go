package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

var (
	ErrInvalidInputData = errors.New("invalid input data")
)

type Handlers struct {
	currencySrv        currencyService
	exchangeRateSrv    exchangeRateService
	currencyConvertSrv currencyConvertService
}

func New(currencySrv currencyService, exchangeRateSrv exchangeRateService, currencyConvertSrv currencyConvertService) *Handlers {
	return &Handlers{
		currencySrv:        currencySrv,
		exchangeRateSrv:    exchangeRateSrv,
		currencyConvertSrv: currencyConvertSrv,
	}
}

func errorJSON(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})

	if err != nil {
		log.Printf("failed to encode error JSON: %v", err)
		http.Error(w, message, statusCode)
	}
}
