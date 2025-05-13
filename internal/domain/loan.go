package domain

import (
	"errors"
	"loan/util"
	"time"
)

type LoanState string

const (
	LoanStateProposed  LoanState = "PROPOSED"
	LoanStateApproved  LoanState = "APPROVED"
	LoanStateInvested  LoanState = "INVESTED"
	LoanStateDisbursed LoanState = "DISBURSED"
)

type Loan struct {
	ID                 string    `json:"id"`
	BorrowerID         string    `json:"borrower_id"`
	PrincipalAmount    float64   `json:"principal_amount"`
	Rate               float64   `json:"rate"`
	ROI                float64   `json:"roi"`
	State              LoanState `json:"state"`
	AgreementLetterURL string    `json:"agreement_letter_url,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`

	Approval     *Approval     `json:"approval,omitempty"`
	Investments  []*Investment `json:"investments,omitempty"`
	Disbursement *Disbursement `json:"disbursement,omitempty"`
}

func NewLoan(borrowerID string, principalAmount, rate, roi float64) *Loan {
	now := time.Now()
	return &Loan{
		ID:              GenerateID(),
		BorrowerID:      borrowerID,
		PrincipalAmount: principalAmount,
		Rate:            rate,
		ROI:             roi,
		State:           LoanStateProposed,
		CreatedAt:       now,
		UpdatedAt:       now,
		Investments:     []*Investment{},
	}
}

func (l *Loan) CanApprove() error {
	if l.State != LoanStateProposed {
		return errors.New("loan must be in PROPOSED state to be approved")
	}
	return nil
}

func (l *Loan) Approve(approval *Approval) error {
	if err := l.CanApprove(); err != nil {
		return err
	}

	l.State = LoanStateApproved
	l.Approval = approval
	l.UpdatedAt = time.Now()
	return nil
}

func (l *Loan) CanAddInvestment(amount float64) error {
	if l.State != LoanStateApproved {
		return errors.New("loan must be in APPROVED state to add investments")
	}

	currentTotal := l.TotalInvestedAmount()
	if currentTotal+amount > l.PrincipalAmount {
		return errors.New("investment would exceed loan principal amount")
	}

	return nil
}

func (l *Loan) AddInvestment(investment *Investment) error {
	if err := l.CanAddInvestment(investment.Amount); err != nil {
		return err
	}

	l.Investments = append(l.Investments, investment)
	l.UpdatedAt = time.Now()

	if l.TotalInvestedAmount() == l.PrincipalAmount {
		l.State = LoanStateInvested
	}

	return nil
}

func (l *Loan) TotalInvestedAmount() float64 {
	var total float64
	for _, investment := range l.Investments {
		total += investment.Amount
	}
	return total
}

func (l *Loan) CanDisburse() error {
	if l.State != LoanStateInvested {
		return errors.New("loan must be in INVESTED state to be disbursed")
	}
	return nil
}

func (l *Loan) Disburse(disbursement *Disbursement) error {
	if err := l.CanDisburse(); err != nil {
		return err
	}

	l.State = LoanStateDisbursed
	l.Disbursement = disbursement
	l.UpdatedAt = time.Now()
	return nil
}

func GenerateID() string {
	return "loan_" + util.GenerateUUID()
}
