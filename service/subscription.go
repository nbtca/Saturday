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
)

type SubscriptionService struct{}

var SubscriptionServiceApp = SubscriptionService{}

// CreateSubscription creates a new event subscription with a generated secret
func (s SubscriptionService) CreateSubscription(
	memberId *string,
	clientId *int64,
	eventTypes []string,
	callbackURL string,
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

	// Validate callback URL
	if callbackURL == "" {
		return model.PublicEventSubscription{}, util.
			MakeServiceError(http.StatusBadRequest).
			SetMessage("Callback URL is required")
	}

	// Generate secret for HMAC signature
	secret, err := generateSecret(32)
	if err != nil {
		return model.PublicEventSubscription{}, util.
			MakeServiceError(http.StatusInternalServerError).
			SetMessage("Failed to generate secret")
	}

	subscription := &model.EventSubscription{
		EventTypes:  eventTypes,
		CallbackURL: callbackURL,
		Secret:      secret,
		Filters:     filters,
		Active:      true,
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
	callbackURL string,
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

	// Update fields
	subscription.EventTypes = eventTypes
	subscription.CallbackURL = callbackURL
	subscription.Filters = filters
	subscription.Active = active

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

// NotifySubscribers sends webhook notifications to all active subscriptions for an event
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

	// Create webhook payload
	payload := createWebhookPayload(event, eventLog, actor, eventType)

	// Send notifications asynchronously
	for _, subscription := range subscriptions {
		go s.sendWebhook(subscription, payload)
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
