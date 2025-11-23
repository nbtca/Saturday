package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/repo"
	"github.com/nbtca/saturday/util"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

type SubscriptionService struct{}

var SubscriptionServiceApp = SubscriptionService{}

// CreateSubscription creates a new event subscription with a generated secret
func (s SubscriptionService) CreateSubscription(
	memberId *string,
	clientId *int64,
	eventTypes []string,
	deliveryMethod string,
	callbackURL *string,
	email *string,
	scope string,
	filters json.RawMessage,
) (model.PublicEventSubscription, error) {
	// Validate that at least one owner is specified
	if (memberId == nil || *memberId == "") && (clientId == nil || *clientId == 0) {
		return model.PublicEventSubscription{}, util.
			MakeServiceError(http.StatusBadRequest).
			SetMessage("Either memberId or clientId must be provided")
	}

	// Validate event types
	if len(eventTypes) == 0 {
		return model.PublicEventSubscription{}, util.
			MakeServiceError(http.StatusBadRequest).
			SetMessage("At least one event type must be specified")
	}

	// Validate delivery method
	if deliveryMethod != "webhook" && deliveryMethod != "email" && deliveryMethod != "both" {
		return model.PublicEventSubscription{}, util.
			MakeServiceError(http.StatusBadRequest).
			SetMessage("Delivery method must be 'webhook', 'email', or 'both'")
	}

	// Validate callback URL for webhook deliveries
	if (deliveryMethod == "webhook" || deliveryMethod == "both") && (callbackURL == nil || *callbackURL == "") {
		return model.PublicEventSubscription{}, util.
			MakeServiceError(http.StatusBadRequest).
			SetMessage("Callback URL is required for webhook delivery")
	}

	// Validate email for email deliveries
	if (deliveryMethod == "email" || deliveryMethod == "both") && (email == nil || *email == "") {
		return model.PublicEventSubscription{}, util.
			MakeServiceError(http.StatusBadRequest).
			SetMessage("Email is required for email delivery")
	}

	// Validate scope
	if scope != "related" && scope != "global" {
		return model.PublicEventSubscription{}, util.
			MakeServiceError(http.StatusBadRequest).
			SetMessage("Scope must be 'related' or 'global'")
	}

	// Generate secret for HMAC signature (for webhooks)
	secret, err := generateSecret(32)
	if err != nil {
		return model.PublicEventSubscription{}, util.
			MakeServiceError(http.StatusInternalServerError).
			SetMessage("Failed to generate secret")
	}

	subscription := &model.EventSubscription{
		EventTypes:     eventTypes,
		DeliveryMethod: deliveryMethod,
		Secret:         secret,
		Scope:          scope,
		Filters:        filters,
		Active:         true,
	}

	if callbackURL != nil && *callbackURL != "" {
		subscription.CallbackURL = sql.NullString{String: *callbackURL, Valid: true}
	}
	if email != nil && *email != "" {
		subscription.Email = sql.NullString{String: *email, Valid: true}
	}
	if memberId != nil && *memberId != "" {
		subscription.MemberId = sql.NullString{String: *memberId, Valid: true}
	}
	if clientId != nil && *clientId != 0 {
		subscription.ClientId = sql.NullInt64{Int64: *clientId, Valid: true}
	}

	if err := repo.CreateSubscription(subscription); err != nil {
		return model.PublicEventSubscription{}, util.
			MakeServiceError(http.StatusInternalServerError).
			SetMessage(err.Error())
	}

	return subscription.ToPublic(), nil
}

// GetSubscription retrieves a subscription by ID
func (s SubscriptionService) GetSubscription(id int64, memberId *string, clientId *int64) (model.PublicEventSubscription, error) {
	subscription, err := repo.GetSubscriptionById(id)
	if err == sql.ErrNoRows {
		return model.PublicEventSubscription{}, util.
			MakeServiceError(http.StatusNotFound).
			SetMessage("Subscription not found")
	}
	if err != nil {
		return model.PublicEventSubscription{}, util.
			MakeServiceError(http.StatusInternalServerError).
			SetMessage(err.Error())
	}

	// Verify ownership
	if !verifyOwnership(subscription, memberId, clientId) {
		return model.PublicEventSubscription{}, util.
			MakeServiceError(http.StatusForbidden).
			SetMessage("You do not have permission to access this subscription")
	}

	return subscription.ToPublic(), nil
}

