package repository

import (
	"context"
	"errors"
	"loan/internal/domain"
	"sync"
)

// MockLoanRepository is an in-memory implementation of LoanRepository
type MockLoanRepository struct {
	loans         map[string]*domain.Loan
	approvals     map[string]*domain.Approval
	investments   map[string]*domain.Investment
	disbursements map[string]*domain.Disbursement
	mutex         sync.RWMutex
}

func NewMockLoanRepository() *MockLoanRepository {
	return &MockLoanRepository{
		loans:         make(map[string]*domain.Loan),
		approvals:     make(map[string]*domain.Approval),
		investments:   make(map[string]*domain.Investment),
		disbursements: make(map[string]*domain.Disbursement),
	}
}

func (r *MockLoanRepository) SaveLoan(ctx context.Context, loan *domain.Loan) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.loans[loan.ID] = loan
	return nil
}

func (r *MockLoanRepository) GetLoanByID(ctx context.Context, id string) (*domain.Loan, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	loan, exists := r.loans[id]
	if !exists {
		return nil, errors.New("loan not found")
	}

	return loan, nil
}

func (r *MockLoanRepository) ListLoans(ctx context.Context, page, pageSize int) ([]*domain.Loan, int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []*domain.Loan
	for _, loan := range r.loans {
		result = append(result, loan)
	}

	total := len(result)

	if page > 0 && pageSize > 0 {
		start := (page - 1) * pageSize
		end := start + pageSize

		if start >= total {
			return []*domain.Loan{}, total, nil
		}

		if end > total {
			end = total
		}

		if start < total {
			result = result[start:end]
		} else {
			result = []*domain.Loan{}
		}
	}

	return result, total, nil
}

func (r *MockLoanRepository) SaveApproval(ctx context.Context, approval *domain.Approval) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	loan, exists := r.loans[approval.LoanID]
	if !exists {
		return errors.New("loan not found")
	}

	r.approvals[approval.LoanID] = approval
	loan.Approval = approval

	return nil
}

func (r *MockLoanRepository) SaveInvestment(ctx context.Context, investment *domain.Investment) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	loan, exists := r.loans[investment.LoanID]
	if !exists {
		return errors.New("loan not found")
	}

	r.investments[investment.ID] = investment

	found := false
	for i, inv := range loan.Investments {
		if inv.ID == investment.ID {
			loan.Investments[i] = investment
			found = true
			break
		}
	}

	if !found {
		loan.Investments = append(loan.Investments, investment)
	}

	return nil
}

func (r *MockLoanRepository) GetLoanInvestments(ctx context.Context, loanID string) ([]*domain.Investment, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	loan, exists := r.loans[loanID]
	if !exists {
		return nil, errors.New("loan not found")
	}

	return loan.Investments, nil
}

func (r *MockLoanRepository) SaveDisbursement(ctx context.Context, disbursement *domain.Disbursement) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	loan, exists := r.loans[disbursement.LoanID]
	if !exists {
		return errors.New("loan not found")
	}

	r.disbursements[disbursement.LoanID] = disbursement
	loan.Disbursement = disbursement

	return nil
}
