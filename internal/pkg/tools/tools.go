package tools

import "github.com/VoyakinH/lokle_backend/internal/models"

func UserToUserRes(user models.User) models.UserRes {
	return models.UserRes{
		Role:          user.Role.String(),
		FirstName:     user.FirstName,
		SecondName:    user.SecondName,
		LastName:      user.LastName,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		Phone:         user.Phone,
	}
}

func ParentToParentRes(parent models.Parent) models.ParentRes {
	return models.ParentRes{
		Passport:         parent.Passport,
		PassportVerified: parent.PassportVerified,
		DirPath:          parent.DirPath,
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

func FullParentPassportReqToSimpleList(reqs []models.ParentPassportReqFull) models.ParentPassportRespList {
	var respList models.ParentPassportRespList
	for _, req := range reqs {
		respList = append(respList, models.ParentPassportResp{
			Type:       req.Type.String(),
			Status:     req.Status,
			CreateTime: req.CreateTime,
			Message:    req.Message,
		})
	}
	return respList
}
