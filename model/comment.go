package model

type Comment struct {
	tableName struct{} `sql:"comment" pg:",discard_unknown_columns"`

	Id            int64       `sql:"id" json:"id"`
	Commentator   string      `sql:"commentator" json:"commentator"`
	Description   string      `sql:"description" json:"description"`
	CommentType   CommentType `sql:"comment_type" json:"commentType"`
	CreatedAt     string      `sql:"created_at" json:"createdAt"`
	UpdatedAt     string      `sql:"updated_at" json:"updatedAt"`
	ApplicationId int64       `sql:"application_id" json:"-"`
}

type CommentType string

const (
	Internal CommentType = "INTERNAL"
	External             = "EXTERNAL"
)

type UserDto struct {
	Id           int64  `json:"id"`
	Username     string `json:"username"`
	Pincode      string `json:"pincode"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Department   string `json:"department"`
	Position     string `json:"position"`
	FlexUserName string `json:"flexUserName"`
	IsActive     bool   `json:"isActive"`
}
