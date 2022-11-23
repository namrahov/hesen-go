package model

type ErrorResponse struct {
	Code   string `json:"code"`
	Status int    `json:"-"`
}

var (
	UnexpectedError        = ErrorResponse{Code: "UNEXPECTED_EXCEPTION", Status: 500}
	CantParseDataError     = ErrorResponse{Code: "CANT_PARSE_DATA", Status: 400}
	CantMarshalObjectError = ErrorResponse{Code: "CANT_MARSHALL_OBJECT", Status: 400}
	AdminClientError       = ErrorResponse{Code: "ADMIN_CLIENT_FAILED", Status: 503}
	LoanRiskClientError    = ErrorResponse{Code: "LOAN_RISK_CLIENT_FAILED", Status: 503}
	InvalidHeaderError     = ErrorResponse{Code: "INVALID_HEADER", Status: 400}
	AccessDeniedError      = ErrorResponse{Code: "ACCESS_DENIED", Status: 403}
	InvalidTypeError       = ErrorResponse{Code: "INVALID_TYPE", Status: 400}
)