// GetSubscriptions retrieves subscriptions with filtering
func (s SubscriptionService) GetSubscriptions(
	memberId *string,
	clientId *int64,
	limit, offset uint64,
) ([]model.PublicEventSubscription, int64, error) {
	filter := repo.SubscriptionFilter{
		Limit:    limit,
		Offset:   offset,
		MemberId: memberId,
		ClientId: clientId,
	}

	subscriptions, total, err := repo.GetSubscriptions(filter)
	if err != nil {
		return nil, 0, util.
			MakeServiceError(http.StatusInternalServerError).
			SetMessage(err.Error())
	}

	publicSubscriptions := make([]model.PublicEventSubscription, len(subscriptions))
	for i, sub := range subscriptions {
		publicSubscriptions[i] = sub.ToPublic()
	}

	return publicSubscriptions, total, nil
}

// UpdateSubscription updates an existing subscription
func (s SubscriptionService) UpdateSubscription(
	id int64,
	memberId *string,
	clientId *int64,
	eventTypes []string,
	deliveryMethod string,
	callbackURL *string,
	email *string,
	scope string,
	filters json.RawMessage,
	active bool,
) (model.PublicEventSubscription, error) {
	subscription, err := repo.GetSubscriptionById(id)
	if err == sql.ErrNoRows {
		return model.PublicEventSubscription{}, util.
			MakeServiceError(http.StatusNotFound).
			SetMessage("Subscription not found")
	}
	if err != nil {
		return model.PublicEventSubscription{}, util.
			MakeServiceError(http.StatusInternalServerError).
			SetMessage(err.Error())
	}

	// Verify ownership
	if !verifyOwnership(subscription, memberId, clientId) {
		return model.PublicEventSubscription{}, util.
			MakeServiceError(http.StatusForbidden).
			SetMessage("You do not have permission to update this subscription")
	}

	// Validate delivery method
	if deliveryMethod != "webhook" && deliveryMethod != "email" && deliveryMethod != "both" {
		return model.PublicEventSubscription{}, util.
			MakeServiceError(http.StatusBadRequest).
			SetMessage("Delivery method must be 'webhook', 'email', or 'both'")
	}

	// Validate callback URL for webhook deliveries
	if (deliveryMethod == "webhook" || deliveryMethod == "both") && (callbackURL == nil || *callbackURL == "") {
		return model.PublicEventSubscription{}, util.
			MakeServiceError(http.StatusBadRequest).
			SetMessage("Callback URL is required for webhook delivery")
	}

	// Validate email for email deliveries
	if (deliveryMethod == "email" || deliveryMethod == "both") && (email == nil || *email == "") {
		return model.PublicEventSubscription{}, util.
			MakeServiceError(http.StatusBadRequest).
			SetMessage("Email is required for email delivery")
	}

	// Validate scope
	if scope != "related" && scope != "global" {
		return model.PublicEventSubscription{}, util.
			MakeServiceError(http.StatusBadRequest).
			SetMessage("Scope must be 'related' or 'global'")
	}

	// Update fields
	subscription.EventTypes = eventTypes
	subscription.DeliveryMethod = deliveryMethod
	subscription.Scope = scope
	subscription.Filters = filters
	subscription.Active = active

	if callbackURL != nil && *callbackURL != "" {
		subscription.CallbackURL = sql.NullString{String: *callbackURL, Valid: true}
	} else {
		subscription.CallbackURL = sql.NullString{Valid: false}
	}
	if email != nil && *email != "" {
		subscription.Email = sql.NullString{String: *email, Valid: true}
	} else {
		subscription.Email = sql.NullString{Valid: false}
	}

	if err := repo.UpdateSubscription(&subscription); err != nil {
		return model.PublicEventSubscription{}, util.
			MakeServiceError(http.StatusInternalServerError).
			SetMessage(err.Error())
	}

	return subscription.ToPublic(), nil
}

