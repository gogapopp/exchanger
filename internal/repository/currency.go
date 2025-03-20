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
	ErrCurrencyNotFound = errors.New("currency not found")
	ErrCurrencyExists   = errors.New("currency already exists")
)

func (r *repository) GetAllCurrencies(ctx context.Context) ([]models.Currency, error) {
	const op = "internal.repository.repository.GetAllCurrencies"

	rows, err := r.conn.QueryContext(ctx, "SELECT ID, code, full_name, sign FROM Currencies")
	if err != nil {
		return []models.Currency{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	hasRows := false

	var currencies []models.Currency
	for rows.Next() {
		hasRows = true
		var c models.Currency
		if err := rows.Scan(&c.ID, &c.Code, &c.Name, &c.Sign); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		currencies = append(currencies, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !hasRows {
		return []models.Currency{}, fmt.Errorf("%s: %w", op, ErrCurrencyNotFound)
	}

	return currencies, nil
}

func (r *repository) GetCurrencyByCode(ctx context.Context, code string) (models.Currency, error) {
	const op = "internal.repository.repository.GetCurrencyByCode"

	var c models.Currency
	err := r.conn.QueryRowContext(ctx, "SELECT ID, full_name, code, sign FROM Currencies WHERE code = ?", code).
		Scan(&c.ID, &c.Name, &c.Code, &c.Sign)
	if err == sql.ErrNoRows {
		return models.Currency{}, fmt.Errorf("%s: %w", op, ErrCurrencyNotFound)
	} else if err != nil {
		return models.Currency{}, fmt.Errorf("%s: %w", op, err)
	}

	return c, nil
}

func (r *repository) AddCurrency(ctx context.Context, currency models.Currency) (models.Currency, error) {
	const op = "internal.repository.repository.AddCurrency"

	var id int
	err := r.conn.QueryRowContext(ctx, "INSERT INTO Currencies (full_name, code, sign) VALUES (?, ?, ?) RETURNING ID",
		currency.Name, currency.Code, currency.Sign).Scan(&id)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				return models.Currency{}, fmt.Errorf("%s: %w", op, ErrCurrencyExists)
			}
		}
		return models.Currency{}, fmt.Errorf("%s: %w", op, err)
	}

	currency.ID = id

	return currency, nil
}
