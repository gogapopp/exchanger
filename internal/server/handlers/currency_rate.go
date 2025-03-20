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

type exchangeRateService interface {
	GetAllExchangeRates(ctx context.Context) ([]models.ExchangeRate, error)
	GetExchangeRate(ctx context.Context, baseCode, targetCode string) (models.ExchangeRate, error)
	AddExchangeRate(ctx context.Context, baseCode, targetCode string, rate float64) (models.ExchangeRate, error)
	UpdateExchangeRate(ctx context.Context, baseCode, targetCode string, rate float64) (models.ExchangeRate, error)
}

func (h *Handlers) GetExchangeRates(w http.ResponseWriter, r *http.Request) {
	const op = "internal.server.handlers.handlers.GetExchangeRates"

	rates, err := h.exchangeRateSrv.GetAllExchangeRates(r.Context())
	if err != nil {
		log.Printf("%s: %v", op, err)
		if errors.Is(err, repository.ErrExchangeRateNotFound) {
			errorJSON(w, "exchange rates not found", http.StatusNotFound)
			return
		}
		errorJSON(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rates)
}

func (h *Handlers) GetExchangeRate(w http.ResponseWriter, r *http.Request) {
	const op = "internal.server.handlers.handlers.GetExchangeRate"

	pair := r.PathValue("pair")
	if len(pair) < 6 {
		log.Printf("%s: %v", op, ErrInvalidInputData)
		errorJSON(w, "invalid currency pair format", http.StatusBadRequest)
		return
	}

	// it supposed that each code is three symbols length
	baseCode := pair[:3]
	targetCode := pair[3:]

	rate, err := h.exchangeRateSrv.GetExchangeRate(r.Context(), baseCode, targetCode)
	if err != nil {
		log.Printf("%s: %v", op, err)
		if errors.Is(err, repository.ErrExchangeRateNotFound) {
			errorJSON(w, "exchange rate not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, repository.ErrCurrencyNotFound) {
			errorJSON(w, "currency not found", http.StatusNotFound)
			return
		}
		errorJSON(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rate)
}

func (h *Handlers) CreateExchangeRate(w http.ResponseWriter, r *http.Request) {
	const op = "internal.server.handlers.handlers.CreateExchangeRate"

	if err := r.ParseForm(); err != nil {
		log.Printf("%s: %v", op, err)
		errorJSON(w, "invalid form data", http.StatusBadRequest)
		return
	}

	baseCode := r.PostFormValue("baseCurrencyCode")
	targetCode := r.PostFormValue("targetCurrencyCode")
	rateStr := r.PostFormValue("rate")

	if baseCode == "" || targetCode == "" || rateStr == "" {
		log.Printf("%s: %v", op, ErrInvalidInputData)
		errorJSON(w, "all fields are required", http.StatusBadRequest)
		return
	}

	rate, err := strconv.ParseFloat(rateStr, 64)
	if err != nil {
		log.Printf("%s: %v", op, err)
		errorJSON(w, "invalid rate format", http.StatusBadRequest)
		return
	}

	createdRate, err := h.exchangeRateSrv.AddExchangeRate(r.Context(), baseCode, targetCode, rate)
	if err != nil {
		log.Printf("%s: %v", op, err)
		if errors.Is(err, repository.ErrCurrencyNotFound) {
			errorJSON(w, "one or both currencies not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, repository.ErrExchangeRateExists) {
			errorJSON(w, "exchange rate already exists", http.StatusConflict)
			return
		}
		errorJSON(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdRate)
}

func (h *Handlers) UpdateExchangeRate(w http.ResponseWriter, r *http.Request) {
	const op = "internal.server.handlers.handlers.UpdateExchangeRate"

	pair := r.PathValue("pair")
	if len(pair) < 6 {
		log.Printf("%s: %v", op, ErrInvalidInputData)
		errorJSON(w, "invalid currency pair format", http.StatusBadRequest)
		return
	}

	// it supposed that each code is three symbols length
	baseCode := pair[:3]
	targetCode := pair[3:]

	if err := r.ParseForm(); err != nil {
		log.Printf("%s: %v", op, err)
		errorJSON(w, "invalid form data", http.StatusBadRequest)
		return
	}

	rateStr := r.PostFormValue("rate")
	if rateStr == "" {
		log.Printf("%s: %v", op, ErrInvalidInputData)
		errorJSON(w, "rate is required", http.StatusBadRequest)
		return
	}

	rate, err := strconv.ParseFloat(rateStr, 64)
	if err != nil {
		log.Printf("%s: %v", op, err)
		errorJSON(w, "invalid rate format", http.StatusBadRequest)
		return
	}

	updatedRate, err := h.exchangeRateSrv.UpdateExchangeRate(r.Context(), baseCode, targetCode, rate)
	if err != nil {
		log.Printf("%s: %v", op, err)
		if errors.Is(err, repository.ErrExchangeRateNotFound) {
			errorJSON(w, "exchange rate not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, repository.ErrCurrencyNotFound) {
			errorJSON(w, "one or both currencies not found", http.StatusNotFound)
			return
		}
		errorJSON(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedRate)
}
