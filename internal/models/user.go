package models

//easyjson:json
type Credentials struct {
	Email    string `json:"email,intern"`
	Password string `json:"password"`
}
