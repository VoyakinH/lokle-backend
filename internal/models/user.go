package models

//easyjson:json
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Role int8

const (
	ParentRole Role = iota
	ChildRole
	ManagerRole
	AdminRole
)

func (r Role) String() string {
	switch r {
	case ParentRole:
		return "PARENT"
	case ChildRole:
		return "CHILD"
	case ManagerRole:
		return "MANAGER"
	case AdminRole:
		return "ADMIN"
	}
	return "UNKNOWN"
}

type Stage int8

const (
	FirstStage Stage = iota + 1
	SecondStage
	ThirdStage
)

//easyjson:json
type User struct {
	ID            uint64 `json:"id"`
	Role          Role   `json:"role"`
	FirstName     string `json:"first_name"`
	SecondName    string `json:"second_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Password      string `json:"password"`
	Phone         string `json:"phone"`
}

//easyjson:json
type UserRes struct {
	ID            uint64 `json:"id"`
	Role          string `json:"role"`
	FirstName     string `json:"first_name"`
	SecondName    string `json:"second_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Phone         string `json:"phone"`
}

//easyjson:json
type UserResList []UserRes

//easyjson:json
type Parent struct {
	ID               uint64 `json:"id"`
	UserID           uint64 `json:"user_id"`
	Role             Role   `json:"role"`
	FirstName        string `json:"first_name"`
	SecondName       string `json:"second_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	EmailVerified    bool   `json:"email_verified"`
	Password         string `json:"password"`
	Phone            string `json:"phone"`
	Passport         string `json:"passport"`
	PassportVerified bool   `json:"passport_verified"`
	DirPath          string `json:"dir_path"`
}

//easyjson:json
type ParentRes struct {
	Passport         string `json:"pasport"`
	PassportVerified bool   `json:"passport_verified"`
	DirPath          string `json:"dir_path"`
}

//easyjson:json
type Child struct {
	ID                  uint64 `json:"id"`
	UserID              uint64 `json:"user_id"`
	Role                Role   `json:"role"`
	FirstName           string `json:"first_name"`
	SecondName          string `json:"second_name"`
	LastName            string `json:"last_name"`
	Email               string `json:"email"`
	EmailVerified       bool   `json:"email_verified"`
	Password            string `json:"password,omitempty"`
	Phone               string `json:"phone"`
	BirthDate           uint64 `json:"birth_date"`
	DoneStage           Stage  `json:"done_stage"`
	Passport            string `json:"passport,omitempty"`
	PlaceOfResidence    string `json:"place_of_residence"`
	PlaceOfRegistration string `json:"place_of_registration"`
	DirPath             string `json:"dir_path,omitempty"`
}

//easyjson:json
type ChildWithRegReq struct {
	Child  Child       `json:"child"`
	RegReq *RegReqResp `json:"reg_req,omitempty"`
}

//easyjson:json
type ChildWithRegReqList []ChildWithRegReq

//easyjson:json
type ChildFullRes struct {
	ID                  uint64 `json:"id"`
	Role                string `json:"role"`
	FirstName           string `json:"first_name"`
	SecondName          string `json:"second_name"`
	LastName            string `json:"last_name"`
	Email               string `json:"email"`
	EmailVerified       bool   `json:"email_verified"`
	Phone               string `json:"phone"`
	BirthDate           uint64 `json:"birth_date"`
	DoneStage           Stage  `json:"done_stage"`
	Passport            string `json:"pasport"`
	PlaceOfResidence    string `json:"place_of_residence"`
	PlaceOfRegistration string `json:"place_of_registration"`
	DirPath             string `json:"dir_path"`
}

//easyjson:json
type ChildRes struct {
	BirthDate           uint64 `json:"birth_date"`
	DoneStage           Stage  `json:"done_stage"`
	Passport            string `json:"pasport"`
	PlaceOfResidence    string `json:"place_of_residence"`
	PlaceOfRegistration string `json:"place_of_registration"`
	DirPath             string `json:"dir_path"`
}
