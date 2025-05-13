package handlers

import (
	"encoding/json"
	"loan/internal/domain"
	"loan/internal/service"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type ApprovalHandler struct {
	loanService *service.LoanService
}

func NewApprovalHandler(loanService *service.LoanService) *ApprovalHandler {
	return &ApprovalHandler{
		loanService: loanService,
	}
}

type ApprovalRequest struct {
	ProofPictureURL  string `json:"proof_picture_url"`
	FieldValidatorID string `json:"field_validator_id"`
	ApprovalDate     string `json:"approval_date"` // Format: YYYY-MM-DD
}

func (h *ApprovalHandler) ApproveLoan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	loanID := vars["id"]

	var req ApprovalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := domain.NewErrorResponse(http.StatusBadRequest, "Invalid request body")
		writeJSON(w, http.StatusBadRequest, response)
		return
	}

	approvalDate, err := time.Parse("2006-01-02", req.ApprovalDate)
	if err != nil {
		response := domain.NewErrorResponse(http.StatusBadRequest, "Invalid approval date format. Use YYYY-MM-DD")
		writeJSON(w, http.StatusBadRequest, response)
		return
	}

	loan, err := h.loanService.ApproveLoan(
		r.Context(),
		loanID,
		req.ProofPictureURL,
		req.FieldValidatorID,
		approvalDate,
	)

	if err != nil {
		response := domain.NewErrorResponse(http.StatusBadRequest, err.Error())
		writeJSON(w, http.StatusBadRequest, response)
		return
	}

	response := domain.NewSuccessResponse(
		http.StatusOK,
		"Loan approved successfully",
		loan,
	)

	writeJSON(w, http.StatusOK, response)
}
