package model

const (
	HeaderKeyCustomerID = "DP-Customer-ID"
	HeaderKeyUserID     = "DP-User-ID"
	HeaderKeyUserAgent  = "User-Agent"
	HeaderKeyUserIP     = "X-Forwarded-For"
	HeaderKeyRequestID  = "requestid"
)

const (
	LoggerKeyRequestID  = "REQUEST_ID"
	LoggerKeyOperation  = "OPERATION"
	LoggerKeyCustomerID = "CUSTOMER_ID"
	LoggerKeyUserID     = "USER_ID"
	LoggerKeyUserIP     = "USER_IP"
	LoggerKeyUserAgent  = "USER_AGENT"
	ContextLogger       = "contextLogger"
	ContextHeader       = "contextHeader"
)

const (
	UserIdHeader                = "User-Id"
	GenerateReportPermissionKey = "lscView"
	ContentTypeString           = "Content-Type"
	ContentDispositionString    = "Content-Disposition"
	JSONType                    = "application/json"
	ExcelType                   = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	POST                        = "POST"
	GET                         = "GET"
	Sheet                       = "Sheet1"
	AttachmentFilename          = "attachment; filename=loan_risks_report.xlsx"
)

const (
	FirstSheet  = "Sheet1"
	SecondSheet = "Sheet2"
)

const Exception = "error.ecourt"
