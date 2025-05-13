package domain

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewSuccessResponse(code int, message string, data interface{}) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(code int, message string) *Response {
	return &Response{
		Code:    code,
		Message: message,
	}
}

type PaginatedResponse struct {
	Items    interface{} `json:"items"`
	Total    int         `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

func NewPaginatedResponse(items interface{}, total, page, pageSize int) *PaginatedResponse {
	return &PaginatedResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}

type InvestmentSummary struct {
	Investments     interface{} `json:"investments"`
	TotalInvested   float64     `json:"total_invested"`
	PrincipalAmount float64     `json:"principal_amount"`
}

func NewInvestmentSummary(investments interface{}, totalInvested, principalAmount float64) *InvestmentSummary {
	return &InvestmentSummary{
		Investments:     investments,
		TotalInvested:   totalInvested,
		PrincipalAmount: principalAmount,
	}
}
