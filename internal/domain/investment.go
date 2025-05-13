package domain

import (
	"errors"
	"loan/util"
	"time"
)

// Investment represents an investment made in a loan
type Investment struct {
	ID         string    `json:"id"`
	LoanID     string    `json:"loan_id"`
	InvestorID string    `json:"investor_id"`
	Amount     float64   `json:"amount"`
	InvestedAt time.Time `json:"invested_at"`
}

func NewInvestment(loanID, investorID string, amount float64) (*Investment, error) {
	if loanID == "" {
		return nil, errors.New("loan ID cannot be empty")
	}

	if investorID == "" {
		return nil, errors.New("investor ID cannot be empty")
	}

	if amount <= 0 {
		return nil, errors.New("investment amount must be greater than zero")
	}

	return &Investment{
		ID:         "inv_" + util.GenerateUUID(),
		LoanID:     loanID,
		InvestorID: investorID,
		Amount:     amount,
		InvestedAt: time.Now(),
	}, nil
}
