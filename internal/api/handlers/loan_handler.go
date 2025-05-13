package handlers

import (
	"encoding/json"
	"loan/internal/domain"
	"loan/internal/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type LoanHandler struct {
	loanService *service.LoanService
}

func NewLoanHandler(loanService *service.LoanService) *LoanHandler {
	return &LoanHandler{
		loanService: loanService,
	}
}

type CreateLoanRequest struct {
	BorrowerID      string  `json:"borrower_id"`
	PrincipalAmount float64 `json:"principal_amount"`
	Rate            float64 `json:"rate"`
	ROI             float64 `json:"roi"`
}

func (h *LoanHandler) CreateLoan(w http.ResponseWriter, r *http.Request) {
	var req CreateLoanRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := domain.NewErrorResponse(http.StatusBadRequest, "Invalid request body")
		writeJSON(w, http.StatusBadRequest, response)
		return
	}

	loan, err := h.loanService.CreateLoan(r.Context(), req.BorrowerID, req.PrincipalAmount, req.Rate, req.ROI)
	if err != nil {
		response := domain.NewErrorResponse(http.StatusBadRequest, err.Error())
		writeJSON(w, http.StatusBadRequest, response)
		return
	}

	response := domain.NewSuccessResponse(
		http.StatusCreated,
		"Loan created successfully",
		loan,
	)

	writeJSON(w, http.StatusCreated, response)
}

func (h *LoanHandler) GetLoan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	loan, err := h.loanService.GetLoan(r.Context(), id)
	if err != nil {
		response := domain.NewErrorResponse(http.StatusNotFound, err.Error())
		writeJSON(w, http.StatusNotFound, response)
		return
	}

	response := domain.NewSuccessResponse(
		http.StatusOK,
		"Loan retrieved successfully",
		loan,
	)

	writeJSON(w, http.StatusOK, response)
}

func (h *LoanHandler) ListLoans(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	filters := make(map[string]interface{})

	if state := query.Get("state"); state != "" {
		filters["state"] = domain.LoanState(state)
	}

	if borrowerID := query.Get("borrower_id"); borrowerID != "" {
		filters["borrower_id"] = borrowerID
	}

	page := 1
	pageSize := 10

	if p := query.Get("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if ps := query.Get("page_size"); ps != "" {
		if parsedPageSize, err := strconv.Atoi(ps); err == nil && parsedPageSize > 0 {
			pageSize = parsedPageSize
		}
	}

	loans, total, err := h.loanService.ListLoans(r.Context(), page, pageSize)
	if err != nil {
		response := domain.NewErrorResponse(http.StatusInternalServerError, err.Error())
		writeJSON(w, http.StatusInternalServerError, response)
		return
	}

	paginatedResponse := domain.NewPaginatedResponse(loans, total, page, pageSize)

	response := domain.NewSuccessResponse(
		http.StatusOK,
		"Loans retrieved successfully",
		paginatedResponse,
	)

	writeJSON(w, http.StatusOK, response)
}
