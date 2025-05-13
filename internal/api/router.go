package api

import (
	"loan/internal/api/handlers"
	"loan/internal/api/middleware"
	"loan/internal/service"

	"github.com/gorilla/mux"
)

func SetupRouter(loanService *service.LoanService) *mux.Router {
	router := mux.NewRouter()

	// middlewares
	router.Use(middleware.Logger)
	router.Use(middleware.ErrorHandler)

	loanHandler := handlers.NewLoanHandler(loanService)
	approvalHandler := handlers.NewApprovalHandler(loanService)
	investmentHandler := handlers.NewInvestmentHandler(loanService)
	disbursementHandler := handlers.NewDisbursementHandler(loanService)

	api := router.PathPrefix("/api/v1").Subrouter()

	// Loan routes
	api.HandleFunc("/loans", loanHandler.CreateLoan).Methods("POST")
	api.HandleFunc("/loans/{id}", loanHandler.GetLoan).Methods("GET")
	api.HandleFunc("/loans", loanHandler.ListLoans).Methods("GET")

	// Approval routes
	api.HandleFunc("/loans/{id}/approve", approvalHandler.ApproveLoan).Methods("POST")

	// Investment routes
	api.HandleFunc("/loans/{id}/investments", investmentHandler.AddInvestment).Methods("POST")
	api.HandleFunc("/loans/{id}/investments", investmentHandler.GetInvestments).Methods("GET")

	// Disbursement routes
	api.HandleFunc("/loans/{id}/disburse", disbursementHandler.DisburseLoan).Methods("POST")

	return router
}
