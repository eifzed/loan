package service

import (
	"context"
	"errors"
	"fmt"
	"loan/internal/domain"
	"loan/internal/repository"
	"log"
	"time"
)

// LoanService handles the business logic for loan operations
type LoanService struct {
	repo         repository.LoanRepository
	emailService EmailService
}

func NewLoanService(repo repository.LoanRepository, emailService EmailService) *LoanService {
	return &LoanService{
		repo:         repo,
		emailService: emailService,
	}
}

func (s *LoanService) CreateLoan(ctx context.Context, borrowerID string, principalAmount, rate, roi float64) (*domain.Loan, error) {
	if borrowerID == "" {
		return nil, errors.New("borrower ID cannot be empty")
	}

	if principalAmount <= 0 {
		return nil, errors.New("principal amount must be greater than zero")
	}

	if rate < 0 {
		return nil, errors.New("rate cannot be negative")
	}

	if roi < 0 {
		return nil, errors.New("ROI cannot be negative")
	}

	loan := domain.NewLoan(borrowerID, principalAmount, rate, roi)

	if err := s.repo.SaveLoan(ctx, loan); err != nil {
		return nil, err
	}

	return loan, nil
}

// GetLoan retrieves a loan by its ID
func (s *LoanService) GetLoan(ctx context.Context, id string) (*domain.Loan, error) {
	return s.repo.GetLoanByID(ctx, id)
}

// ListLoans retrieves loans based on filters and pagination
func (s *LoanService) ListLoans(ctx context.Context, page, pageSize int) ([]*domain.Loan, int, error) {
	return s.repo.ListLoans(ctx, page, pageSize)
}

// ApproveLoan changes a loan state from PROPOSED to APPROVED
func (s *LoanService) ApproveLoan(ctx context.Context, loanID, proofPictureURL, fieldValidatorID string, approvalDate time.Time) (*domain.Loan, error) {
	loan, err := s.repo.GetLoanByID(ctx, loanID)
	if err != nil {
		return nil, err
	}

	approval, err := domain.NewApproval(loanID, proofPictureURL, fieldValidatorID, approvalDate)
	if err != nil {
		return nil, err
	}

	if err := loan.Approve(approval); err != nil {
		return nil, err
	}

	// use transaction in real implementation with commit and rollback defer function

	if err := s.repo.SaveApproval(ctx, approval); err != nil {
		return nil, err
	}

	if err := s.repo.SaveLoan(ctx, loan); err != nil {
		return nil, err
	}

	return loan, nil
}

func (s *LoanService) AddInvestment(ctx context.Context, loanID, investorID string, amount float64) (*domain.Investment, error) {
	loan, err := s.repo.GetLoanByID(ctx, loanID)
	if err != nil {
		return nil, err
	}

	// use transaction in real implementation with commit and rollback defer function

	investment, err := domain.NewInvestment(loanID, investorID, amount)
	if err != nil {
		return nil, err
	}
	fmt.Println("investment id", investment.ID)

	// inv_9ea12a9b-85f1-416b-8737-a9ba0c4e1165

	if err := loan.AddInvestment(investment); err != nil {
		return nil, err
	}

	if err := s.repo.SaveInvestment(ctx, investment); err != nil {
		return nil, err
	}

	if err := s.repo.SaveLoan(ctx, loan); err != nil {
		return nil, err
	}

	// If the loan has transitioned to INVESTED state, send notifications to all investors
	if loan.State == domain.LoanStateInvested {
		for _, inv := range loan.Investments {
			err := s.emailService.SendInvestmentNotification(ctx, inv.InvestorID, loanID, loan.AgreementLetterURL)
			if err != nil {
				// Log error but continue processing
				// In a real implementation, we might use a retry mechanism
				// or queue for handling notification failures (using NSQ, Kafka, etc.)
				log.Printf("Failed to send notification to investor %s: %v", inv.InvestorID, err)
			}
		}
	}

	return investment, nil
}

func (s *LoanService) GetLoanInvestments(ctx context.Context, loanID string) ([]*domain.Investment, error) {
	return s.repo.GetLoanInvestments(ctx, loanID)
}

func (s *LoanService) DisburseLoan(ctx context.Context, loanID, agreementDocumentURL, fieldOfficerID string, disbursementDate time.Time) (*domain.Loan, error) {
	loan, err := s.repo.GetLoanByID(ctx, loanID)
	if err != nil {
		return nil, err
	}

	disbursement, err := domain.NewDisbursement(loanID, agreementDocumentURL, fieldOfficerID, disbursementDate)
	if err != nil {
		return nil, err
	}

	if err := loan.Disburse(disbursement); err != nil {
		return nil, err
	}

	if err := s.repo.SaveDisbursement(ctx, disbursement); err != nil {
		return nil, err
	}

	if err := s.repo.SaveLoan(ctx, loan); err != nil {
		return nil, err
	}

	return loan, nil
}
