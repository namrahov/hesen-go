package model

type ErrorResponse struct {
	Code   string `json:"code"`
	Status int    `json:"-"`
}
