package handlers

import (
	"encoding/json"
	"loan/internal/domain"
	"loan/internal/service"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type DisbursementHandler struct {
	loanService *service.LoanService
}

func NewDisbursementHandler(loanService *service.LoanService) *DisbursementHandler {
	return &DisbursementHandler{
		loanService: loanService,
	}
}

type DisbursementRequest struct {
	AgreementDocumentURL string `json:"agreement_document_url"`
	FieldOfficerID       string `json:"field_officer_id"`
	DisbursementDate     string `json:"disbursement_date"` // Format: YYYY-MM-DD
}

func (h *DisbursementHandler) DisburseLoan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	loanID := vars["id"]

	var req DisbursementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := domain.NewErrorResponse(http.StatusBadRequest, "Invalid request body")
		writeJSON(w, http.StatusBadRequest, response)
		return
	}

	disbursementDate, err := time.Parse("2006-01-02", req.DisbursementDate)
	if err != nil {
		response := domain.NewErrorResponse(http.StatusBadRequest, "Invalid disbursement date format. Use YYYY-MM-DD")
		writeJSON(w, http.StatusBadRequest, response)
		return
	}

	loan, err := h.loanService.DisburseLoan(
		r.Context(),
		loanID,
		req.AgreementDocumentURL,
		req.FieldOfficerID,
		disbursementDate,
	)

	if err != nil {
		response := domain.NewErrorResponse(http.StatusBadRequest, err.Error())
		writeJSON(w, http.StatusBadRequest, response)
		return
	}

	response := domain.NewSuccessResponse(
		http.StatusOK,
		"Loan disbursed successfully",
		loan,
	)

	writeJSON(w, http.StatusOK, response)
}
