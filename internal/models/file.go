package models

//easyjson:json
type DonwloadReq struct {
	UserID   uint64 `json:"user_id"`
	FileName string `json:"file_name"`
}
