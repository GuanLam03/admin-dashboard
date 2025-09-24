package gmail

import (
	"fmt"
	"encoding/base64"
	   "sync"
	   "strings"
	   "time"
	"golang.org/x/sync/errgroup"
	"context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
"google.golang.org/api/option"
	"github.com/goravel/framework/facades"
    "github.com/goravel/framework/contracts/http"
    "goravel/app/models"

)

type GmailController struct{}

func NewGmailController() *GmailController {
	return &GmailController{}
}

// OAuth2 config
func getOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     facades.Config().Env("GOOGLE_CLIENT_ID").(string),
		ClientSecret: facades.Config().Env("GOOGLE_CLIENT_SECRET").(string),
		RedirectURL:  facades.Config().Env("GOOGLE_REDIRECT_URI").(string),
		Scopes: []string{
			gmail.GmailReadonlyScope,
			gmail.GmailSendScope,
		},
		Endpoint: google.Endpoint,
	}
}

// Step 1: Redirect user to Google
func (c *GmailController) RedirectToGoogle(ctx http.Context) http.Response {
	account := ctx.Request().Input("account")

	url := getOAuthConfig().AuthCodeURL(account, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	return ctx.Response().Redirect(http.StatusFound, url)
}

// Step 2: Handle callback
func (c *GmailController) HandleCallback(ctx http.Context) http.Response {
	code := ctx.Request().Query("code")


	token, err := getOAuthConfig().Exchange(context.Background(), code)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Failed to exchange token",
		})
	}

	// Create Gmail service with the new token
	config := getOAuthConfig()
	client := config.Client(context.Background(), token)

	srv, err := gmail.New(client)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create Gmail service",
		})
	}

	// Get Gmail profile (fetch the real email)
	profile, err := srv.Users.GetProfile("me").Do()
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch Gmail profile",
		})
	}

	// Save to DB
	err = facades.Orm().Query().Create(&models.GmailAccount{
		Email:        profile.EmailAddress, // ✅ real Gmail email
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       &token.Expiry,

	})
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Failed to save Gmail account",
		})
	}

	return ctx.Response().Json(http.StatusOK, map[string]interface{}{
		"message": "Mailbox connected successfully!",
		"account": profile.EmailAddress,
	
	})
}


func GetClientFromDB(accountEmail string) (*gmail.Service, error) {
	var account models.GmailAccount
	err := facades.Orm().Query().Where("email", accountEmail).First(&account)
	if err != nil {
		return nil, err
	}

	config := getOAuthConfig()
	token := &oauth2.Token{
		AccessToken:  account.AccessToken,
		RefreshToken: account.RefreshToken,
		Expiry:       *account.Expiry,
	}

	// Handle expiry automatically
	tokenSource := config.TokenSource(context.Background(), token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}

	// If refreshed → update DB
	if newToken.AccessToken != account.AccessToken {
		account.AccessToken = newToken.AccessToken
		account.Expiry = &newToken.Expiry
		facades.Orm().Query().Save(&account)
	}

	// Return Gmail service ready to use
	client := oauth2.NewClient(context.Background(), tokenSource)
	return gmail.NewService(context.Background(), option.WithHTTPClient(client))
}

// func (c *GmailController) ListMessages(ctx http.Context) http.Response {
//     email := ctx.Request().Query("email") // /gmail/messages?email=hgledgetech@gmail.com


//     srv, err := GetClientFromDB(email) // ✅ already *gmail.Service
//     if err != nil {
//         return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
//             "error": "Failed to get Gmail client",
//         })
//     }

//     // res, _ := srv.Users.Messages.List("me").MaxResults(10).Do()
// 	res, err := srv.Users.Messages.List("me").MaxResults(10).Do()
// 	if err != nil {
// 		facades.Log().Errorf("Gmail API error for %s: %v", email, err)
// 		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
// 			"error": err.Error(),
// 		})
// 	}


//     messages := []map[string]string{}
//     for _, m := range res.Messages {
//         msg, _ := srv.Users.Messages.Get("me", m.Id).Do()
//         messages = append(messages, map[string]string{
//             "id":      m.Id,
//             "snippet": msg.Snippet,
//         })
//     }

