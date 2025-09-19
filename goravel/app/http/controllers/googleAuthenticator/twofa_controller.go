package googleAuthenticator

import (
   
    "goravel/app/models"

    "github.com/pquerna/otp/totp"
    "github.com/skip2/go-qrcode"

    "github.com/goravel/framework/contracts/http"
    "github.com/goravel/framework/facades"
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
            "error": "unauthorized",
        })
    }

    // if already enabled don’t allow regenerate
    if user.TwoFactorEnabled {
        return ctx.Response().Json(400, http.Json{"error": "2FA already enabled"})
    }

    key, err := totp.Generate(totp.GenerateOpts{
        Issuer:      "Admin Dashboard",
        AccountName: user.Email,
    })
    if err != nil {
        return ctx.Response().Json(500, http.Json{
            "error": "failed to generate secret",
        })
    }

    encrypted, err := facades.Crypt().EncryptString(key.Secret())
    if err != nil {
        return ctx.Response().Json(500, http.Json{
            "error": "failed to encrypt secret",
        })
    }

    user.TwoFactorSecret = encrypted
    if err := facades.Orm().Query().Save(&user); err != nil {
        return ctx.Response().Json(500, http.Json{
            "error": "failed to save user secret",
        })
    }

    png, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
    if err != nil {
        return ctx.Response().Json(500, http.Json{
            "error": "failed to generate QR",
        })
    }

    return ctx.Response().Data(200, "image/png", png)
    
}

func (c *TwoFAController) ConfirmEnable(ctx http.Context) http.Response {
    code := ctx.Request().Input("code")

    var user models.User
    if err := facades.Auth(ctx).User(&user); err != nil {
        return ctx.Response().Json(401, http.Json{
            "error": "unauthorized",
        })
    }

    secret, err := facades.Crypt().DecryptString(user.TwoFactorSecret)
    if err != nil {
        return ctx.Response().Json(500, http.Json{
            "error": "failed to decrypt secret",
        })
    }

    if !totp.Validate(code, secret) {
        return ctx.Response().Json(400, http.Json{
            "error": "invalid code",
        })
    }

    user.TwoFactorEnabled = true

    if err := facades.Orm().Query().Save(&user); err != nil {
        return ctx.Response().Json(500, http.Json{
            "error": "failed to save user",
        })
    }

    return ctx.Response().Json(200, http.Json{
        "message":        "2FA enabled successfully",
        
    })
}


func (c *TwoFAController) ConfirmDisable(ctx http.Context) http.Response {
    code := ctx.Request().Input("code")

    var user models.User
    if err := facades.Auth(ctx).User(&user); err != nil {
        return ctx.Response().Json(401, http.Json{
            "error": "unauthorized",
        })
    }

    if !user.TwoFactorEnabled || user.TwoFactorSecret == "" {
        return ctx.Response().Json(400, http.Json{
            "error": "2FA is not enabled",
        })
    }

    decryptedSecret, err := facades.Crypt().DecryptString(user.TwoFactorSecret)
    if err != nil {
        return ctx.Response().Json(500, http.Json{
            "error": "failed to decrypt secret",
        })
    }

    // validate OTP before disabling
    if !totp.Validate(code, decryptedSecret) {
        return ctx.Response().Json(400, http.Json{
            "error": "invalid code",
        })
    }

    // reset 2FA fields
    user.TwoFactorEnabled = false
    user.TwoFactorSecret = ""

    if err := facades.Orm().Query().Save(&user); err != nil {
        return ctx.Response().Json(500, http.Json{
            "error": "failed to update user",
        })
    }

    return ctx.Response().Json(200, http.Json{
        "message": "2FA disabled successfully",
    })
}
