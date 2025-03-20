package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"exchanger/internal/models"
	"exchanger/internal/repository"
	"log"
	"net/http"
)

type currencyService interface {
	GetAllCurrencies(ctx context.Context) ([]models.Currency, error)
	GetCurrencyByCode(ctx context.Context, code string) (models.Currency, error)
	AddCurrency(ctx context.Context, currency models.Currency) (models.Currency, error)
}

func (h *Handlers) GetCurrencies(w http.ResponseWriter, r *http.Request) {
	const op = "internal.server.handlers.handlers.GetCurrencies"

	currencies, err := h.currencySrv.GetAllCurrencies(r.Context())
	if err != nil {
		log.Printf("%s: %v", op, err)
		if errors.Is(err, repository.ErrCurrencyNotFound) {
			errorJSON(w, "currency not found", http.StatusNotFound)
			return
		}
		errorJSON(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currencies)
}

func (h *Handlers) GetCurrency(w http.ResponseWriter, r *http.Request) {
	const op = "internal.server.handlers.handlers.GetCurrency"

	code := r.PathValue("code")
	if code == "" {
		log.Printf("%s: %v", op, ErrInvalidInputData)
		errorJSON(w, "currency code is required", http.StatusBadRequest)
		return
	}

	currency, err := h.currencySrv.GetCurrencyByCode(r.Context(), code)
	if err != nil {
		log.Printf("%s: %v", op, err)
		if errors.Is(err, repository.ErrCurrencyNotFound) {
			errorJSON(w, "currency not found", http.StatusNotFound)
			return
		}
		errorJSON(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currency)
}

func (h *Handlers) CreateCurrency(w http.ResponseWriter, r *http.Request) {
	const op = "internal.server.handlers.handlers.CreateCurrency"

	if err := r.ParseForm(); err != nil {
		errorJSON(w, "invalid form data", http.StatusBadRequest)
		return
	}

	currency := models.Currency{
		Name: r.PostFormValue("name"),
		Code: r.PostFormValue("code"),
		Sign: r.PostFormValue("sign"),
	}

	if currency.Name == "" || currency.Code == "" || currency.Sign == "" {
		log.Printf("%s: %v", op, errors.New("invalid input data"))
		errorJSON(w, "all fields are required", http.StatusBadRequest)
		return
	}

	createdCurrency, err := h.currencySrv.AddCurrency(r.Context(), currency)
	if err != nil {
		log.Printf("%s: %v", op, err)
		if errors.Is(err, repository.ErrCurrencyExists) {
			errorJSON(w, "currency already exists", http.StatusConflict)
			return
		}
		errorJSON(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdCurrency)
}