// 	facades.Log().Infof("Fetching Gmail messages for %s", email)
// 	if err != nil {
//     facades.Log().Errorf("Gmail error: %v", err)
//     return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
//         "error": err.Error(),
//     })
// }



//     return ctx.Response().Json(http.StatusOK, messages)
// }

// func (c *GmailController) ListMessages(ctx http.Context) http.Response {
// 	email := ctx.Request().Query("email") // /gmail/messages?email=hgledgetech@gmail.com

// 	// Get the Gmail client from DB
// 	srv, err := GetClientFromDB(email)
// 	if err != nil {
// 		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
// 			"error": "Failed to get Gmail client",
// 		})
// 	}

// 	// Fetch the list of message IDs (this is a blocking call, but not too slow)
// 	res, err := srv.Users.Messages.List("me").MaxResults(10).Do()
// 	if err != nil {
// 		facades.Log().Errorf("Gmail API error for %s: %v", email, err)
// 		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
// 			"error": err.Error(),
// 		})
// 	}

// 	// Prepare the error group for goroutines
// 	var g errgroup.Group
// 	messages := make([]map[string]any, len(res.Messages))

// 	// Loop through the messages and fetch them concurrently
// 	for i, m := range res.Messages {
// 		// Capture the index to ensure correct result placement
// 		i, m := i, m
// 		g.Go(func() error {
// 			// Fetch full message with headers
// 			msg, err := srv.Users.Messages.Get("me", m.Id).Format("metadata").Do()
// 			if err != nil {
// 				facades.Log().Errorf("Failed to fetch message %s: %v", m.Id, err)
// 				return err
// 			}

// 			var from, subject, date string
// 			for _, h := range msg.Payload.Headers {
// 				switch h.Name {
// 				case "From":
// 					from = h.Value
// 				case "Subject":
// 					subject = h.Value
// 				case "Date":
// 					date = h.Value
// 				}
// 			}

// 			// Unread status check
// 			unread := false
// 			for _, label := range msg.LabelIds {
// 				if label == "UNREAD" {
// 					unread = true
// 					break
// 				}
// 			}

// 			messages[i] = map[string]any{
// 				"id":      msg.Id,
// 				"from":    from,
// 				"subject": subject,
// 				"snippet": msg.Snippet,
// 				"date":    date,
// 				"unread":  unread,
// 			}
// 			return nil
// 		})

// 	}

// 	// Wait for all goroutines to finish, and handle errors if any
// 	if err := g.Wait(); err != nil {
// 		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
// 			"error": fmt.Sprintf("Failed to fetch some messages: %v", err),
// 		})
// 	}

// 	// Log successful message fetch
// 	facades.Log().Infof("Fetched %d Gmail messages for %s", len(messages), email)

// 	// Return the fetched messages as a JSON response
// 	return ctx.Response().Json(http.StatusOK, messages)
// }

//reply message will not get, only get the initial message
// func (c *GmailController) ListMessages(ctx http.Context) http.Response {
// 	email := ctx.Request().Query("email")
// 	srv, err := GetClientFromDB(email)
// 	if err != nil {
// 		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{"error": "Failed to get Gmail client"})
// 	}

// 	res, err := srv.Users.Messages.List("me").MaxResults(50).Do()
// 	if err != nil {
// 		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{"error": err.Error()})
// 	}

// 	var g errgroup.Group
// 	threadMap := make(map[string]*gmail.Message)
// 	mu := &sync.Mutex{}

// 	for _, m := range res.Messages {
// 		m := m // capture range variable
// 		g.Go(func() error {
// 			msg, err := srv.Users.Messages.Get("me", m.Id).Format("metadata").Do()
// 			if err != nil {
// 				facades.Log().Errorf("Failed to fetch message %s: %v", m.Id, err)
// 				return nil // ignore failed message
// 			}

// 			mu.Lock()
// 			defer mu.Unlock()

// 			// If thread not in map, or this msg is earlier, store it
// 			if existing, exists := threadMap[msg.ThreadId]; !exists || msg.InternalDate < existing.InternalDate {
// 				threadMap[msg.ThreadId] = msg
// 			}

