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
			gmail.GmailModifyScope,
			gmail.GmailSendScope,
		},
		Endpoint: google.Endpoint,
	}
}

func (c *GmailController) RedirectToGoogle(ctx http.Context) http.Response {
    team := ctx.Request().Query("team") // technical / support / etc.

    url := getOAuthConfig().AuthCodeURL(team, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
    return ctx.Response().Json(http.StatusOK, url)
}

func (c *GmailController) HandleCallback(ctx http.Context) http.Response {
    code := ctx.Request().Query("code")
    team := ctx.Request().Query("state") // <-- from AuthCodeURL

    token, err := getOAuthConfig().Exchange(context.Background(), code)
    if err != nil {
        return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
            "error": models.GmailAccountErrorMessage["auth_failed"],
        })
    }

    config := getOAuthConfig()
    client := config.Client(context.Background(), token)

    srv, err := gmail.New(client)
    if err != nil {
        return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
            "error": models.GmailAccountErrorMessage["internal_error"],
        })
    }

    profile, err := srv.Users.GetProfile("me").Do()
    if err != nil {
        return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
            "error": models.GmailAccountErrorMessage["fetch_failed"],
        })
    }

    //  Check if this Gmail is already linked to another department
    var existing models.GmailAccount
    err = facades.Orm().Query().Where("email", profile.EmailAddress).First(&existing)
    if err == nil && existing.ID != 0 {
        return ctx.Response().Json(http.StatusBadRequest, map[string]string{
            "error": models.GmailAccountErrorMessage["already_linked"],
        })
    }
    
    // Save new account
    err = facades.Orm().Query().Create(&models.GmailAccount{
        Email:        profile.EmailAddress,
        Team:         &team,
        AccessToken:  token.AccessToken,
        RefreshToken: token.RefreshToken,
        Expiry:       &token.Expiry,
    })
    if err != nil {
        return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
            "error": models.GmailAccountErrorMessage["create_failed"],
        })
    }

    return ctx.Response().Json(http.StatusOK, map[string]interface{}{
        "message": "Mailbox connected successfully!",
        "account": profile.EmailAddress,
        "team":   team,
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



func (c *GmailController) ListMessages(ctx http.Context) http.Response {
	email := ctx.Request().Query("email")
	pageToken := ctx.Request().Query("pageToken")
	labelFilter  := ctx.Request().Query("label") // inbox, starred

	srv, err := GetClientFromDB(email)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError,
			map[string]string{"error": models.GmailAccountErrorMessage["not_found"]})
	}

	// Request 20 threads per page
	call := srv.Users.Threads.List("me").MaxResults(20)

	if labelFilter == "starred" {
		call = call.LabelIds("STARRED")
	} else if labelFilter == "inbox" {
		call = call.LabelIds("INBOX")
	}

	if pageToken != "" {
		call = call.PageToken(pageToken)
	}
	res, err := call.Do()
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError,
			map[string]string{"error": models.GmailAccountErrorMessage["kist_failed"]})
	}

	var g errgroup.Group
	threadMap := make(map[string]*gmail.Thread)
	replyCountMap := make(map[string]int)
	mu := &sync.Mutex{}

	// fetch each thread concurrently
	for _, t := range res.Threads {
		t := t
		g.Go(func() error {
			thread, err := srv.Users.Threads.Get("me", t.Id).Do()
			if err != nil {
				facades.Log().Errorf("Failed to fetch thread %s: %v", t.Id, err)
				return nil
			}
			mu.Lock()
			threadMap[t.Id] = thread
			replyCountMap[t.Id] = len(thread.Messages) - 1
			mu.Unlock()
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError,
			map[string]string{"error": models.GmailAccountErrorMessage["fetch_failed"]})
	}

	// build response
	messages := []map[string]any{}
	for threadId, thread := range threadMap {
		firstMessage := thread.Messages[0]
		latestMessage := thread.Messages[len(thread.Messages)-1]

		// parse headers
		var from, subject, date string
		from = getHeader(latestMessage.Payload.Headers,"From")
		if strings.Contains(from, "<") {
			from = strings.TrimSpace(strings.Split(from, "<")[0])
		}
		subject = getHeader(firstMessage.Payload.Headers,"Subject")
		date = getHeader(latestMessage.Payload.Headers,"Date")


		// unread and starred check
		isUnread := false
		isStarred := false
		isPromotions := false // Track if the message belongs to CATEGORY_PROMOTIONS
		for _, m := range thread.Messages {
			fmt.Println("Label: ",m.LabelIds);
			for _, lbl := range m.LabelIds {
				if lbl == "UNREAD" {
					isUnread = true
				}
				if lbl == "STARRED" {
					isStarred = true
				}
				if lbl == "CATEGORY_PROMOTIONS" {
					isPromotions = true
				}
			}
		}

		// Skip the message if it not has "CATEGORY_PROMOTIONS" label
		if isPromotions {
			continue // Skip this thread if it's promotions
		}

		// response row
		messages = append(messages, map[string]any{
			"id":           latestMessage.Id,
			"threadId":     threadId,
			"from":         from,
			"subject":      subject,
			"snippet":      latestMessage.Snippet,
			"date":         date,
			"replyCount":   replyCountMap[threadId],
			"isUnread":     isUnread,
			"isStarred":    isStarred,
		})
	}

	return ctx.Response().Json(http.StatusOK, map[string]any{
		"messages":      messages,
		"nextPageToken": res.NextPageToken,
	})
}


