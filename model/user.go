package model

type User struct {
	tableName struct{} `sql:"userr" pg:",discard_unknown_columns"`

	Id        string `sql:"id"  json:"id"`
	UserName  string `sql:"user_name" json:"userName"`
	FirstName string `sql:"first_name" json:"firstName"`
	LastName  string `sql:"last_name" json:"lastName"`
	Password  []byte `sql:"password" json:"-"`
	CreatedAt string `sql:"created_at" json:"createdAt"`
	UpdatedAt string `sql:"updated_at" json:"updatedAt"`
}

type UserRegister struct {
	Id        string `json:"id"`
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

type Session struct {
	tableName struct{} `sql:"sessionn" pg:",discard_unknown_columns"`

	Id        int64  `sql:"id"  json:"id"`
	SessionId string `sql:"session_id" json:"sessionId"`
	UserId    string `sql:"user_id" json:"userId"`
}
