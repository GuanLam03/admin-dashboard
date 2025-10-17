package gmail

import (
	"fmt"
	"encoding/base64"
	"strings"
	"time"
	"google.golang.org/api/gmail/v1"
    "github.com/goravel/framework/contracts/http"
  

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
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Missing required parameters",
		})
	}

	// Get Gmail client
	srv, err := GetClientFromDB(email)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get Gmail client",
		})
	}

	// Fetch the original message to get threadId, subject, and other details
	origMsg, err := srv.Users.Messages.Get("me", messageID).Format("full").Do()
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to fetch original message: %v", err),
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
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to send reply: %v", err),
		})
	}

	// Return success response
	return ctx.Response().Json(http.StatusOK, map[string]string{
		"message": "Reply sent successfully",
	})
}