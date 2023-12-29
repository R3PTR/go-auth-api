package auth

import (
	"time"
)

type User struct {
	Id              string    `bson:"_id,omitempty"`
	Username        string    `bson:"username"`
	Password        string    `bson:"password"`
	OneTimePassword string    `bson:"oneTimePassword,omitempty"`
	Role            string    `bson:"role"`
	State           string    `bson:"state"`
	Personnelnumber string    `bson:"personnelnumber,omitempty"`
	TotpSecret      string    `bson:"totpSecret,omitempty"`
	TotpActive      bool      `bson:"totpActive,omitempty"`
	BackupCodes     []string  `bson:"backupCodes,omitempty"`
	InsertedAt      time.Time `bson:"insertedAt"`
	UpdatedAt       time.Time `bson:"updatedAt"`
	ResetValidUntil time.Time `bson:"resetValidUntil,omitempty"`
}

type tokenModel struct {
	UserId         string    `bson:"user_id"`
	Token          string    `bson:"token"`
	Requires2FA    bool      `bson:"requires2FA,omitempty"`
	TwoFAConfirmed bool      `bson:"TwoFAConfirmed,omitempty"`
	InsertedAt     time.Time `bson:"insertedAt"`
	UpdatedAt      time.Time `bson:"updatedAt"`
	Expires        time.Time `bson:"expires"`
}

type loginReturn struct {
	Username string `bson:"username"`
	Token    string `bson:"token"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ActivateUserRequest struct {
	OneTimePassword string `json:"oneTimePassword"`
	NewPassword     string `json:"newPassword"`
}

type ResetPasswordRequest struct {
	OneTimePassword string `json:"oneTimePassword"`
	NewPassword     string `json:"newPassword"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type DeleteOtherUserRequest struct {
	UsernameToDelete string `json:"usernameToDelete"`
}

type ForgotPasswordRequest struct {
	Username string `json:"username"`
}

type ActivateTOTPRequest struct {
	TOTP string `json:"totp"`
}
