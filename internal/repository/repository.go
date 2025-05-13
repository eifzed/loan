package repository

import (
	"context"
	"loan/internal/domain"
)

// LoanRepository defines the interface for loan data operations
type LoanRepository interface {
	SaveLoan(ctx context.Context, loan *domain.Loan) error
	GetLoanByID(ctx context.Context, id string) (*domain.Loan, error)
	ListLoans(ctx context.Context, page, pageSize int) ([]*domain.Loan, int, error)

	SaveApproval(ctx context.Context, approval *domain.Approval) error

	SaveInvestment(ctx context.Context, investment *domain.Investment) error
	GetLoanInvestments(ctx context.Context, loanID string) ([]*domain.Investment, error)

	SaveDisbursement(ctx context.Context, disbursement *domain.Disbursement) error
}
