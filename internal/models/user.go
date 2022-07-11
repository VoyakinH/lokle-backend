package models

//easyjson:json
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	IsParent bool   `json:"is_parent"`
}

type Session struct {
	Email    string
	IsParent bool
}

//easyjson:json
type Parent struct {
	ID            uint64 `json:"id"`
	FirstName     string `json:"first_name"`
	SecondName    string `json:"second_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Password      string `json:"password,omitempty"`
	Phone         string `json:"phone"`
}
