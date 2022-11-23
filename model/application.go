package model

type Application struct {
	tableName struct{} `sql:"application" pg:",discard_unknown_columns"`

	Id             int64        `sql:"id"  json:"id"`
	RequestId      int64        `sql:"request_id" json:"requestId"`
	CheckedId      int64        `sql:"checked_id" json:"checkedId"`
	Person         string       `sql:"person" json:"person"`
	CustomerType   CustomerType `sql:"customer_type" json:"customerType"`
	CustomerName   string       `sql:"customer_name" json:"customerName"`
	FilePath       string       `sql:"file_path" json:"filePath"`
	CourtName      string       `sql:"court_name" json:"courtName"`
	JudgeName      string       `sql:"judge_name" json:"judgeName"`
	DecisionNumber string       `sql:"decision_number" json:"decisionNumber"`
	DecisionDate   string       `sql:"decision_date" json:"decisionDate"`
	IsChecked      bool         `sql:"is_checked" json:"isChecked"`
	Note           string       `sql:"note" json:"note"`
	Status         Status       `sql:"status" json:"status"`
	Deadline       string       `sql:"deadline" json:"deadline"`
	AssigneeId     int64        `sql:"assignee_id" json:"assigneeId"`
	Priority       Priority     `sql:"priority" json:"priority"`
	MailSent       bool         `sql:"mail_sent" json:"mailSent"`
	Comments       []Comment    `json:"comments"`
	AssigneeName   string       `sql:"assignee_name" json:"assigneeName"`
	BeginDate      string       `sql:"begin_date" json:"beginDate"`
	EndDate        string       `sql:"end_date" json:"endDate"`
	CreatedAt      string       `sql:"created_at" json:"createdAt"`
}

type CustomerType string

const (
	Person   CustomerType = "PERSON"
	Taxpayer              = "TAXPAYER"
)

type Status string

const (
	Received   Status = "RECEIVED"
	Inprogress        = "IN_PROGRESS"
	Sent              = "SENT"
	Hold              = "HOLD"
)

type Priority string

const (
	Standard Priority = "STANDARD"
	High              = "HIGH"
)

type ChangeStatusRequest struct {
	Status      Status `json:"status"`
	Description string `json:"description"`
}

type FilterInfo struct {
	Courts []string `json:"courts"`
	Judges []string `json:"judges"`
}