func (c *GmailController) ReadMessage(ctx http.Context) http.Response {
	email := ctx.Request().Query("email") // /gmail/message?id=xxxx&email=hgledgetech@gmail.com
	messageID := ctx.Request().Route("id")

	if messageID == "" || email == "" {
		facades.Log().Warningf("Missing parameters: email=%s, messageID=%s", email, messageID)
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": models.GmailAccountErrorMessage["invalid_request"],
		})
	}

	// Get Gmail client
	srv, err := GetClientFromDB(email)
	if err != nil {
		facades.Log().Errorf("Failed to get Gmail client for %s: %v", email, err)
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": models.GmailAccountErrorMessage["not_found"],
		})
	}

	// Fetch the thread directly using messageID
	// Trick: we don't know the threadId yet, so we fetch the message in metadata (cheap + small payload)
	msg, err := srv.Users.Messages.Get("me", messageID).Format("metadata").Do()
	if err != nil {
		facades.Log().Errorf("Failed to fetch message metadata (%s): %v", messageID, err)
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": models.GmailAccountErrorMessage["read_failed"],
		})
	}

	// Now fetch the whole thread (only ONE full request)
	thread, err := srv.Users.Threads.Get("me", msg.ThreadId).Format("full").Do()
	if err != nil {
		facades.Log().Errorf("Failed to fetch thread (%s): %v", msg.ThreadId, err)
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": models.GmailAccountErrorMessage["thread_not_found"],
		})
	}

	// Prepare messages array
	messages := []map[string]any{}
	for _, m := range thread.Messages {
		from := getHeader(m.Payload.Headers, "From")
		subject := getHeader(m.Payload.Headers, "Subject")
		date := getHeader(m.Payload.Headers, "Date")
		body := extractMessageBody(m.Payload)

		// Default to the raw header first
		formattedDate := date

		parsedDate, err := time.Parse(time.RFC1123Z, date)
		if err != nil {
			// Gmail sometimes includes "(PDT)" or "(UTC)" etc, strip it
			parts := strings.Split(date, "(")
			cleanDate := strings.TrimSpace(parts[0])

			parsedDate, err = time.Parse(time.RFC1123Z, cleanDate)
			if err != nil {
				// try fallback without zone offset
				parsedDate, err = time.Parse(time.RFC1123, cleanDate)
			}
		}

		if err == nil {
			formattedDate = parsedDate.Local().Format("Jan 02, 2006, 3:04 PM")
		}


		messages = append(messages, map[string]any{
			"id":      m.Id,
			"from":    from,
			"subject": subject,
			"date":    formattedDate,
			"snippet": m.Snippet,
			"body":    body,
			"labels":  m.LabelIds, // <-- contains "UNREAD"
		})

		fmt.Println("m: ",m.Id);
		// If this message was unread → mark it as read
		if contains(m.LabelIds, "UNREAD") {
			_, err := srv.Users.Messages.Modify("me", m.Id, &gmail.ModifyMessageRequest{
				RemoveLabelIds: []string{"UNREAD"},
			}).Do()
			if err != nil {
				facades.Log().Errorf("Failed to mark message %s as read: %v", messageID, err)
			}
		}
	}

	result := map[string]any{
		"threadId": thread.Id,
		"messages": messages,
	}

	return ctx.Response().Json(http.StatusOK, result)
}