// 			return nil
// 		})
// 	}

// 	if err := g.Wait(); err != nil {
// 		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to fetch messages: %v", err)})
// 	}

// 	// Build response
// 	messages := []map[string]any{}
// 	for _, msg := range threadMap {
// 		var from, subject, date string
// 		for _, h := range msg.Payload.Headers {
// 			switch h.Name {
// 			case "From":
// 				from = h.Value
// 			case "Subject":
// 				subject = h.Value
// 			case "Date":
// 				date = h.Value
// 			}
// 		}
// 		unread := false
// 		for _, label := range msg.LabelIds {
// 			if label == "UNREAD" {
// 				unread = true
// 				break
// 			}
// 		}
// 		messages = append(messages, map[string]any{
// 			"id":       msg.Id,
// 			"threadId": msg.ThreadId,
// 			"from":     from,
// 			"subject":  subject,
// 			"snippet":  msg.Snippet,
// 			"date":     date,
// 			"unread":   unread,
// 		})
// 	}

// 	return ctx.Response().Json(http.StatusOK, messages)
// }

//get thread initial messages
func (c *GmailController) ListMessages(ctx http.Context) http.Response {
    email := ctx.Request().Query("email")
	pageToken := ctx.Request().Query("pageToken") // for pagination
    srv, err := GetClientFromDB(email)
    if err != nil {
        return ctx.Response().Json(http.StatusInternalServerError, map[string]string{"error": "Failed to get Gmail client"})
    }

      // Request 50 threads per page
    call := srv.Users.Threads.List("me").MaxResults(50)
    if pageToken != "" {
        call = call.PageToken(pageToken)
    }

    res, err := call.Do()
    if err != nil {
        return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }

    var g errgroup.Group
    threadMap := make(map[string]*gmail.Thread)  // To store full thread information
    replyCountMap := make(map[string]int)       // To track the number of replies per thread
	mu := &sync.Mutex{}

    for _, t := range res.Threads {
        t := t // capture range variable
        g.Go(func() error {
            // Fetch the full thread details
            thread, err := srv.Users.Threads.Get("me", t.Id).Do()
            if err != nil {
                facades.Log().Errorf("Failed to fetch thread %s: %v", t.Id, err)
                return nil // ignore failed thread
            }

            mu.Lock()
            defer mu.Unlock()

            // Store the thread
            threadMap[t.Id] = thread

            // Count the number of replies in the thread (total messages minus 1 for the initial message)
            replyCountMap[t.Id] = len(thread.Messages) - 1

            return nil
        })
    }

    // Wait for all goroutines to finish
    if err := g.Wait(); err != nil {
        return ctx.Response().Json(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to fetch threads: %v", err)})
    }

    // Build the response, showing the initial message and number of replies
    messages := []map[string]any{}
    for threadId, thread := range threadMap {
        // Get the first message in the thread (typically the earliest one)
        firstMessage := thread.Messages[0]

        // Extract headers like From, Subject, Date from the first message
        var from, subject, date string
        for _, h := range firstMessage.Payload.Headers {
            switch h.Name {
            case "From":
                from = h.Value
            case "Subject":
                subject = h.Value
            case "Date":
                date = h.Value
            }
        }

        // Get the number of replies (stored in replyCountMap)
        replyCount := replyCountMap[threadId]

        // Build the response object for this thread
        messages = append(messages, map[string]any{
            "id":         firstMessage.Id,
            "threadId":   threadId,
            "from":       from,
            "subject":    subject,
            "snippet":    firstMessage.Snippet,
            "date":       date,
            "replyCount": replyCount,
        })

		

    }

	response := map[string]any{
		"messages":      messages,           
		"nextPageToken": res.NextPageToken,  
	}

    return ctx.Response().Json(http.StatusOK, response)
}


func (c *GmailController) ReadMessage(ctx http.Context) http.Response {
	email := ctx.Request().Query("email") // /gmail/message?id=xxxx&email=hgledgetech@gmail.com
	messageID := ctx.Request().Route("id")

	if messageID == "" || email == "" {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Missing required parameters (id, email)",
		})
	}

	// Get Gmail client
	srv, err := GetClientFromDB(email)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get Gmail client",
		})
	}

	// Fetch the single message first to get threadId
	msg, err := srv.Users.Messages.Get("me", messageID).Format("full").Do()
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to fetch message %s: %v", messageID, err),
		})
	}

	// Fetch the full thread
	thread, err := srv.Users.Threads.Get("me", msg.ThreadId).Format("full").Do()
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to fetch thread %s: %v", msg.ThreadId, err),
		})
	}

	// Prepare messages array
	messages := []map[string]any{}
	for _, m := range thread.Messages {
		from := getHeader(m.Payload.Headers, "From")
		subject := getHeader(m.Payload.Headers, "Subject")
		date := getHeader(m.Payload.Headers, "Date")
		body := extractMessageBody(m.Payload)

		messages = append(messages, map[string]any{
			"id":      m.Id,
			"from":    from,
			"subject": subject,
			"date":    date,
			"snippet": m.Snippet,
			"body":    body,
		})
	}

	// Build response
	result := map[string]any{
		"threadId": thread.Id,
		"messages": messages,
	}

	return ctx.Response().Json(http.StatusOK, result)
}

