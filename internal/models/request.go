package models

type RegReqType int8

const (
	ParentPassportVerification RegReqType = iota
)

func (r RegReqType) String() string {
	switch r {
	case ParentPassportVerification:
		return "Подтверждение паспорта родителя"
	}
	return "UNKNOWN"
}

//easyjson:json
type ParentPassportReq struct {
	Passport string `json:"passport"`
}

//easyjson:json
type ParentPassportReqFull struct {
	ID         uint64     `json:"id"`
	UserID     uint64     `json:"user_id"`
	ManagerID  uint64     `json:"manager_id"`
	Type       RegReqType `json:"type"`
	Status     string     `json:"status"`
	CreateTime uint64     `json:"create_time"`
	Message    string     `json:"message"`
}

//easyjson:json
type ParentPassportResp struct {
	Type       string `json:"type"`
	Status     string `json:"status"`
	CreateTime uint64 `json:"create_time"`
	Message    string `json:"message"`
}

//easyjson:json
type ParentPassportRespList []ParentPassportResp