// helper
func contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
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





// GET /gmail/accounts
func (c *GmailController) ListAccounts(ctx http.Context) http.Response {
	var accounts []models.GmailAccount
	if err := facades.Orm().Query().Get(&accounts); err != nil {
		facades.Log().Errorf("Failed to fetch Gmail accounts: %v", err)
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": models.GmailAccountErrorMessage["fetch_failed"],
		})
	}
	return ctx.Response().Json(http.StatusOK, accounts)
}

// DELETE /gmail/accounts/:email
func (c *GmailController) DeleteAccount(ctx http.Context) http.Response {
	email := ctx.Request().Route("email")

	if email == "" {
		facades.Log().Warning("Missing email parameter in DeleteAccount request")
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": models.GmailAccountErrorMessage["invalid_request"],
		})
	}

	if _,err := facades.Orm().Query().Where("email", email).Delete(&models.GmailAccount{}); err != nil {
		facades.Log().Errorf("Failed to delete Gmail account (%s): %v", email, err)
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": models.GmailAccountErrorMessage["delete_failed"],
		})
	}
	return ctx.Response().Json(http.StatusOK, map[string]string{
		"message": "Account removed",
	})
}


func (c *GmailController) GetGmailAccountTeams(ctx http.Context) http.Response{
	return ctx.Response().Json(http.StatusOK, models.GmailAccountTeams); 
}



func (c *GmailController) ToggleStar(ctx http.Context) http.Response {
	threadID := ctx.Request().Route("id")
	email := ctx.Request().Query("email")

	if threadID == "" || email == "" {
		facades.Log().Warningf("Missing parameters in ToggleStar: threadID=%s, email=%s", threadID, email)
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": models.GmailAccountErrorMessage["invalid_request"],
		})
	}

	srv, err := GetClientFromDB(email)
	if err != nil {
		facades.Log().Errorf("Failed to get Gmail client for %s: %v", email, err)
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": models.GmailAccountErrorMessage["not_found"],
		})
	}

	// 1. Get the thread with messages
	thread, err := srv.Users.Threads.Get("me", threadID).Format("full").Do()
	if err != nil {
		facades.Log().Errorf("Failed to fetch Gmail thread %s: %v", threadID, err)
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": models.GmailAccountErrorMessage["thread_not_found"],
		})
	}

	// fmt.Printf("Thread ID: %s\n", threadID)
	// fmt.Printf("Number of messages: %d\n", len(thread.Messages))
	// for i, msg := range thread.Messages {
	// 	fmt.Printf("Message %d ID: %s, Labels: %v\n", i, msg.Id, msg.LabelIds)
	// }

	if len(thread.Messages) == 0 {
		facades.Log().Warningf("Thread %s has no messages", threadID)
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": models.GmailAccountErrorMessage["thread_not_found"],
		})
	}

	// 2. Latest message = last in slice
	latest := thread.Messages[len(thread.Messages)-1]

	// 3. Check if latest is starred
	isStarred := false
	for _, lbl := range latest.LabelIds {
		if lbl == "STARRED" {
			isStarred = true
			break
		}
	}

	if isStarred {
		// 4a. Unstar = remove STARRED from all messages in thread
		for _, msg := range thread.Messages {
			_, err := srv.Users.Messages.Modify("me", msg.Id, &gmail.ModifyMessageRequest{
				RemoveLabelIds: []string{"STARRED"},
			}).Do()
			if err != nil {
				facades.Log().Errorf("Failed to unstar message %s in thread %s: %v", msg.Id, threadID, err)
				return ctx.Response().Json(http.StatusInternalServerError,
					map[string]string{"error": models.GmailAccountErrorMessage["update_failed"]})
			}
		}
	} else {
		// 4b. Star = add STARRED only to latest message
		_, err := srv.Users.Messages.Modify("me", latest.Id, &gmail.ModifyMessageRequest{
			AddLabelIds: []string{"STARRED"},
		}).Do()
		if err != nil {
			facades.Log().Errorf("Failed to star message %s in thread %s: %v", latest.Id, threadID, err)
			return ctx.Response().Json(http.StatusInternalServerError,
				map[string]string{"error": models.GmailAccountErrorMessage["update_failed"]})
		}
	}

	// 5. Return response
	return ctx.Response().Json(http.StatusOK, map[string]any{
		"threadId":  threadID,
		"messageId": latest.Id,
		"starred":   !isStarred, // new status
	})
}