// DeleteSubscription soft-deletes a subscription
func (s SubscriptionService) DeleteSubscription(id int64, memberId *string, clientId *int64) error {
	subscription, err := repo.GetSubscriptionById(id)
	if err == sql.ErrNoRows {
		return util.
			MakeServiceError(http.StatusNotFound).
			SetMessage("Subscription not found")
	}
	if err != nil {
		return util.
			MakeServiceError(http.StatusInternalServerError).
			SetMessage(err.Error())
	}

	// Verify ownership
	if !verifyOwnership(subscription, memberId, clientId) {
		return util.
			MakeServiceError(http.StatusForbidden).
			SetMessage("You do not have permission to delete this subscription")
	}

	if err := repo.DeleteSubscription(id); err != nil {
		return util.
			MakeServiceError(http.StatusInternalServerError).
			SetMessage(err.Error())
	}

	return nil
}

// NotifySubscribers sends notifications to all active subscriptions for an event
func (s SubscriptionService) NotifySubscribers(event model.Event, eventLog model.EventLog, actor *model.Identity) error {
	eventType := fmt.Sprintf("event.%s", eventLog.Action)

	// Get active subscriptions for this event type
	subscriptions, err := repo.GetActiveSubscriptionsByEventType(eventType)
	if err != nil {
		util.Logger.Error("Failed to get subscriptions: ", err)
		return err
	}

	if len(subscriptions) == 0 {
		util.Logger.Debug("No active subscriptions for event type: ", eventType)
		return nil
	}

	// Filter subscriptions based on scope
	filteredSubscriptions := filterSubscriptionsByScope(subscriptions, event, actor)

	if len(filteredSubscriptions) == 0 {
		util.Logger.Debug("No subscriptions match the scope for event: ", event.EventId)
		return nil
	}

	// Create webhook payload
	payload := createWebhookPayload(event, eventLog, actor, eventType)

	// Send notifications asynchronously based on delivery method
	for _, subscription := range filteredSubscriptions {
		sub := subscription // capture loop variable
		if sub.DeliveryMethod == "webhook" || sub.DeliveryMethod == "both" {
			go s.sendWebhook(sub, payload)
		}
		if sub.DeliveryMethod == "email" || sub.DeliveryMethod == "both" {
			go s.sendEmail(sub, event, eventLog, actor)
		}
	}

	return nil
}

// sendWebhook sends a webhook notification with retry logic
func (s SubscriptionService) sendWebhook(subscription model.EventSubscription, payload model.WebhookPayload) {
	const maxRetries = 3

	// Create delivery record
	delivery := &model.EventSubscriptionDelivery{
		SubscriptionId: subscription.SubscriptionId,
		EventType:      payload.EventType,
		Status:         "pending",
		Attempts:       0,
	}

	// Extract event ID if present
	if payload.Data.EventId > 0 {
		delivery.EventId = sql.NullInt64{Int64: payload.Data.EventId, Valid: true}
	}

	if err := repo.CreateDelivery(delivery); err != nil {
		util.Logger.Error("Failed to create delivery record: ", err)
		return
	}

	// Attempt delivery with retries
	for attempt := 1; attempt <= maxRetries; attempt++ {
		delivery.Attempts = attempt
		now := sql.NullTime{Time: time.Now(), Valid: true}
		delivery.LastAttempt = now

		// Marshal payload
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			delivery.Status = "failed"
			delivery.ErrorMessage = sql.NullString{String: err.Error(), Valid: true}
			repo.UpdateDelivery(delivery)
			return
		}

		// Create request
		req, err := http.NewRequest("POST", subscription.CallbackURL, bytes.NewBuffer(payloadBytes))
		if err != nil {
			delivery.Status = "failed"
			delivery.ErrorMessage = sql.NullString{String: err.Error(), Valid: true}
			repo.UpdateDelivery(delivery)
			return
		}

		// Add headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Saturday-Webhook/1.0")
		req.Header.Set("X-Event-Type", payload.EventType)

		// Generate HMAC signature
		signature := generateHMACSignature(payloadBytes, subscription.Secret)
		req.Header.Set("X-Signature", signature)

		// Send request
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			delivery.ErrorMessage = sql.NullString{String: err.Error(), Valid: true}

			// Retry with exponential backoff
			if attempt < maxRetries {
				backoff := time.Duration(attempt*attempt) * time.Second
				time.Sleep(backoff)
				repo.UpdateDelivery(delivery)
				continue
			}

			delivery.Status = "failed"
			repo.UpdateDelivery(delivery)
			return
		}

		defer resp.Body.Close()

		// Read response
		responseBody, _ := io.ReadAll(resp.Body)
		delivery.ResponseCode = sql.NullInt32{Int32: int32(resp.StatusCode), Valid: true}
		if len(responseBody) > 0 {
			delivery.ResponseBody = sql.NullString{String: string(responseBody), Valid: true}
		}

		// Check if successful
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			delivery.Status = "success"
			repo.UpdateDelivery(delivery)
			util.Logger.Info(fmt.Sprintf("Webhook delivered successfully to %s", subscription.CallbackURL))
			return
		}

		// Retry on 5xx errors
		if resp.StatusCode >= 500 && attempt < maxRetries {
			backoff := time.Duration(attempt*attempt) * time.Second
			time.Sleep(backoff)
			repo.UpdateDelivery(delivery)
			continue
		}

		// Non-retryable error
		delivery.Status = "failed"
		delivery.ErrorMessage = sql.NullString{
			String: fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(responseBody)),
			Valid:  true,
		}
		repo.UpdateDelivery(delivery)
		return
	}

	// All retries exhausted
	delivery.Status = "failed"
	delivery.ErrorMessage = sql.NullString{String: "Max retries exhausted", Valid: true}
	repo.UpdateDelivery(delivery)
}

