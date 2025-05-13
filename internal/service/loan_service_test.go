package service_test

import (
	"context"
	"loan/internal/domain"
	"loan/internal/repository"
	"loan/internal/service"
	"testing"
	"time"
)

// MockEmailService that tracks notifications
type MockEmailServiceWithTracking struct {
	NotificationsSent map[string]bool
}

func NewMockEmailServiceWithTracking() *MockEmailServiceWithTracking {
	return &MockEmailServiceWithTracking{
		NotificationsSent: make(map[string]bool),
	}
}

func (s *MockEmailServiceWithTracking) SendInvestmentNotification(ctx context.Context, investorID, loanID string, agreementLetterURL string) error {
	key := investorID + ":" + loanID
	s.NotificationsSent[key] = true
	return nil
}

func TestCreateLoan(t *testing.T) {
	// Arrange
	repo := repository.NewMockLoanRepository()
	emailService := service.NewMockEmailService()
	loanService := service.NewLoanService(repo, emailService)

	// Act
	loan, err := loanService.CreateLoan(context.Background(), "borrower123", 1000.0, 0.1, 0.08)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if loan.BorrowerID != "borrower123" {
		t.Errorf("Expected BorrowerID to be borrower123, got %s", loan.BorrowerID)
	}
	if loan.State != domain.LoanStateProposed {
		t.Errorf("Expected State to be PROPOSED, got %s", loan.State)
	}

	// Verify loan was saved
	savedLoan, err := loanService.GetLoan(context.Background(), loan.ID)
	if err != nil {
		t.Fatalf("Expected to retrieve loan, got error: %v", err)
	}

	if savedLoan.ID != loan.ID {
		t.Errorf("Expected to retrieve same loan ID, got %s", savedLoan.ID)
	}
}

func TestInvalidLoanCreation(t *testing.T) {
	// Arrange
	repo := repository.NewMockLoanRepository()
	emailService := service.NewMockEmailService()
	loanService := service.NewLoanService(repo, emailService)

	// Act & Assert
	_, err := loanService.CreateLoan(context.Background(), "", 1000.0, 0.1, 0.08)
	if err == nil {
		t.Error("Expected error for empty borrower ID, got nil")
	}

	_, err = loanService.CreateLoan(context.Background(), "borrower123", -100, 0.1, 0.08)
	if err == nil {
		t.Error("Expected error for negative principal amount, got nil")
	}

	_, err = loanService.CreateLoan(context.Background(), "borrower123", 1000.0, -0.1, 0.08)
	if err == nil {
		t.Error("Expected error for negative rate, got nil")
	}
}

func TestApproveLoan(t *testing.T) {
	// Arrange
	repo := repository.NewMockLoanRepository()
	emailService := service.NewMockEmailService()
	loanService := service.NewLoanService(repo, emailService)

	// Create a loan
	loan, _ := loanService.CreateLoan(context.Background(), "borrower123", 1000.0, 0.1, 0.08)

	// Act
	approvalDate := time.Now()
	updatedLoan, err := loanService.ApproveLoan(
		context.Background(),
		loan.ID,
		"proof.jpg",
		"validator123",
		approvalDate,
	)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error when approving loan, got %v", err)
	}

	if updatedLoan.State != domain.LoanStateApproved {
		t.Errorf("Expected loan state to be APPROVED, got %s", updatedLoan.State)
	}

	if updatedLoan.Approval == nil {
		t.Fatal("Expected loan approval to be set, got nil")
	}

	if updatedLoan.Approval.ProofPictureURL != "proof.jpg" {
		t.Errorf("Expected proof picture URL to be proof.jpg, got %s", updatedLoan.Approval.ProofPictureURL)
	}

	if updatedLoan.Approval.FieldValidatorID != "validator123" {
		t.Errorf("Expected field validator ID to be validator123, got %s", updatedLoan.Approval.FieldValidatorID)
	}
}

