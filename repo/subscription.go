package repo

import (
	"database/sql"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/util"
)

var subscriptionFields = []string{
	"subscription_id", "member_id", "client_id", "event_types",
	"delivery_method", "callback_url", "email", "secret", "scope",
	"filters", "active", "gmt_create", "gmt_modified",
}

var deliveryFields = []string{
	"delivery_id", "subscription_id", "event_id", "event_type",
	"status", "attempts", "last_attempt", "response_code",
	"response_body", "error_message", "gmt_create",
}

// SubscriptionFilter defines filter criteria for querying subscriptions
type SubscriptionFilter struct {
	Offset   uint64
	Limit    uint64
	MemberId *string
	ClientId *int64
	Active   *bool
	Order    string
}

// CreateSubscription creates a new event subscription
func CreateSubscription(subscription *model.EventSubscription) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	subscription.GmtCreate = now
	subscription.GmtModified = now

	sql, args, err := sq.Insert("event_subscription").
		Columns("member_id", "client_id", "event_types", "delivery_method", "callback_url", "email", "secret", "scope", "filters", "active", "gmt_create", "gmt_modified").
		Values(
			subscription.MemberId,
			subscription.ClientId,
			pq.Array(subscription.EventTypes),
			subscription.DeliveryMethod,
			subscription.CallbackURL,
			subscription.Email,
			subscription.Secret,
			subscription.Scope,
			subscription.Filters,
			subscription.Active,
			subscription.GmtCreate,
			subscription.GmtModified,
		).
		Suffix("RETURNING subscription_id").
		ToSql()

	if err != nil {
		return err
	}

	return db.QueryRow(sql, args...).Scan(&subscription.SubscriptionId)
}

// GetSubscriptionById retrieves a subscription by ID
func GetSubscriptionById(id int64) (model.EventSubscription, error) {
	sql, args, err := sq.Select(subscriptionFields...).
		From("event_subscription").
		Where(squirrel.Eq{"subscription_id": id}).
		ToSql()

	if err != nil {
		return model.EventSubscription{}, err
	}

	var subscription model.EventSubscription
	err = db.Get(&subscription, sql, args...)
	if err == sql.ErrNoRows {
		return model.EventSubscription{}, sql.ErrNoRows
	}

	return subscription, err
}

// GetSubscriptions retrieves subscriptions based on filter criteria
func GetSubscriptions(filter SubscriptionFilter, conditions ...squirrel.Eq) ([]model.EventSubscription, int64, error) {
	stat := sq.Select(subscriptionFields...).From("event_subscription")

	// Apply filter conditions
	if filter.MemberId != nil {
		stat = stat.Where(squirrel.Eq{"member_id": *filter.MemberId})
	}
	if filter.ClientId != nil {
		stat = stat.Where(squirrel.Eq{"client_id": *filter.ClientId})
	}
	if filter.Active != nil {
		stat = stat.Where(squirrel.Eq{"active": *filter.Active})
	}

	// Apply additional conditions
	for _, condition := range conditions {
		stat = stat.Where(condition)
	}

	// Get total count
	countSql, countArgs, err := stat.ToSql()
	if err != nil {
		return nil, 0, err
	}

	var total int64
	countQuery := "SELECT COUNT(*) FROM (" + countSql + ") AS count_query"
	if err := db.Get(&total, countQuery, countArgs...); err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	if filter.Order != "" {
		stat = stat.OrderBy(filter.Order)
	} else {
		stat = stat.OrderBy("gmt_create DESC")
	}

	if filter.Limit > 0 {
		stat = stat.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		stat = stat.Offset(filter.Offset)
	}

	sql, args, err := stat.ToSql()
	if err != nil {
		return nil, 0, err
	}

	var subscriptions []model.EventSubscription
	if err := db.Select(&subscriptions, sql, args...); err != nil {
		return nil, 0, err
	}

	return subscriptions, total, nil
}

