package repository

import (
	"strings"
	"time"

	"github.com/aicacia/pickup/app/util"
	"github.com/jmoiron/sqlx"
)

type UserRowST struct {
	Id                        int32     `db:"id"`
	Username                  string    `db:"username"`
	Email                     string    `db:"email"`
	TermsOfServiceAcknowledge bool      `db:"terms_of_service_acknowledge"`
	Birthdate                 time.Time `db:"birthdate"`
	Sex                       string    `db:"sex"`
	UpdatedAt                 time.Time `db:"updated_at"`
	CreatedAt                 time.Time `db:"created_at"`
}

func GetUserById(id int32) (*UserRowST, error) {
	return GetOptional[UserRowST](`SELECT u.*
		FROM users u
		WHERE u.id = $1
		LIMIT 1;`,
		id)
}

func GetUserByEmail(email string) (*UserRowST, error) {
	return GetOptional[UserRowST](`SELECT u.*
		FROM users u
		WHERE u.email = $1
		LIMIT 1;`,
		email)
}

func GetUserByUsername(username string) (*UserRowST, error) {
	return GetOptional[UserRowST](`SELECT u.*
		FROM users u
		WHERE u.username = $1
		LIMIT 1;`,
		username)
}

func CreateUserFromEmail(email string) (UserRowST, error) {
	return Transaction[UserRowST](func(tx *sqlx.Tx) (UserRowST, error) {
		var result UserRowST

		emailParts := strings.Split(email, "@")
		username := emailParts[0]
		for {
			user, err := GetUserByUsername(username)
			if err != nil {
				return result, err
			}
			if user == nil {
				break
			}
			hex, err := util.GenerateRandomHex(2)
			if err != nil {
				return result, err
			}
			username += hex
		}

		userRow := tx.QueryRowx(`INSERT INTO users (username, email)
			VALUES ($1, $2)
			RETURNING *;`,
			username, email)
		if userRow.Err() != nil {
			return result, userRow.Err()
		}
		err := userRow.StructScan(&result)
		if err != nil {
			return result, err
		}

		return result, nil
	})
}

type UpdateUserST struct {
	Username                  *string    `json:"username"`
	TermsOfServiceAcknowledge *bool      `json:"terms_of_service_acknowledge"`
	Birthdate                 *time.Time `json:"birthdate" format:"date-time"`
	Sex                       *string    `json:"sex"`
}

func UpdateUser(userId int32, update UpdateUserST) (*UserRowST, error) {
	return GetOptional[UserRowST](`UPDATE users
	SET username = COALESCE($2, username),
		terms_of_service_acknowledge = COALESCE($3, terms_of_service_acknowledge),
		birthdate = COALESCE($4, birthdate),
		sex = COALESCE($5, sex)
	WHERE id = $1
	RETURNING *;`,
		userId, update.Username, update.TermsOfServiceAcknowledge, update.Birthdate, update.Sex)
}
