package model

import (
	"time"

	"github.com/aicacia/pickup/app/repository"
)

type UserST struct {
	Id                        int32     `json:"id" validate:"required"`
	Email                     string    `json:"email"`
	Username                  string    `json:"username" validate:"required"`
	TermsOfServiceAcknowledge bool      `json:"terms_of_service_acknowledge" validate:"required"`
	Birthdate                 time.Time `json:"birthdate" validate:"required" format:"date-time"`
	Sex                       string    `json:"sex" validate:"required"`
	UpdatedAt                 time.Time `json:"updated_at" validate:"required" format:"date-time"`
	CreatedAt                 time.Time `json:"created_at" validate:"required" format:"date-time"`
} // @name User

func UserFromUserRow(userRow repository.UserRowST) UserST {
	return UserST{
		Id:                        userRow.Id,
		Email:                     userRow.Email,
		Username:                  userRow.Username,
		TermsOfServiceAcknowledge: userRow.TermsOfServiceAcknowledge,
		Birthdate:                 userRow.Birthdate,
		Sex:                       userRow.Sex,
		UpdatedAt:                 userRow.UpdatedAt,
		CreatedAt:                 userRow.CreatedAt,
	}
}

type CompleteUserST struct {
	Username                  string    `json:"username" validate:"required"`
	TermsOfServiceAcknowledge bool      `json:"terms_of_service_acknowledge" validate:"required"`
	Birthdate                 time.Time `json:"birthdate" validate:"required" format:"date-time"`
	Sex                       string    `json:"sex" validate:"required"`
} // @name CompleteUser

type UpdateUserST struct {
	Username  string    `json:"username" validate:"required"`
	Birthdate time.Time `json:"birthdate" validate:"required" format:"date-time"`
	Sex       string    `json:"sex" validate:"required"`
} // @name UpdateUser
