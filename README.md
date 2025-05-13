# Loan Service API Design

## Overview
This design outlines a RESTful API for a loan management system that handles the full lifecycle of a loan from proposal to disbursement.

## Domain Model

### Loan States
- `PROPOSED` - Initial state when loan is created
- `APPROVED` - After staff approval with required documentation
- `INVESTED` - When total investment equals loan principal
- `DISBURSED` - When loan is given to borrower with signed agreement

### Core Entities

#### Loan
- ID (unique identifier)
- BorrowerID (borrower identity number)
- PrincipalAmount (loan amount)
- Rate (defines total interest borrower will pay)
- ROI (return on investment for investors)
- State (current loan state)
- AgreementLetterURL (link to generated agreement letter)
- CreatedAt (timestamp)
- UpdatedAt (timestamp)

#### Approval
- LoanID (reference to the loan)
- ProofPictureURL (evidence of field validator visit)
- FieldValidatorID (employee who validated)
- ApprovalDate (date of approval)

#### Investment
- ID (unique identifier)
- LoanID (reference to the loan)
- InvestorID (reference to investor)
- Amount (invested amount)
- InvestedAt (timestamp)

#### Disbursement
- LoanID (reference to the loan)
- AgreementDocumentURL (signed loan agreement)
- FieldOfficerID (employee who handled disbursement)
- DisbursementDate (date of disbursement)

## API Endpoints

### Loans

#### POST /api/v1/loans
Creates a new loan in the PROPOSED state.

Request:
```json
{
  "borrower_id": "string",
  "principal_amount": float,
  "rate": float,
  "roi": float
}
```

Response:
```json
{
  "id": "string",
  "borrower_id": "string",
  "principal_amount": float,
  "rate": float,
  "roi": float,
  "state": "PROPOSED",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

#### GET /api/v1/loans/{id}
Retrieves a specific loan by ID.

Response:
```json
{
  "id": "string",
  "borrower_id": "string",
  "principal_amount": float,
  "rate": float,
  "roi": float,
  "state": "string",
  "agreement_letter_url": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

#### GET /api/v1/loans
List all loans with optional filtering.

Query Parameters:
- state (optional): Filter by loan state
- borrower_id (optional): Filter by borrower ID

Response:
```json
{
  "loans": [
    {
      "id": "string",
      "borrower_id": "string",
      "principal_amount": float,
      "rate": float,
      "roi": float,
      "state": "string",
      "created_at": "timestamp",
      "updated_at": "timestamp"
    }
  ],
  "total": integer,
  "page": integer,
  "page_size": integer
}
```

### Loan Approval

#### POST /api/v1/loans/{id}/approve
Approves a loan, changing state from PROPOSED to APPROVED.

Request:
```json
{
  "proof_picture_url": "string",
  "field_validator_id": "string",
  "approval_date": "date"
}
```

Response:
```json
{
  "id": "string",
  "state": "APPROVED",
  "updated_at": "timestamp",
  "approval": {
    "proof_picture_url": "string",
    "field_validator_id": "string",
    "approval_date": "date"
  }
}
```

### Investments

#### POST /api/v1/loans/{id}/investments
Adds an investment to a loan.

Request:
```json
{
  "investor_id": "string",
  "amount": float
}
```

Response:
```json
{
  "id": "string",
  "loan_id": "string",
  "investor_id": "string",
  "amount": float,
  "invested_at": "timestamp"
}
```

#### GET /api/v1/loans/{id}/investments
Lists all investments for a loan.

Response:
```json
{
  "investments": [
    {
      "id": "string",
      "loan_id": "string",
      "investor_id": "string",
      "amount": float,
      "invested_at": "timestamp"
    }
  ],
  "total_invested": float,
  "principal_amount": float
}
```

### Loan Disbursement

#### POST /api/v1/loans/{id}/disburse
Disburses a loan, changing state from INVESTED to DISBURSED.

Request:
```json
{
  "agreement_document_url": "string",
  "field_officer_id": "string",
  "disbursement_date": "date"
}
```

Response:
```json
{
  "id": "string",
  "state": "DISBURSED",
  "updated_at": "timestamp",
  "disbursement": {
    "agreement_document_url": "string",
    "field_officer_id": "string",
    "disbursement_date": "date"
  }
}
```

## Business Rules Implementation

1. Loans can only move forward in state (PROPOSED → APPROVED → INVESTED → DISBURSED)
2. Approval requires proof picture, field validator ID, and approval date
3. Investment total cannot exceed loan principal amount
4. When total investment equals principal amount, loan state changes to INVESTED
5. Disbursement requires agreement document, field officer ID, and disbursement date
6. When loan becomes INVESTED, email notifications are sent to all investors

## Assumptions

1. Authentication and authorization mechanisms are handled by an external system
2. The API assumes valid input formats; detailed input validation errors will be provided
3. File uploads for documents and images are handled by a separate service
4. Email notification service is available as a dependency
5. Agreement letter generation is handled by a separate service