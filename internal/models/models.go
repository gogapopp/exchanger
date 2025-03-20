package models

type Currency struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
	Sign string `json:"sign"`
}

type ExchangeRate struct {
	ID             int      `json:"id"`
	BaseCurrency   Currency `json:"baseCurrency"`
	TargetCurrency Currency `json:"targetCurrency"`
	Rate           float64  `json:"rate"`
}

type CurrencyConversion struct {
	BaseCurrency    Currency `json:"baseCurrency"`
	TargetCurrency  Currency `json:"targetCurrency"`
	Rate            float64  `json:"rate"`
	Amount          float64  `json:"amount"`
	ConvertedAmount float64  `json:"convertedAmount"`
}
