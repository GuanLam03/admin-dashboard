package gmail

import (
	"fmt"
	"encoding/base64"
	"time"
	"google.golang.org/api/gmail/v1"
    "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	// "goravel/app/models"
	"goravel/app/messages"
  

)

type ForwardGmailController struct{}

func NewForwardGmailController() *ForwardGmailController {
	return &ForwardGmailController{}
}



func (c *ForwardGmailController) ForwardMessage(ctx http.Context) http.Response {
	email := ctx.Request().Input("email") // sender Gmail account
	to := ctx.Request().Input("to")       // recipient
	subject := ctx.Request().Input("subject")
	body := ctx.Request().Input("body")   // full HTML (already prepared by frontend)

	if email == "" || to == "" || body == "" {
		facades.Log().Warningf("Missing required parameters")
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": messages.GetError("validation.invalid_request"),
			
		})
	}

	// Gmail client
	srv, err := GetClientFromDB(email)
	if err != nil {
		facades.Log().Errorf("Failed to get Gmail client for %s: %v", email, err)
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": messages.GetError("validation.gmail_account_not_found"),
		})
	}

	// Build MIME message
	msgStr := []byte(
		fmt.Sprintf("To: %s\r\n", to) +
			fmt.Sprintf("Subject: %s\r\n", subject) +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			body,
	)

	var message gmail.Message
	message.Raw = base64.URLEncoding.EncodeToString(msgStr)

	_, err = srv.Users.Messages.Send("me", &message).Do()
	if err != nil {
		facades.Log().Errorf("Failed to forward message: %v", err)
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error":messages.GetError("validation.gmail_account_forward_failed"),
			
		})
	}

	return ctx.Response().Json(http.StatusOK, map[string]string{
		"status": "Message forwarded successfully",
		"time":   time.Now().Format(time.RFC3339),
	})
}