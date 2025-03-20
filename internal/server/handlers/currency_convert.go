package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"exchanger/internal/models"
	"exchanger/internal/repository"
	"log"
	"net/http"
	"strconv"
)

type currencyConvertService interface {
	ConvertCurrency(ctx context.Context, fromCode, toCode string, amount float64) (models.CurrencyConversion, error)
}

func (h *Handlers) ExchangeCurrency(w http.ResponseWriter, r *http.Request) {
	const op = "internal.server.handlers.handlers.ExchangeCurrency"

	fromCode := r.URL.Query().Get("from")
	toCode := r.URL.Query().Get("to")
	amountStr := r.URL.Query().Get("amount")

	if fromCode == "" || toCode == "" || amountStr == "" {
		log.Printf("%s: %v", op, ErrInvalidInputData)
		errorJSON(w, "from, to, and amount parameters are required", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		log.Printf("%s: %v", op, err)
		errorJSON(w, "invalid amount format", http.StatusBadRequest)
		return
	}

	if amount <= 0 {
		log.Printf("%s: %v", op, ErrInvalidInputData)
		errorJSON(w, "amount must be greater than zero", http.StatusBadRequest)
		return
	}

	result, err := h.currencyConvertSrv.ConvertCurrency(r.Context(), fromCode, toCode, amount)
	if err != nil {
		log.Printf("%s: %v", op, err)
		if errors.Is(err, repository.ErrCurrencyNotFound) {
			errorJSON(w, "currency not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, repository.ErrExchangeRateNotFound) {
			errorJSON(w, "exchange rate not found", http.StatusNotFound)
			return
		}
		errorJSON(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
