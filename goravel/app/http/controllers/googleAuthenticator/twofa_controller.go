package googleAuthenticator

import (
	"goravel/app/models"

	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/messages"
)

type TwoFAController struct{}

func NewTwoFAController() *TwoFAController {
	return &TwoFAController{}
}

// Generate a secret and return QR PNG (authenticated route)
func (c *TwoFAController) GenerateQRCode(ctx http.Context) http.Response {
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return ctx.Response().Json(401, http.Json{
			"error": messages.GetError("validation.unauthorized"),
		})
	}

	// if already enabled don’t allow regenerate
	if user.TwoFactorEnabled {
		return ctx.Response().Json(400, http.Json{"error": messages.GetError("validation.twofa_already_enabled")})
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Admin Dashboard",
		AccountName: user.Email,
	})
	if err != nil {
		return ctx.Response().Json(500, http.Json{
			"error": messages.GetError("validation.twofa_generate_failed"),
		})
	}

	encrypted, err := facades.Crypt().EncryptString(key.Secret())
	if err != nil {
		return ctx.Response().Json(500, http.Json{
			"error": messages.GetError("validation.twofa_encrypt_failed"),
		})
	}

	user.TwoFactorSecret = encrypted
	if err := facades.Orm().Query().Save(&user); err != nil {
		return ctx.Response().Json(500, http.Json{
			"error": messages.GetError("validation.twofa_save_failed"),
		})
	}

	png, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
	if err != nil {
		return ctx.Response().Json(500, http.Json{
			"error": messages.GetError("validation.twofa_qr_failed"),
		})
	}

	return ctx.Response().Data(200, "image/png", png)

}

func (c *TwoFAController) ConfirmEnable(ctx http.Context) http.Response {
	code := ctx.Request().Input("code")

	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return ctx.Response().Json(401, http.Json{
			"error": messages.GetError("validation.unauthorized"),
		})
	}

	secret, err := facades.Crypt().DecryptString(user.TwoFactorSecret)
	if err != nil {
		return ctx.Response().Json(500, http.Json{
			"error": messages.GetError("validation.twofa_decrypt_failed"),
		})
	}

	if !totp.Validate(code, secret) {
		return ctx.Response().Json(400, http.Json{
			"error": messages.GetError("validation.twofa_invalid_code"),
		})
	}

	user.TwoFactorEnabled = true

	if err := facades.Orm().Query().Save(&user); err != nil {
		return ctx.Response().Json(500, http.Json{
			"error": messages.GetError("validation.twofa_save_failed"),
		})
	}

	return ctx.Response().Json(200, http.Json{
		"message": messages.GetSuccess("twofa_enabled"),
	})
}

func (c *TwoFAController) ConfirmDisable(ctx http.Context) http.Response {
	code := ctx.Request().Input("code")

	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return ctx.Response().Json(401, http.Json{
			"error": messages.GetError("validation.unauthorized"),
		})
	}

	if !user.TwoFactorEnabled || user.TwoFactorSecret == "" {
		return ctx.Response().Json(400, http.Json{
			"error": messages.GetError("validation.twofa_not_enabled"),
		})
	}

	decryptedSecret, err := facades.Crypt().DecryptString(user.TwoFactorSecret)
	if err != nil {
		return ctx.Response().Json(500, http.Json{
			"error": messages.GetError("validation.twofa_decrypt_failed"),
		})
	}

	// validate OTP before disabling
	if !totp.Validate(code, decryptedSecret) {
		return ctx.Response().Json(400, http.Json{
			"error": messages.GetError("validation.twofa_invalid_code"),
		})
	}

	// reset 2FA fields
	user.TwoFactorEnabled = false
	user.TwoFactorSecret = ""

	if err := facades.Orm().Query().Save(&user); err != nil {
		return ctx.Response().Json(500, http.Json{
			"error": messages.GetError("validation.twofa_save_failed"),
		})
	}

	return ctx.Response().Json(200, http.Json{
		"message": messages.GetSuccess("twofa_disabled"),
	})
}
