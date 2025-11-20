package gmail

import (
	"encoding/base64"
	"fmt"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"google.golang.org/api/gmail/v1"
	"strings"
	"time"
	// "goravel/app/models"
	"goravel/app/messages"
)

type ReplyGmailController struct{}

func NewReplyGmailController() *ReplyGmailController {
	return &ReplyGmailController{}
}

// ReplyMessage handles replying to an email
func (c *ReplyGmailController) ReplyMessage(ctx http.Context) http.Response {
	messageID := ctx.Request().Route("id")
	email := ctx.Request().Input("email")
	body := ctx.Request().Input("body")

	// Validate inputs
	if messageID == "" || email == "" || body == "" {
		facades.Log().Warningf("Missing required parameters")
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": messages.GetError("invalid_request"),
		})
	}

	// Get Gmail client
	srv, err := GetClientFromDB(email)
	if err != nil {
		facades.Log().Errorf("Failed to get Gmail client for %s: %v", email, err)
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": messages.GetError("gmail_account_not_found"),
		})
	}

	// Fetch the original message to get threadId, subject, and other details
	origMsg, err := srv.Users.Messages.Get("me", messageID).Format("full").Do()
	if err != nil {
		facades.Log().Errorf("Failed to fetch original message for %s: %v", messageID, err)
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": messages.GetError("gmail_account_read_failed"),
		})
	}

	// Extract "From" and "Subject" headers
	to := getHeader(origMsg.Payload.Headers, "From")

	subject := getHeader(origMsg.Payload.Headers, "Subject")

	// Make sure the subject has "Re:" prefix for reply
	if !strings.HasPrefix(subject, "Re:") {
		subject = "Re: " + subject
	}

	// Message-ID for "In-Reply-To" and "References"
	messageIDHeader := getHeader(origMsg.Payload.Headers, "Message-ID")
	if !strings.HasPrefix(messageIDHeader, "<") {
		messageIDHeader = "<" + messageIDHeader + ">"
	}

	// Get the timestamp of the original message
	timestamp := time.Unix(origMsg.InternalDate/1000, 0).Format("Mon, Jan 2, 2006 at 3:04PM")

	// Get the original HTML body
	originalHTMLBody := getHTMLBody(origMsg)

	// Construct the reply body in HTML
	htmlReply := fmt.Sprintf(`
		<div dir="ltr">%s</div><br>
		<div class="gmail_quote">
		<div dir="ltr" class="gmail_attr">
			On %s, %s wrote:
		</div>
		<blockquote class="gmail_quote" style="margin:0 0 0 .8ex;border-left:1px solid #ccc;padding-left:1ex">
			%s
		</blockquote>
		</div>`,
		body, timestamp, to, originalHTMLBody)

	// Construct the raw MIME email with proper headers
	raw := fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\nIn-Reply-To: %s\r\nReferences: %s\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s",
		to, subject, messageIDHeader, messageIDHeader, htmlReply,
	)

	// Base64 encode the raw message
	encoded := base64.URLEncoding.EncodeToString([]byte(raw))

	// Create the Gmail message object
	message := &gmail.Message{
		Raw:      encoded,
		ThreadId: origMsg.ThreadId,
	}

	// Send the reply using Gmail API
	_, err = srv.Users.Messages.Send("me", message).Do()
	if err != nil {
		facades.Log().Errorf("Failed to send reply for %s: %v", email, err)
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": messages.GetError("gmail_account_send_failed"),
		})
	}

	// Return success response
	return ctx.Response().Json(http.StatusOK, map[string]string{
		"message": messages.GetSuccess("gmail_reply_sent"),
	})
}