// GetActiveSubscriptionsByEventType retrieves active subscriptions for a specific event type
func GetActiveSubscriptionsByEventType(eventType string) ([]model.EventSubscription, error) {
	// Use PostgreSQL array containment operator @>
	sql := `SELECT ` + strings.Join(subscriptionFields, ", ") + `
		FROM event_subscription
		WHERE active = true
		AND $1 = ANY(event_types)`

	var subscriptions []model.EventSubscription
	if err := db.Select(&subscriptions, sql, eventType); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// UpdateSubscription updates an existing subscription
func UpdateSubscription(subscription *model.EventSubscription) error {
	subscription.GmtModified = time.Now().Format("2006-01-02 15:04:05")

	sql, args, err := sq.Update("event_subscription").
		Set("event_types", pq.Array(subscription.EventTypes)).
		Set("delivery_method", subscription.DeliveryMethod).
		Set("callback_url", subscription.CallbackURL).
		Set("email", subscription.Email).
		Set("scope", subscription.Scope).
		Set("filters", subscription.Filters).
		Set("active", subscription.Active).
		Set("gmt_modified", subscription.GmtModified).
		Where(squirrel.Eq{"subscription_id": subscription.SubscriptionId}).
		ToSql()

	if err != nil {
		return err
	}

	_, err = db.Exec(sql, args...)
	return err
}

// DeleteSubscription soft-deletes a subscription by setting active to false
func DeleteSubscription(id int64) error {
	sql, args, err := sq.Update("event_subscription").
		Set("active", false).
		Set("gmt_modified", time.Now().Format("2006-01-02 15:04:05")).
		Where(squirrel.Eq{"subscription_id": id}).
		ToSql()

	if err != nil {
		return err
	}

	_, err = db.Exec(sql, args...)
	return err
}

// CreateDelivery creates a new delivery record
func CreateDelivery(delivery *model.EventSubscriptionDelivery) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	delivery.GmtCreate = now

	sql, args, err := sq.Insert("event_subscription_delivery").
		Columns("subscription_id", "event_id", "event_type", "status", "attempts", "last_attempt", "response_code", "response_body", "error_message", "gmt_create").
		Values(
			delivery.SubscriptionId,
			delivery.EventId,
			delivery.EventType,
			delivery.Status,
			delivery.Attempts,
			delivery.LastAttempt,
			delivery.ResponseCode,
			delivery.ResponseBody,
			delivery.ErrorMessage,
			delivery.GmtCreate,
		).
		Suffix("RETURNING delivery_id").
		ToSql()

	if err != nil {
		return err
	}

	return db.QueryRow(sql, args...).Scan(&delivery.DeliveryId)
}

// UpdateDelivery updates an existing delivery record
func UpdateDelivery(delivery *model.EventSubscriptionDelivery) error {
	sql, args, err := sq.Update("event_subscription_delivery").
		Set("status", delivery.Status).
		Set("attempts", delivery.Attempts).
		Set("last_attempt", delivery.LastAttempt).
		Set("response_code", delivery.ResponseCode).
		Set("response_body", delivery.ResponseBody).
		Set("error_message", delivery.ErrorMessage).
		Where(squirrel.Eq{"delivery_id": delivery.DeliveryId}).
		ToSql()

	if err != nil {
		return err
	}

	_, err = db.Exec(sql, args...)
	return err
}

// GetDeliveriesBySubscription retrieves delivery history for a subscription
func GetDeliveriesBySubscription(subscriptionId int64, limit, offset uint64) ([]model.EventSubscriptionDelivery, error) {
	stat := sq.Select(deliveryFields...).
		From("event_subscription_delivery").
		Where(squirrel.Eq{"subscription_id": subscriptionId}).
		OrderBy("gmt_create DESC")

	if limit > 0 {
		stat = stat.Limit(limit)
	}
	if offset > 0 {
		stat = stat.Offset(offset)
	}

	sql, args, err := stat.ToSql()
	if err != nil {
		return nil, err
	}

	var deliveries []model.EventSubscriptionDelivery
	if err := db.Select(&deliveries, sql, args...); err != nil {
		return nil, err
	}

	return deliveries, nil
}
