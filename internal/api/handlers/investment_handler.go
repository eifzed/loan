package handlers

import (
	"encoding/json"
	"loan/internal/domain"
	"loan/internal/service"
	"net/http"

	"github.com/gorilla/mux"
)

type InvestmentHandler struct {
	loanService *service.LoanService
}

func NewInvestmentHandler(loanService *service.LoanService) *InvestmentHandler {
	return &InvestmentHandler{
		loanService: loanService,
	}
}

type InvestmentRequest struct {
	InvestorID string  `json:"investor_id"`
	Amount     float64 `json:"amount"`
}

func (h *InvestmentHandler) AddInvestment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	loanID := vars["id"]

	var req InvestmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := domain.NewErrorResponse(http.StatusBadRequest, "Invalid request body")
		writeJSON(w, http.StatusBadRequest, response)
		return
	}

	investment, err := h.loanService.AddInvestment(
		r.Context(),
		loanID,
		req.InvestorID,
		req.Amount,
	)

	if err != nil {
		response := domain.NewErrorResponse(http.StatusBadRequest, err.Error())
		writeJSON(w, http.StatusBadRequest, response)
		return
	}

	response := domain.NewSuccessResponse(
		http.StatusCreated,
		"Investment added successfully",
		investment,
	)

	writeJSON(w, http.StatusCreated, response)
}

func (h *InvestmentHandler) GetInvestments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	loanID := vars["id"]

	investments, err := h.loanService.GetLoanInvestments(r.Context(), loanID)
	if err != nil {
		response := domain.NewErrorResponse(http.StatusInternalServerError, err.Error())
		writeJSON(w, http.StatusInternalServerError, response)
		return
	}

	loan, err := h.loanService.GetLoan(r.Context(), loanID)
	if err != nil {
		response := domain.NewErrorResponse(http.StatusInternalServerError, err.Error())
		writeJSON(w, http.StatusInternalServerError, response)
		return
	}

	var totalInvested float64
	for _, inv := range investments {
		totalInvested += inv.Amount
	}

	investmentSummary := domain.NewInvestmentSummary(
		investments,
		totalInvested,
		loan.PrincipalAmount,
	)

	response := domain.NewSuccessResponse(
		http.StatusOK,
		"Investments retrieved successfully",
		investmentSummary,
	)

	writeJSON(w, http.StatusOK, response)
}
