package models

type RegReqType int8

const (
	ParentPassportVerification RegReqType = iota + 1
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
type FixParentPassportReq struct {
	ReqID    uint64 `json:"req_id"`
	Passport string `json:"passport"`
}

//easyjson:json
type ChildFirstRegReq struct {
	Child     Child `json:"child"`
	IsStudent bool  `json:"is_student"`
}

//easyjson:json
type FixChildFirstRegReq struct {
	ReqID     uint64 `json:"req_id"`
	Child     Child  `json:"child"`
	IsStudent bool   `json:"is_student"`
}

//easyjson:json
type ChildSecondRegReq struct {
	Child        Child  `json:"child"`
	Relationship string `json:"relationship"`
}

//easyjson:json
type FixChildSecondRegReq struct {
	ReqID        uint64 `json:"req_id"`
	Child        Child  `json:"child"`
	Relationship string `json:"relationship"`
}

//easyjson:json
type ChildThirdRegReq struct {
	Child Child `json:"child"`
}

//easyjson:json
type FixChildThirdRegReq struct {
	ReqID uint64 `json:"req_id"`
	Child Child  `json:"child"`
}

//easyjson:json
type RegReqFull struct {
	ID         uint64     `json:"id"`
	UserID     uint64     `json:"user_id"`
	ManagerID  uint64     `json:"manager_id,omitempty"`
	Type       RegReqType `json:"type"`
	Status     string     `json:"status"`
	CreateTime uint64     `json:"create_time"`
	Message    string     `json:"message"`
}

//easyjson:json
type RegReqResp struct {
	ID         uint64 `json:"id"`
	UserID     uint64 `json:"user_id"`
	ManagerID  uint64 `json:"manager_id,omitempty"`
	Type       string `json:"type"`
	Status     string `json:"status"`
	CreateTime uint64 `json:"create_time"`
	Message    string `json:"message"`
}

//easyjson:json
type RegReqRespList []RegReqResp

//easyjson:json
type RegReqWithUser struct {
	ID          uint64     `json:"id"`
	User        User       `json:"user"`
	Manager     *User      `json:"manager,omitempty"`
	Type        RegReqType `json:"type"`
	Status      string     `json:"status"`
	TimeInQueue uint32     `json:"time_in_queue"`
	CreateTime  uint64     `json:"create_time"`
	Message     string     `json:"message"`
}

//easyjson:json
type RegReqWithUserResp struct {
	ID          uint64   `json:"id"`
	User        UserRes  `json:"user"`
	Manager     *UserRes `json:"manager,omitempty"`
	Type        string   `json:"type"`
	Status      string   `json:"status"`
	TimeInQueue uint32   `json:"time_in_queue"`
	CreateTime  uint64   `json:"create_time"`
	Message     string   `json:"message"`
}

//easyjson:json
type RegReqWithUserRespList []RegReqWithUserResp

//easyjson:json
type FailedReq struct {
	ReqId         uint64 `json:"req_id"`
	FailedMessage string `json:"failed_message"`
}