// helper to extract a specific header
func getHeader(headers []*gmail.MessagePartHeader, key string) string {
	for _, h := range headers {
		if h.Name == key {
			return h.Value
		}
	}
	return ""
}

// recursive body extraction (your existing function)
func extractMessageBody(payload *gmail.MessagePart) string {
	if payload == nil {
		return ""
	}

	if payload.Body != nil && payload.Body.Data != "" {
		data, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err == nil {
			return string(data)
		}
	}

	// Prefer HTML over plain text
	var plain, html string
	for _, part := range payload.Parts {
		if part.MimeType == "text/plain" {
			plain = extractMessageBody(part)
		} else if part.MimeType == "text/html" {
			html = extractMessageBody(part)
		} else if len(part.Parts) > 0 {
			sub := extractMessageBody(part)
			if html == "" {
				html = sub
			}
		}
	}

	if html != "" {
		return html
	}
	return plain
}



// func (c *GmailController) ReplyMessage(ctx http.Context) http.Response {
// 	messageID := ctx.Request().Route("id")
// 	email := ctx.Request().Input("email")
// 	body := ctx.Request().Input("body")

// 	if messageID == "" || email == "" || body == "" {
// 		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
// 			"error": "Missing required parameters",
// 		})
// 	}

// 	// Get Gmail client
// 	srv, err := GetClientFromDB(email)
// 	if err != nil {
// 		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
// 			"error": "Failed to get Gmail client",
// 		})
// 	}

// 	// Fetch original message to get threadId and recipient
// 	origMsg, err := srv.Users.Messages.Get("me", messageID).Format("full").Do()
// 	if err != nil {
// 		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
// 			"error": fmt.Sprintf("Failed to fetch original message: %v", err),
// 		})
// 	}

// 	// Get original "From" to reply to
// 	to := getHeader(origMsg.Payload.Headers, "From")
// 	subject := getHeader(origMsg.Payload.Headers, "Subject")
// 	if subject[:3] != "Re:" {
// 		subject = "Re: " + subject
// 	}

// 	// Build raw RFC2822 message
// 	raw := fmt.Sprintf("To: %s\r\nSubject: %s\r\nIn-Reply-To: %s\r\nReferences: %s\r\n\r\n%s",
// 		to, subject, origMsg.Id, origMsg.Id, body)

// 	// Encode in base64 URL encoding
// 	message := &gmail.Message{
// 		Raw: base64.URLEncoding.EncodeToString([]byte(raw)),
// 		ThreadId: origMsg.ThreadId,
// 	}

// 	_, err = srv.Users.Messages.Send("me", message).Do()
// 	if err != nil {
// 		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
// 			"error": fmt.Sprintf("Failed to send reply: %v", err),
// 		})
// 	}

// 	return ctx.Response().Json(http.StatusOK, map[string]string{
// 		"message": "Reply sent successfully",
// 	})
// }

