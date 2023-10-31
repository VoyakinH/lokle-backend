package tools

import "github.com/VoyakinH/lokle_backend/internal/models"

func UserToUserRes(user models.User) models.UserRes {
	return models.UserRes{
		ID:            user.ID,
		Role:          user.Role.String(),
		FirstName:     user.FirstName,
		SecondName:    user.SecondName,
		LastName:      user.LastName,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		Phone:         user.Phone,
	}
}

func UsersToUserResList(users []models.User) models.UserResList {
	var resp models.UserResList
	for _, user := range users {
		resp = append(resp, UserToUserRes(user))
	}
	return resp
}

func ParentToParentRes(parent models.Parent) models.ParentRes {
	return models.ParentRes{
		Passport:         parent.Passport,
		PassportVerified: parent.PassportVerified,
	}
}

func ChildToChildFullRes(child models.Child) models.ChildFullRes {
	return models.ChildFullRes{
		ID:                  child.ID,
		Role:                child.Role.String(),
		FirstName:           child.FirstName,
		SecondName:          child.SecondName,
		LastName:            child.LastName,
		Email:               child.Email,
		EmailVerified:       child.EmailVerified,
		Phone:               child.Phone,
		BirthDate:           child.BirthDate,
		DoneStage:           child.DoneStage,
		Passport:            child.Passport,
		PlaceOfResidence:    child.PlaceOfResidence,
		PlaceOfRegistration: child.PlaceOfRegistration,
		DirPath:             child.DirPath,
	}
}

func ChildToChildRes(child models.Child) models.ChildRes {
	return models.ChildRes{
		BirthDate:           child.BirthDate,
		DoneStage:           child.DoneStage,
		Passport:            child.Passport,
		PlaceOfResidence:    child.PlaceOfResidence,
		PlaceOfRegistration: child.PlaceOfRegistration,
		DirPath:             child.DirPath,
	}
}

func ChildToUser(child models.Child) models.User {
	return models.User{
		ID:            child.ID,
		Role:          child.Role,
		FirstName:     child.FirstName,
		SecondName:    child.SecondName,
		LastName:      child.LastName,
		Email:         child.Email,
		EmailVerified: child.EmailVerified,
		Password:      child.Password,
		Phone:         child.Phone,
	}
}

func FullRegReqToSimpleResp(req models.RegReqFull) models.RegReqResp {
	return models.RegReqResp{
		ID:         req.ID,
		UserID:     req.UserID,
		Type:       req.Type.String(),
		Status:     req.Status,
		CreateTime: req.CreateTime,
		Message:    req.Message,
	}
}

func FullRegReqToSimpleRespList(reqs []models.RegReqFull) models.RegReqRespList {
	var respList models.RegReqRespList
	for _, req := range reqs {
		respList = append(respList, models.RegReqResp{
			ID:         req.ID,
			UserID:     req.UserID,
			Type:       req.Type.String(),
			Status:     req.Status,
			CreateTime: req.CreateTime,
			Message:    req.Message,
		})
	}
	return respList
}

func RegReqsWithUserToRespList(reqs []models.RegReqWithUser) models.RegReqWithUserRespList {
	var respList models.RegReqWithUserRespList
	tempManager := &models.UserRes{}
	var tempManagerRes models.UserRes
	for _, req := range reqs {
		if req.Manager != nil {
			tempManagerRes = UserToUserRes(*req.Manager)
			tempManager = &tempManagerRes
		} else {
			tempManager = nil
		}
		respList = append(respList, models.RegReqWithUserResp{
			ID:          req.ID,
			User:        UserToUserRes(req.User),
			Manager:     tempManager,
			Type:        req.Type.String(),
			Status:      req.Status,
			CreateTime:  req.CreateTime,
			TimeInQueue: req.TimeInQueue,
			Message:     req.Message,
		})
	}
	return respList
}
