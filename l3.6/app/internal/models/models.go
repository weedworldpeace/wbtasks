package models

import (
	"errors"
	"time"
)

type TransactionType string

const (
	TransactionTypeIncome  TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"
)

type Transaction struct {
	ID          string          `json:"id" db:"id" `
	UserID      string          `json:"user_id" db:"user_id"`
	Amount      float64         `json:"amount" db:"amount"`
	Type        TransactionType `json:"type" db:"type"`
	Category    string          `json:"category" db:"category"`
	Description string          `json:"description" db:"description"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

type AnalytiscHelper struct {
	Amount  float64 `json:"amount"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

type Analytics struct {
	Total        int             `json:"total"`
	Sum          AnalytiscHelper `json:"sum"`
	Average      AnalytiscHelper `json:"average"`
	Median       AnalytiscHelper `json:"median"`
	Percentile90 AnalytiscHelper `json:"percentile_90"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	ErrInvalidRequestBody            = errors.New("invalid request body")
	ErrInvalidTransactionID          = errors.New("invalid transaction ID")
	ErrInvalidTransactionAmount      = errors.New("invalid transaction amount")
	ErrInvalidTransactionType        = errors.New("invalid transaction type")
	ErrInvalidTransactionCategory    = errors.New("invalid transaction category")
	ErrInvalidTransactionDescription = errors.New("invalid transaction description")
	ErrInvalidTransactionGenre       = errors.New("invalid transaction genre")
	ErrOnDatabase                    = errors.New("error on database")
	ErrTransactionNotFound           = errors.New("transaction not found")
	ErrInvalidTimestamp              = errors.New("invalid timestamp")
	ErrInvalidLimit                  = errors.New("invalid limit")
	ErrInvalidOffset                 = errors.New("invalid offset")
	ErrListIsEmpty                   = errors.New("list is empty")
	ErrInvalidPort                   = errors.New("invalid port")
	ErrInvalidReleaseMode            = errors.New("invalid release mode")
)

// type DailySummary struct {
// 	SummaryDate       time.Time `json:"summary_date" db:"summary_date"`
// 	TotalTransactions int       `json:"total_transactions" db:"total_transactions"`
// 	TotalAmount       float64   `json:"total_amount" db:"total_amount"`
// 	TotalIncome       float64   `json:"total_income" db:"total_income"`
// 	TotalExpense      float64   `json:"total_expense" db:"total_expense"`
// 	UniqueUsers       int       `json:"unique_users" db:"unique_users"`
// 	TopCategory       *string   `json:"top_category,omitempty" db:"top_category"`
// 	CreatedAt         time.Time `json:"created_at" db:"created_at"`
// }