func TestAddInvestment(t *testing.T) {
	// Arrange
	repo := repository.NewMockLoanRepository()
	emailService := NewMockEmailServiceWithTracking()
	loanService := service.NewLoanService(repo, emailService)

	// Create and approve a loan
	loan, _ := loanService.CreateLoan(context.Background(), "borrower123", 1000.0, 0.1, 0.08)
	_, _ = loanService.ApproveLoan(
		context.Background(),
		loan.ID,
		"proof.jpg",
		"validator123",
		time.Now(),
	)

	// Act
	investment, err := loanService.AddInvestment(
		context.Background(),
		loan.ID,
		"investor123",
		600.0,
	)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error when adding investment, got %v", err)
	}

	if investment.LoanID != loan.ID {
		t.Errorf("Expected investment loan ID to be %s, got %s", loan.ID, investment.LoanID)
	}

	if investment.InvestorID != "investor123" {
		t.Errorf("Expected investor ID to be investor123, got %s", investment.InvestorID)
	}

	if investment.Amount != 600.0 {
		t.Errorf("Expected investment amount to be 600.0, got %f", investment.Amount)
	}

	// Verify loan state (should still be APPROVED after partial investment)
	updatedLoan, _ := loanService.GetLoan(context.Background(), loan.ID)
	if updatedLoan.State != domain.LoanStateApproved {
		t.Errorf("Expected loan state to still be APPROVED after partial investment, got %s", updatedLoan.State)
	}

	// Add second investment to fully fund the loan
	_, err = loanService.AddInvestment(
		context.Background(),
		loan.ID,
		"investor456",
		400.0,
	)
	if err != nil {
		t.Fatalf("Expected no error when adding second investment, got %v", err)
	}

	// Verify loan state (should now be INVESTED)
	updatedLoan, _ = loanService.GetLoan(context.Background(), loan.ID)
	if updatedLoan.State != domain.LoanStateInvested {
		t.Errorf("Expected loan state to be INVESTED after full investment, got %s", updatedLoan.State)
	}

	// Verify notifications were sent to both investors
	if !emailService.NotificationsSent["investor123:"+loan.ID] {
		t.Error("Expected notification to be sent to first investor")
	}
	if !emailService.NotificationsSent["investor456:"+loan.ID] {
		t.Error("Expected notification to be sent to second investor")
	}
}

func TestDisburseLoan(t *testing.T) {
	// Arrange
	repo := repository.NewMockLoanRepository()
	emailService := service.NewMockEmailService()
	loanService := service.NewLoanService(repo, emailService)

	// Create, approve, and invest in a loan
	loan, _ := loanService.CreateLoan(context.Background(), "borrower123", 1000.0, 0.1, 0.08)
	_, _ = loanService.ApproveLoan(
		context.Background(),
		loan.ID,
		"proof.jpg",
		"validator123",
		time.Now(),
	)
	_, _ = loanService.AddInvestment(
		context.Background(),
		loan.ID,
		"investor123",
		1000.0,
	)

	// Act
	disbursementDate := time.Now()
	updatedLoan, err := loanService.DisburseLoan(
		context.Background(),
		loan.ID,
		"agreement.pdf",
		"officer123",
		disbursementDate,
	)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error when disbursing loan, got %v", err)
	}

	if updatedLoan.State != domain.LoanStateDisbursed {
		t.Errorf("Expected loan state to be DISBURSED, got %s", updatedLoan.State)
	}

	if updatedLoan.Disbursement == nil {
		t.Fatal("Expected loan disbursement to be set, got nil")
	}

	if updatedLoan.Disbursement.AgreementDocumentURL != "agreement.pdf" {
		t.Errorf("Expected agreement document URL to be agreement.pdf, got %s", updatedLoan.Disbursement.AgreementDocumentURL)
	}

	if updatedLoan.Disbursement.FieldOfficerID != "officer123" {
		t.Errorf("Expected field officer ID to be officer123, got %s", updatedLoan.Disbursement.FieldOfficerID)
	}
}

func TestGetLoanInvestments(t *testing.T) {
	// Arrange
	repo := repository.NewMockLoanRepository()
	emailService := service.NewMockEmailService()
	loanService := service.NewLoanService(repo, emailService)

	// Create and approve a loan
	loan, _ := loanService.CreateLoan(context.Background(), "borrower123", 1000.0, 0.1, 0.08)
	_, _ = loanService.ApproveLoan(
		context.Background(),
		loan.ID,
		"proof.jpg",
		"validator123",
		time.Now(),
	)

	// Add investments
	_, _ = loanService.AddInvestment(context.Background(), loan.ID, "investor1", 400.0)
	_, _ = loanService.AddInvestment(context.Background(), loan.ID, "investor2", 600.0)

	// Act
	investments, err := loanService.GetLoanInvestments(context.Background(), loan.ID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error when getting loan investments, got %v", err)
	}

	if len(investments) != 2 {
		t.Fatalf("Expected 2 investments, got %d", len(investments))
	}

	// Verify total investment amount
	var totalAmount float64
	for _, inv := range investments {
		totalAmount += inv.Amount
	}

	if totalAmount != 1000.0 {
		t.Errorf("Expected total investment amount to be 1000.0, got %f", totalAmount)
	}
}
