package models

import (
	"time"

	"github.com/Lucasdesouzat/ReservasCondominio/pkg/constants"
)

type User struct {
	ID             int                  `json:"id" db:"id"`
	FirstName      string               `json:"first_name" db:"first_name"`
	LastName       string               `json:"last_name" db:"last_name"`
	CPF            string               `json:"cpf" db:"cpf"`
	BirthDate      constants.CustomDate `json:"birth_date" db:"birth_date" time_format:"2006-01-02"`
	ProfilePicture *string              `json:"profile_picture" db:"profile_picture"`
	Phone1         string               `json:"phone_1" db:"phone_1"`
	Phone2         *string              `json:"phone_2,omitempty" db:"phone_2"`
	Email          string               `json:"email" db:"email"`
	Password       string               `json:"password,omitempty" db:"password"`
	Role           constants.UserRole   `json:"role" db:"role"`
	Status         constants.UserStatus `json:"status" db:"status"`
	AuthMethod     string               `json:"auth_method" db:"auth_method"`
	CreatedAt      time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at" db:"updated_at"`
}
