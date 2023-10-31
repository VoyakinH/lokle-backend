package models

//easyjson:json
type DonwloadReq struct {
	UserID       uint64   `json:"user_id"`
	FileName     []string `json:"file_name"`
	ResponseType string   `json:"response_type"`
}

//easyjson:json
type DeleteReq struct {
	UserID   uint64 `json:"user_id"`
	FileName string `json:"file_name"`
}

//easyjson:json
type DonwloadResp struct {
	Files []FileStruct `json:"files"`
}

//easyjson:json
type FileStruct struct {
	File string `json:"file"`
	Type string `json:"type"`
}
