package repository

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type repository struct {
	conn *sql.DB
}

func New(ctx context.Context) (*repository, error) {
	const op = "internal.repository.repository.New"

	db, err := sql.Open("sqlite3", "storage.db")
	if err != nil {
		return &repository{}, fmt.Errorf("%s: %v", op, err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return &repository{}, fmt.Errorf("%s: %v", op, err)
	}

	_, err = db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS Currencies (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		code TEXT NOT NULL UNIQUE,
		full_name TEXT NOT NULL,
		sign TEXT NOT NULL
	);`)
	if err != nil {
		return &repository{}, fmt.Errorf("%s: %v", op, err)
	}

	_, err = db.ExecContext(ctx, `
	CREATE INDEX IF NOT EXISTS idx_currencies_code 
	ON Currencies(code);`)
	if err != nil {
		return &repository{}, fmt.Errorf("%s: %v", op, err)
	}

	_, err = db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS ExchangeRates (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		base_currency_id INTEGER NOT NULL,
		target_currency_id INTEGER NOT NULL,
		rate REAL NOT NULL,

		FOREIGN KEY (base_currency_id) REFERENCES Currencies(ID),
		FOREIGN KEY (target_currency_id) REFERENCES Currencies(ID),
		
		UNIQUE (base_currency_id, target_currency_id)
	);`)
	if err != nil {
		return &repository{}, fmt.Errorf("%s: %v", op, err)
	}

	_, err = db.ExecContext(ctx, "INSERT OR IGNORE INTO Currencies (code, full_name, sign) VALUES ('USD', 'United States dollar', '$')")
	if err != nil {
		return &repository{}, fmt.Errorf("%s: %v", op, err)
	}
	_, err = db.ExecContext(ctx, "INSERT OR IGNORE INTO Currencies (code, full_name, sign) VALUES ('RUB', 'Russian Ruble', '₽')")
	if err != nil {
		return &repository{}, fmt.Errorf("%s: %v", op, err)
	}
	_, err = db.ExecContext(ctx, "INSERT OR IGNORE INTO Currencies (code, full_name, sign) VALUES ('EUR', 'Euro', '€')")
	if err != nil {
		return &repository{}, fmt.Errorf("%s: %v", op, err)
	}

	return &repository{conn: db}, nil
}

func (r *repository) Close() error {
	return r.conn.Close()
}
