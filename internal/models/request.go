package models

type RegReqType int8

const (
	ParentPassportVerification RegReqType = iota
	ChildFirstStageForStudent
	ChildFirstStage
	ChildSecondStage
	ChildThirdStage
)

func (r RegReqType) String() string {
	switch r {
	case ParentPassportVerification:
		return "Подтверждение паспорта родителя"
	case ChildFirstStageForStudent:
		return "Регистрация ребенка (этап 1 / старый ученик)"
	case ChildFirstStage:
		return "Регистрация ребенка (этап 1 / новый ученик)"
	case ChildSecondStage:
		return "Регистрация ребенка (этап 2)"
	case ChildThirdStage:
		return "Регистрация ребенка (этап 3)"
	}
	return "UNKNOWN"
}

//easyjson:json
type ParentPassportReq struct {
	Passport string `json:"passport"`
}

//easyjson:json
type ChildFirstRegReq struct {
	Child     Child `json:"child"`
	IsStudent bool  `json:"is_student"`
}

//easyjson:json
type ChildSecondRegReq struct {
	Child        Child  `json:"child"`
	Relationship string `json:"relationship"`
}

//easyjson:json
type RegReqFull struct {
	ID         uint64     `json:"id"`
	UserID     uint64     `json:"user_id"`
	ManagerID  uint64     `json:"manager_id"`
	Type       RegReqType `json:"type"`
	Status     string     `json:"status"`
	CreateTime uint64     `json:"create_time"`
	Message    string     `json:"message"`
}

//easyjson:json
type RegReqResp struct {
	ID         uint64 `json:"id"`
	UserID     uint64 `json:"user_id"`
	Type       string `json:"type"`
	Status     string `json:"status"`
	CreateTime uint64 `json:"create_time"`
	Message    string `json:"message"`
}

//easyjson:json
type RegReqRespList []RegReqResp