// func (c *GmailController) ReplyMessage(ctx http.Context) http.Response {
// 	messageID := ctx.Request().Route("id")
// 	email := ctx.Request().Input("email")
// 	body := ctx.Request().Input("body")

// 	if messageID == "" || email == "" || body == "" {
// 		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
// 			"error": "Missing required parameters",
// 		})
// 	}

// 	srv, err := GetClientFromDB(email)
// 	if err != nil {
// 		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
// 			"error": "Failed to get Gmail client",
// 		})
// 	}

// 	origMsg, err := srv.Users.Messages.Get("me", messageID).Format("full").Do()
// 	if err != nil {
// 		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
// 			"error": fmt.Sprintf("Failed to fetch original message: %v", err),
// 		})
// 	}

// 	to := getHeader(origMsg.Payload.Headers, "From")
// 	subject := getHeader(origMsg.Payload.Headers, "Subject")
// 	if !strings.HasPrefix(subject, "Re:") {
// 		subject = "Re: " + subject
// 	}

// 	messageIDHeader := getHeader(origMsg.Payload.Headers, "Message-ID")
// 	if !strings.HasPrefix(messageIDHeader, "<") {
// 		messageIDHeader = "<" + messageIDHeader + ">"
// 	}

// 	// Use original message date (format it properly)
// 	timestamp := time.Unix(origMsg.InternalDate/1000, 0).Format("Mon, Jan 2, 2006 at 3:04PM")
// 	originalBody := getMessageBody(origMsg)

// 	// Build HTML body with quoted section
// 	quotedHTML := fmt.Sprintf(`
// <div dir="ltr">%s</div><br>
// <div class="gmail_quote">
//   <div dir="ltr" class="gmail_attr">On %s, %s wrote:</div>
//   <blockquote class="gmail_quote" style="margin:0 0 0 .8ex;border-left:1px solid #ccc;padding-left:1ex">
//     <div dir="ltr">%s</div>
//   </blockquote>
// </div>`,
// 		body, timestamp, to, originalBody)

// 	raw := fmt.Sprintf(
// 		"To: %s\r\nSubject: %s\r\nIn-Reply-To: %s\r\nReferences: %s\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s",
// 		to, subject, messageIDHeader, messageIDHeader, quotedHTML,
// 	)

// 	encoded := base64.URLEncoding.EncodeToString([]byte(raw))
// 	message := &gmail.Message{
// 		Raw:      encoded,
// 		ThreadId: origMsg.ThreadId,
// 	}

// 	_, err = srv.Users.Messages.Send("me", message).Do()
// 	if err != nil {
// 		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
// 			"error": fmt.Sprintf("Failed to send reply: %v", err),
// 		})
// 	}

// 	return ctx.Response().Json(http.StatusOK, map[string]string{
// 		"message": "Reply sent successfully",
// 	})
// }

// // Extract text/plain body if exists, else fallback to snippet
// func getMessageBody(msg *gmail.Message) string {
//     if msg.Payload != nil && len(msg.Payload.Parts) > 0 {
//         for _, part := range msg.Payload.Parts {
//             if part.MimeType == "text/plain" && part.Body != nil && part.Body.Data != "" {
//                 // Gmail’s body data is base64 URL safe encoded
//                 data, err := base64.URLEncoding.DecodeString(part.Body.Data)
//                 if err == nil {
//                     return string(data)
//                 }
//             }
//         }
//     }
//     // fallback
//     return msg.Snippet
// }



func getHTMLBody(msg *gmail.Message) string {
	if msg.Payload.MimeType == "text/html" && msg.Payload.Body != nil && msg.Payload.Body.Data != "" {
		data, err := base64.URLEncoding.DecodeString(msg.Payload.Body.Data)
		if err == nil {
			return string(data)
		}
	}

	for _, part := range msg.Payload.Parts {
		if part.MimeType == "text/html" && part.Body != nil && part.Body.Data != "" {
			data, err := base64.URLEncoding.DecodeString(part.Body.Data)
			if err == nil {
				return string(data)
			}
		}
	}

	return "<i>(no original content found)</i>"
}


// ReplyMessage handles replying to an email
func (c *GmailController) ReplyMessage(ctx http.Context) http.Response {
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