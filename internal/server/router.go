package server

import (
	"exchanger/internal/server/handlers"
	"net/http"
)

func Routes(h *handlers.Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /currencies", h.GetCurrencies)
	mux.HandleFunc("GET /currency/{code}", h.GetCurrency)
	mux.HandleFunc("POST /currencies", h.CreateCurrency)

	mux.HandleFunc("GET /exchangeRates", h.GetExchangeRates)
	mux.HandleFunc("GET /exchangeRate/{pair}", h.GetExchangeRate)
	mux.HandleFunc("POST /exchangeRates", h.CreateExchangeRate)
	mux.HandleFunc("PATCH /exchangeRate/{pair}", h.UpdateExchangeRate)

	mux.HandleFunc("GET /exchange", h.ExchangeCurrency)

	return mux
}
