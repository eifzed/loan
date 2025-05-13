package domain

import (
	"errors"
	"time"
)

type Approval struct {
	LoanID           string    `json:"loan_id"`
	ProofPictureURL  string    `json:"proof_picture_url"`
	FieldValidatorID string    `json:"field_validator_id"`
	ApprovalDate     time.Time `json:"approval_date"`
}

func NewApproval(loanID, proofPictureURL, fieldValidatorID string, approvalDate time.Time) (*Approval, error) {
	if loanID == "" {
		return nil, errors.New("loan ID cannot be empty")
	}

	if proofPictureURL == "" {
		return nil, errors.New("proof picture URL cannot be empty")
	}

	if fieldValidatorID == "" {
		return nil, errors.New("field validator ID cannot be empty")
	}

	if approvalDate.IsZero() {
		return nil, errors.New("approval date cannot be empty")
	}

	return &Approval{
		LoanID:           loanID,
		ProofPictureURL:  proofPictureURL,
		FieldValidatorID: fieldValidatorID,
		ApprovalDate:     approvalDate,
	}, nil
}
