package auth

import (
	"time"

	"github.com/pquerna/otp/totp"
)

// GenerateTOTP generates a TOTP token for a given key.
func GenerateTOTP() (string, error) {
	// Create an OTP configuration
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "YourApp",
		AccountName: "user@example.com",
	})
	if err != nil {
		return "", err
	}

	// Get the current TOTP code
	otpCode, err := totp.GenerateCode(key.Secret(), time.Now())
	if err != nil {
		return "", err
	}

	return otpCode, nil
}

// VerifyTOTP verifies a TOTP token against the expected value.
func VerifyTOTP(secret, token string) (bool, error) {
	// Parse the secret to create an OTP key
	key, err := totp.Parse(secret)
	if err != nil {
		return false, err
	}

	// Verify the TOTP token
	valid := totp.Validate(token, key.Secret())
	return valid, nil
}
