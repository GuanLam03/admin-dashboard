package models

import "github.com/goravel/framework/database/orm"

type User struct {
    orm.Model
    Name     string `json:"name"`
    Email    string `json:"email" gorm:"unique"`
    Password string `json:"password"`

    TwoFactorSecret string `gorm:"column:two_factor_secret"`
    TwoFactorEnabled bool  `gorm:"column:two_factor_enabled"`
}



var UserErrorMessage = map[string]string{
    "not_found":         "User not found.",
    "create_failed":     "Failed to create user account.",
    "update_failed":     "Failed to update user account.",
    "delete_failed":     "Failed to delete user account.",
    "email_exists":      "This email is already registered.",
    "invalid_credentials": "Invalid email or password.",
    "unauthorized":      "Unauthorized access. Please log in again.",
    "validation_failed": "Invalid input. Please check the fields and try again",
	"invalid_request":   "Invalid request body. Please check your JSON format.",
    "current_password_incorrect": "Current password is incorrect",
    "password_mismatch":    "New passwords do not match",
    "internal_error":    "Something went wrong. Please try again later.",
}





var TwofaErrorMessage = map[string]string{
    "already_enabled":   "Two-factor authentication is already enabled.",
    "not_enabled":       "Two-factor authentication is not enabled.",
    "invalid_code":      "Invalid verification code.",
    "generate_failed":   "Failed to generate 2FA secret key.",
    "encrypt_failed":    "Failed to encrypt 2FA secret.",
    "decrypt_failed":    "Failed to decrypt 2FA secret.",
    "save_failed":       "Failed to save 2FA settings.",
    "qr_failed":         "Failed to generate QR code.",
    "internal_error":    "Something went wrong. Please try again later.",
}