// createWebhookPayload creates a webhook payload from event data
func createWebhookPayload(event model.Event, eventLog model.EventLog, actor *model.Identity, eventType string) model.WebhookPayload {
	payload := model.WebhookPayload{
		EventType: eventType,
		Timestamp: time.Now().Format(time.RFC3339),
		Data: model.WebhookEventData{
			EventId: event.EventId,
			Status:  event.Status,
			Action:  eventLog.Action,
			Event:   model.CreatePublicEvent(event),
		},
	}

	// Add actor information if available
	if actor != nil {
		webhookActor := &model.WebhookActor{}
		if actor.IsMember() {
			webhookActor.MemberId = actor.MemberId
			if event.Member != nil {
				webhookActor.Alias = event.Member.Alias
			}
		} else if actor.IsClient() {
			webhookActor.ClientId = actor.ClientId
		}
		payload.Data.Actor = webhookActor
	}

	return payload
}

// generateSecret generates a random secret for HMAC signatures
func generateSecret(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// generateHMACSignature generates an HMAC SHA-256 signature
func generateHMACSignature(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

// verifyOwnership checks if the user has permission to access the subscription
func verifyOwnership(subscription model.EventSubscription, memberId *string, clientId *int64) bool {
	// Check member ownership
	if memberId != nil && *memberId != "" {
		if subscription.MemberId.Valid && subscription.MemberId.String == *memberId {
			return true
		}
	}

	// Check client ownership
	if clientId != nil && *clientId != 0 {
		if subscription.ClientId.Valid && subscription.ClientId.Int64 == *clientId {
			return true
		}
	}

	return false
}

// GetDeliveryHistory retrieves webhook delivery history for a subscription
func (s SubscriptionService) GetDeliveryHistory(
	subscriptionId int64,
	memberId *string,
	clientId *int64,
	limit, offset uint64,
) ([]model.EventSubscriptionDelivery, error) {
	// Verify subscription ownership
	subscription, err := repo.GetSubscriptionById(subscriptionId)
	if err == sql.ErrNoRows {
		return nil, util.
			MakeServiceError(http.StatusNotFound).
			SetMessage("Subscription not found")
	}
	if err != nil {
		return nil, util.
			MakeServiceError(http.StatusInternalServerError).
			SetMessage(err.Error())
	}

	if !verifyOwnership(subscription, memberId, clientId) {
		return nil, util.
			MakeServiceError(http.StatusForbidden).
			SetMessage("You do not have permission to access this subscription")
	}

	deliveries, err := repo.GetDeliveriesBySubscription(subscriptionId, limit, offset)
	if err != nil {
		return nil, util.
			MakeServiceError(http.StatusInternalServerError).
			SetMessage(err.Error())
	}

	return deliveries, nil
}

// filterSubscriptionsByScope filters subscriptions based on their scope (related vs global)
func filterSubscriptionsByScope(subscriptions []model.EventSubscription, event model.Event, actor *model.Identity) []model.EventSubscription {
	filtered := []model.EventSubscription{}

	for _, sub := range subscriptions {
		// Global subscriptions always match
		if sub.Scope == "global" {
			filtered = append(filtered, sub)
			continue
		}

		// For "related" scope, check if the subscription owner is related to the event
		if sub.Scope == "related" {
			isRelated := false

			// Member subscriptions: check if member is assigned to event or performed the action
			if sub.MemberId.Valid {
				if event.MemberId == sub.MemberId.String {
					isRelated = true
				}
				if event.ClosedBy == sub.MemberId.String {
					isRelated = true
				}
				if actor != nil && actor.IsMember() && actor.MemberId == sub.MemberId.String {
					isRelated = true
				}
			}

			// Client subscriptions: check if client owns the event
			if sub.ClientId.Valid {
				if event.ClientId == sub.ClientId.Int64 {
					isRelated = true
				}
				if actor != nil && actor.IsClient() && actor.ClientId == sub.ClientId.Int64 {
					isRelated = true
				}
			}

			if isRelated {
				filtered = append(filtered, sub)
			}
		}
	}

	return filtered
}

// sendEmail sends an email notification for a subscription
func (s SubscriptionService) sendEmail(subscription model.EventSubscription, event model.Event, eventLog model.EventLog, actor *model.Identity) {
	if !subscription.Email.Valid || subscription.Email.String == "" {
		util.Logger.Error("Email not configured for subscription: ", subscription.SubscriptionId)
		return
	}

	// Determine the subject based on the action
	subject := fmt.Sprintf("[NBTCA] Event #%d: %s", event.EventId, eventLog.Action)

	// Build HTML email body
	htmlBody := fmt.Sprintf(`
<html>
<body>
<h3>Event Notification</h3>
<table style="border-collapse: collapse; width: 100%%;">
	<tr><td style="padding: 8px; font-weight: bold;">Event ID:</td><td style="padding: 8px;">%d</td></tr>
	<tr><td style="padding: 8px; font-weight: bold;">Action:</td><td style="padding: 8px;">%s</td></tr>
	<tr><td style="padding: 8px; font-weight: bold;">Status:</td><td style="padding: 8px;">%s</td></tr>
	<tr><td style="padding: 8px; font-weight: bold;">Problem:</td><td style="padding: 8px;">%s</td></tr>
	<tr><td style="padding: 8px; font-weight: bold;">Model:</td><td style="padding: 8px;">%s</td></tr>
	<tr><td style="padding: 8px; font-weight: bold;">Time:</td><td style="padding: 8px;">%s</td></tr>`,
		event.EventId, eventLog.Action, event.Status, event.Problem, event.Model, eventLog.GmtCreate)

	if actor != nil {
		if actor.IsMember() {
			htmlBody += fmt.Sprintf(`
	<tr><td style="padding: 8px; font-weight: bold;">Performed by:</td><td style="padding: 8px;">Member %s</td></tr>`, actor.MemberId)
		} else if actor.IsClient() {
			htmlBody += fmt.Sprintf(`
	<tr><td style="padding: 8px; font-weight: bold;">Performed by:</td><td style="padding: 8px;">Client %d</td></tr>`, actor.ClientId)
		}
	}

	htmlBody += `
</table>
<div style="margin-top: 20px;">
`

	if event.GithubIssueNumber.Valid && event.GithubIssueNumber.Int64 > 0 {
		htmlBody += fmt.Sprintf(`<p><a href="https://github.com/nbtca/Saturday/issues/%d">View GitHub Issue</a></p>`, event.GithubIssueNumber.Int64)
	}

	htmlBody += fmt.Sprintf(`<p><a href="https://repair.nbtca.space/events/%d">View Event Details</a></p>
</div>
</body>
</html>`, event.EventId)

	// Create email message
	m := gomail.NewMessage()
	m.SetHeader("To", subscription.Email.String)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	// Get SMTP configuration from viper
	smtpHost := viper.GetString("smtp.host")
	smtpPort := viper.GetInt("smtp.port")
	smtpUser := viper.GetString("smtp.user")
	smtpPassword := viper.GetString("smtp.password")

	if smtpHost == "" {
		util.Logger.Error("SMTP host not configured")
		return
	}

	// Create dialer and send
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPassword)
	if err := d.DialAndSend(m); err != nil {
		util.Logger.Error(fmt.Sprintf("Failed to send email to %s: %v", subscription.Email.String, err))
	} else {
		util.Logger.Info(fmt.Sprintf("Email sent successfully to %s for event %d", subscription.Email.String, event.EventId))
	}
}
