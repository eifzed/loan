package domain

import (
	"errors"
	"time"
)

type Disbursement struct {
	LoanID               string    `json:"loan_id"`
	AgreementDocumentURL string    `json:"agreement_document_url"`
	FieldOfficerID       string    `json:"field_officer_id"`
	DisbursementDate     time.Time `json:"disbursement_date"`
}

func NewDisbursement(loanID, agreementDocumentURL, fieldOfficerID string, disbursementDate time.Time) (*Disbursement, error) {
	if loanID == "" {
		return nil, errors.New("loan ID cannot be empty")
	}

	if agreementDocumentURL == "" {
		return nil, errors.New("agreement document URL cannot be empty")
	}

	if fieldOfficerID == "" {
		return nil, errors.New("field officer ID cannot be empty")
	}

	if disbursementDate.IsZero() {
		return nil, errors.New("disbursement date cannot be empty")
	}

	return &Disbursement{
		LoanID:               loanID,
		AgreementDocumentURL: agreementDocumentURL,
		FieldOfficerID:       fieldOfficerID,
		DisbursementDate:     disbursementDate,
	}, nil
}
