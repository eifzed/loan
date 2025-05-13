package service

import (
	"context"
	"fmt"
)

// EmailService is an interface for sending email notifications
type EmailService interface {
	SendInvestmentNotification(ctx context.Context, investorID, loanID string, agreementLetterURL string) error
}

type MockEmailService struct{}

func NewMockEmailService() *MockEmailService {
	return &MockEmailService{}
}

func (s *MockEmailService) SendInvestmentNotification(ctx context.Context, investorID, loanID string, agreementLetterURL string) error {
	// Printing sending email as simulation
	fmt.Printf("Sending email to investor %s for loan %s with agreement letter: %s\n",
		investorID, loanID, agreementLetterURL)
	return nil
}
