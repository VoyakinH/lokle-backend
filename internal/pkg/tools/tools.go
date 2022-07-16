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
