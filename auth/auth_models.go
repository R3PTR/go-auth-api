package auth

import (
	"time"
)

type User struct {
	Id                  string    `bson:"_id,omitempty"`
	FirstName           string    `bson:"firstName"`
	LastName            string    `bson:"lastName"`
	Username            string    `bson:"username"`
	Password            string    `bson:"password"`
	OneTimePassword     string    `bson:"oneTimePassword,omitempty"`
	Role                string    `bson:"role"`
	State               string    `bson:"state"`
	Personnelnumber     string    `bson:"personnelnumber,omitempty"`
	VacationDaysPerYear int       `bson:"vacationDaysPerYear"`
	TargetHoursPerWeek  float32   `bson:"targetHoursPerWeek"`
	MaximumHoursPerWeek float32   `bson:"MaximumHoursPerWeek,omitempty"`
	TotpSecret          string    `bson:"totpSecret,omitempty"`
	TotpActive          bool      `bson:"totpActive,omitempty"`
	BackupCodes         []string  `bson:"backupCodes,omitempty"`
	InsertedAt          time.Time `bson:"insertedAt"`
	UpdatedAt           time.Time `bson:"updatedAt"`
	ResetValidUntil     time.Time `bson:"resetValidUntil,omitempty"`
}

type UserOutputAll struct {
	Id                  string    `bson:"_id,omitempty"`
	Username            string    `bson:"username"`
	FirstName           string    `bson:"firstName"`
	LastName            string    `bson:"lastName"`
	Role                string    `bson:"role"`
	State               string    `bson:"state"`
	Personnelnumber     string    `bson:"personnelnumber,omitempty"`
	VacationDaysPerYear int       `bson:"vacationDaysPerYear"`
	TargetHoursPerWeek  float32   `bson:"targetHoursPerWeek"`
	MaximumHoursPerWeek float32   `bson:"MaximumHoursPerWeek,omitempty"`
	InsertedAt          time.Time `bson:"insertedAt"`
	UpdatedAt           time.Time `bson:"updatedAt"`
}

type tokenModel struct {
	UserId         string    `bson:"user_id"`
	Token          string    `bson:"token"`
	Requires2FA    bool      `bson:"requires2FA,omitempty"`
	TwoFAConfirmed bool      `bson:"TwoFAConfirmed,omitempty"`
	TokenType      string    `bson:"tokenType"`
	InsertedAt     time.Time `bson:"insertedAt"`
	UpdatedAt      time.Time `bson:"updatedAt"`
	Expires        time.Time `bson:"expires"`
}

type CreateUserRequest struct {
	Username            string  `json:"username"`
	FirstName           string  `json:"firstName"`
	LastName            string  `json:"lastName"`
	Role                string  `json:"role"`
	Personnelnumber     string  `json:"personnelnumber,omitempty"`
	VacationDaysPerYear int     `bson:"vacationDaysPerYear,omitempty"`
	TargetHoursPerWeek  float32 `bson:"targetHoursPerWeek,omitempty"`
	MaximumHoursPerWeek float32 `bson:"MaximumHoursPerWeek,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ActivateUserRequest struct {
	NewPassword string `json:"newPassword"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"newPassword"`
}

type ChangePasswordRequest struct {
	NewPassword string `json:"newPassword"`
}

type DeleteOtherUserRequest struct {
	UserId string `json:"id"`
}

type ForgotPasswordRequest struct {
	Username string `json:"username"`
}

type ActivateTOTPRequest struct {
	TOTP string `json:"totp"`
}

type UpdateOtherUserRequest struct {
	Id                  string  `json:"id,omitempty"`
	Username            string  `json:"username,omitempty"`
	FirstName           string  `json:"firstName,omitempty"`
	LastName            string  `json:"lastName,omitempty"`
	Role                string  `json:"role,omitempty"`
	Personnelnumber     string  `json:"personnelnumber,omitempty"`
	VacationDaysPerYear int     `bson:"vacationDaysPerYear,omitempty"`
	TargetHoursPerWeek  float32 `bson:"targetHoursPerWeek,omitempty"`
	MaximumHoursPerWeek float32 `bson:"MaximumHoursPerWeek,omitempty"`
}

type UpdateOwnUserRequest struct {
	Username  string `json:"username,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}
