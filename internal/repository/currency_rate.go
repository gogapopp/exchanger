package repository

import (
	"context"
	"database/sql"
	"errors"
	"exchanger/internal/models"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrExchangeRateNotFound = errors.New("exchange rate not found")
	ErrExchangeRateExists   = errors.New("exchange rate already exists")
)

func (r *repository) GetAllExchangeRates(ctx context.Context) ([]models.ExchangeRate, error) {
	const op = "internal.repository.repository.GetAllExchangeRates"

	query := `
	SELECT er.ID, er.rate, 
		bc.ID, bc.code, bc.full_name, bc.sign,
		tc.ID, tc.code, tc.full_name, tc.sign 
	FROM ExchangeRates er
	JOIN Currencies bc ON er.base_currency_id = bc.ID
	JOIN Currencies tc ON er.target_currency_id = tc.ID
	`

	rows, err := r.conn.QueryContext(ctx, query)
	if err != nil {
		return []models.ExchangeRate{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	hasRows := false

	var rates []models.ExchangeRate
	for rows.Next() {
		hasRows = true

		var er models.ExchangeRate
		var bc models.Currency
		var tc models.Currency

		if err := rows.Scan(
			&er.ID, &er.Rate,
			&bc.ID, &bc.Code, &bc.Name, &bc.Sign,
			&tc.ID, &tc.Code, &tc.Name, &tc.Sign); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		er.BaseCurrency = bc
		er.TargetCurrency = tc

		rates = append(rates, er)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !hasRows {
		return []models.ExchangeRate{}, fmt.Errorf("%s: %w", op, ErrExchangeRateNotFound)
	}

	return rates, nil
}

func (r *repository) GetExchangeRate(ctx context.Context, baseCode, targetCode string) (models.ExchangeRate, error) {
	const op = "internal.repository.repository.GetExchangeRate"

	query := `
	SELECT er.ID, er.rate, 
		bc.ID, bc.code, bc.full_name, bc.sign,
		tc.ID, tc.code, tc.full_name, tc.sign 
	FROM ExchangeRates er
	JOIN Currencies bc ON er.base_currency_id = bc.ID
	JOIN Currencies tc ON er.target_currency_id = tc.ID
	WHERE bc.code = ? AND tc.code = ?
	`

	var er models.ExchangeRate
	var bc models.Currency
	var tc models.Currency

	err := r.conn.QueryRowContext(ctx, query, baseCode, targetCode).Scan(
		&er.ID, &er.Rate,
		&bc.ID, &bc.Code, &bc.Name, &bc.Sign,
		&tc.ID, &tc.Code, &tc.Name, &tc.Sign)
	if err == sql.ErrNoRows {
		return models.ExchangeRate{}, fmt.Errorf("%s: %w", op, ErrExchangeRateNotFound)
	} else if err != nil {
		return models.ExchangeRate{}, fmt.Errorf("%s: %w", op, err)
	}

	er.BaseCurrency = bc
	er.TargetCurrency = tc

	return er, nil
}

func (r *repository) AddExchangeRate(ctx context.Context, baseCode, targetCode string, rate float64) (models.ExchangeRate, error) {
	const op = "internal.repository.repository.AddExchangeRate"

	baseCurrency, err := r.GetCurrencyByCode(ctx, baseCode)
	if err != nil {
		return models.ExchangeRate{}, fmt.Errorf("%s: %w", op, err)
	}

	targetCurrency, err := r.GetCurrencyByCode(ctx, targetCode)
	if err != nil {
		return models.ExchangeRate{}, fmt.Errorf("%s: %w", op, err)
	}

	var id int
	err = r.conn.QueryRowContext(
		ctx,
		"INSERT INTO ExchangeRates (base_currency_id, target_currency_id, rate) VALUES (?, ?, ?) RETURNING ID",
		baseCurrency.ID, targetCurrency.ID, rate,
	).Scan(&id)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				return models.ExchangeRate{}, fmt.Errorf("%s: %w", op, ErrExchangeRateExists)
			}
		}

		return models.ExchangeRate{}, fmt.Errorf("%s: %w", op, err)
	}

	return models.ExchangeRate{
		ID:             id,
		BaseCurrency:   baseCurrency,
		TargetCurrency: targetCurrency,
		Rate:           rate,
	}, nil
}

func (r *repository) UpdateExchangeRate(ctx context.Context, baseCode, targetCode string, rate float64) (models.ExchangeRate, error) {
	const op = "internal.repository.repository.UpdateExchangeRate"

	baseCurrency, err := r.GetCurrencyByCode(ctx, baseCode)
	if err != nil {
		return models.ExchangeRate{}, fmt.Errorf("%s: %w", op, err)
	}

	targetCurrency, err := r.GetCurrencyByCode(ctx, targetCode)
	if err != nil {
		return models.ExchangeRate{}, fmt.Errorf("%s: %w", op, err)
	}

	result, err := r.conn.ExecContext(
		ctx,
		"UPDATE ExchangeRates SET rate = ? WHERE base_currency_id = ? AND target_currency_id = ?",
		rate, baseCurrency.ID, targetCurrency.ID,
	)
	if err != nil {
		return models.ExchangeRate{}, fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.ExchangeRate{}, fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return models.ExchangeRate{}, fmt.Errorf("%s: %w", op, ErrExchangeRateNotFound)
	}

	var id int
	err = r.conn.QueryRowContext(
		ctx,
		"SELECT ID FROM ExchangeRates WHERE base_currency_id = ? AND target_currency_id = ?",
		baseCurrency.ID, targetCurrency.ID,
	).Scan(&id)
	if err != nil {
		return models.ExchangeRate{}, fmt.Errorf("%s: %w", op, err)
	}

	return models.ExchangeRate{
		ID:             id,
		BaseCurrency:   baseCurrency,
		TargetCurrency: targetCurrency,
		Rate:           rate,
	}, nil
}
