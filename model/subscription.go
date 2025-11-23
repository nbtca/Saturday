package model

import (
	"database/sql"
	"encoding/json"
)

// EventSubscription represents a subscription to event notifications
type EventSubscription struct {
	SubscriptionId int64           `json:"subscriptionId" db:"subscription_id"`
	MemberId       sql.NullString  `json:"memberId,omitempty" db:"member_id"`
	ClientId       sql.NullInt64   `json:"clientId,omitempty" db:"client_id"`
	EventTypes     []string        `json:"eventTypes" db:"event_types"`
	CallbackURL    string          `json:"callbackUrl" db:"callback_url"`
	Secret         string          `json:"secret" db:"secret"`
	Filters        json.RawMessage `json:"filters,omitempty" db:"filters"`
	Active         bool            `json:"active" db:"active"`
	GmtCreate      string          `json:"gmtCreate" db:"gmt_create"`
	GmtModified    string          `json:"gmtModified" db:"gmt_modified"`
}

// SubscriptionFilters represents the filtering options for subscriptions
type SubscriptionFilters struct {
	Status   []string `json:"status,omitempty"`
	MemberId string   `json:"memberId,omitempty"`
	ClientId int64    `json:"clientId,omitempty"`
}

// EventSubscriptionDelivery represents a webhook delivery attempt
type EventSubscriptionDelivery struct {
	DeliveryId     int64         `json:"deliveryId" db:"delivery_id"`
	SubscriptionId int64         `json:"subscriptionId" db:"subscription_id"`
	EventId        sql.NullInt64 `json:"eventId,omitempty" db:"event_id"`
	EventType      string        `json:"eventType" db:"event_type"`
	Status         string        `json:"status" db:"status"`
	Attempts       int           `json:"attempts" db:"attempts"`
	LastAttempt    sql.NullTime  `json:"lastAttempt,omitempty" db:"last_attempt"`
	ResponseCode   sql.NullInt32 `json:"responseCode,omitempty" db:"response_code"`
	ResponseBody   sql.NullString `json:"responseBody,omitempty" db:"response_body"`
	ErrorMessage   sql.NullString `json:"errorMessage,omitempty" db:"error_message"`
	GmtCreate      string        `json:"gmtCreate" db:"gmt_create"`
}

// WebhookPayload represents the data sent to webhook endpoints
type WebhookPayload struct {
	EventType string          `json:"event_type"`
	Timestamp string          `json:"timestamp"`
	Data      WebhookEventData `json:"data"`
}

// WebhookEventData contains the event details in the webhook payload
type WebhookEventData struct {
	EventId int64        `json:"event_id"`
	Status  string       `json:"status"`
	Action  string       `json:"action"`
	Actor   *WebhookActor `json:"actor,omitempty"`
	Event   PublicEvent  `json:"event"`
}

// WebhookActor represents who performed the action
type WebhookActor struct {
	MemberId string `json:"member_id,omitempty"`
	ClientId int64  `json:"client_id,omitempty"`
	Alias    string `json:"alias,omitempty"`
}

// PublicEventSubscription is a client-safe version without the secret
type PublicEventSubscription struct {
	SubscriptionId int64           `json:"subscriptionId" db:"subscription_id"`
	MemberId       sql.NullString  `json:"memberId,omitempty" db:"member_id"`
	ClientId       sql.NullInt64   `json:"clientId,omitempty" db:"client_id"`
	EventTypes     []string        `json:"eventTypes" db:"event_types"`
	CallbackURL    string          `json:"callbackUrl" db:"callback_url"`
	Filters        json.RawMessage `json:"filters,omitempty" db:"filters"`
	Active         bool            `json:"active" db:"active"`
	GmtCreate      string          `json:"gmtCreate" db:"gmt_create"`
	GmtModified    string          `json:"gmtModified" db:"gmt_modified"`
}

// ToPublic converts an EventSubscription to PublicEventSubscription (hides secret)
func (s EventSubscription) ToPublic() PublicEventSubscription {
	return PublicEventSubscription{
		SubscriptionId: s.SubscriptionId,
		MemberId:       s.MemberId,
		ClientId:       s.ClientId,
		EventTypes:     s.EventTypes,
		CallbackURL:    s.CallbackURL,
		Filters:        s.Filters,
		Active:         s.Active,
		GmtCreate:      s.GmtCreate,
		GmtModified:    s.GmtModified,
	}
}
