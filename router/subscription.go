package router

import (
	"context"
	"encoding/json"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nbtca/saturday/middleware"
	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/model/dto"
	"github.com/nbtca/saturday/service"
	"github.com/nbtca/saturday/util"
)

type SubscriptionRouter struct{}

// Input types for subscription endpoints

type SubscriptionPathInput struct {
	SubscriptionId int64 `path:"SubscriptionId" example:"123" doc:"Subscription ID"`
}

type CreateSubscriptionInput struct {
	AuthenticatedInput
	Body struct {
		EventTypes  []string        `json:"eventTypes" minItems:"1" doc:"Event types to subscribe to (e.g., event.created, event.accepted)"`
		CallbackURL string          `json:"callbackUrl" format:"uri" doc:"Webhook callback URL"`
		Filters     json.RawMessage `json:"filters,omitempty" doc:"Optional filters for events"`
	}
}

type GetSubscriptionsInput struct {
	AuthenticatedInput
	dto.PageRequest
}

type GetSubscriptionInput struct {
	AuthenticatedInput
	SubscriptionPathInput
}

type UpdateSubscriptionInput struct {
	AuthenticatedInput
	SubscriptionPathInput
	Body struct {
		EventTypes  []string        `json:"eventTypes" minItems:"1" doc:"Event types to subscribe to"`
		CallbackURL string          `json:"callbackUrl" format:"uri" doc:"Webhook callback URL"`
		Filters     json.RawMessage `json:"filters,omitempty" doc:"Optional filters for events"`
		Active      bool            `json:"active" doc:"Whether the subscription is active"`
	}
}

type DeleteSubscriptionInput struct {
	AuthenticatedInput
	SubscriptionPathInput
}

type GetDeliveryHistoryInput struct {
	AuthenticatedInput
	SubscriptionPathInput
	dto.PageRequest
}

// CreateSubscription creates a new event subscription
func (sr SubscriptionRouter) CreateSubscription(ctx context.Context, input *CreateSubscriptionInput) (*util.CommonResponse[model.PublicEventSubscription], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "member", "admin", "client")
	if err != nil {
		return nil, err
	}

	var memberId *string
	var clientId *int64

	// Determine if this is a member or client subscription
	if auth.Role == "member" || auth.Role == "admin" {
		memberId = &auth.ID
	} else if auth.Role == "client" {
		cid, err := middleware.GetClientIdFromAuth(auth)
		if err != nil {
			return nil, err
		}
		clientId = &cid
	}

	subscription, err := service.SubscriptionServiceApp.CreateSubscription(
		memberId,
		clientId,
		input.Body.EventTypes,
		input.Body.CallbackURL,
		input.Body.Filters,
	)

	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	return util.MakeCommonResponse(subscription), nil
}

// GetSubscriptions retrieves all subscriptions for the authenticated user
func (sr SubscriptionRouter) GetSubscriptions(ctx context.Context, input *GetSubscriptionsInput) (*util.CommonResponse[[]model.PublicEventSubscription], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "member", "admin", "client")
	if err != nil {
		return nil, err
	}

	var memberId *string
	var clientId *int64

	if auth.Role == "member" || auth.Role == "admin" {
		memberId = &auth.ID
	} else if auth.Role == "client" {
		cid, err := middleware.GetClientIdFromAuth(auth)
		if err != nil {
			return nil, err
		}
		clientId = &cid
	}

	subscriptions, total, err := service.SubscriptionServiceApp.GetSubscriptions(
		memberId,
		clientId,
		input.Limit,
		input.Offset,
	)

	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	return util.MakePaginatedResponse(subscriptions, total, input.Offset, input.Limit), nil
}

// GetSubscription retrieves a single subscription by ID
func (sr SubscriptionRouter) GetSubscription(ctx context.Context, input *GetSubscriptionInput) (*util.CommonResponse[model.PublicEventSubscription], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "member", "admin", "client")
	if err != nil {
		return nil, err
	}

	var memberId *string
	var clientId *int64

	if auth.Role == "member" || auth.Role == "admin" {
		memberId = &auth.ID
	} else if auth.Role == "client" {
		cid, err := middleware.GetClientIdFromAuth(auth)
		if err != nil {
			return nil, err
		}
		clientId = &cid
	}

	subscription, err := service.SubscriptionServiceApp.GetSubscription(
		input.SubscriptionId,
		memberId,
		clientId,
	)

	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	return util.MakeCommonResponse(subscription), nil
}

// UpdateSubscription updates an existing subscription
func (sr SubscriptionRouter) UpdateSubscription(ctx context.Context, input *UpdateSubscriptionInput) (*util.CommonResponse[model.PublicEventSubscription], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "member", "admin", "client")
	if err != nil {
		return nil, err
	}

	var memberId *string
	var clientId *int64

	if auth.Role == "member" || auth.Role == "admin" {
		memberId = &auth.ID
	} else if auth.Role == "client" {
		cid, err := middleware.GetClientIdFromAuth(auth)
		if err != nil {
			return nil, err
		}
		clientId = &cid
	}

	subscription, err := service.SubscriptionServiceApp.UpdateSubscription(
		input.SubscriptionId,
		memberId,
		clientId,
		input.Body.EventTypes,
		input.Body.CallbackURL,
		input.Body.Filters,
		input.Body.Active,
	)

	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	return util.MakeCommonResponse(subscription), nil
}

// DeleteSubscription deletes a subscription
func (sr SubscriptionRouter) DeleteSubscription(ctx context.Context, input *DeleteSubscriptionInput) (*util.CommonResponse[struct{ Success bool }], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "member", "admin", "client")
	if err != nil {
		return nil, err
	}

	var memberId *string
	var clientId *int64

	if auth.Role == "member" || auth.Role == "admin" {
		memberId = &auth.ID
	} else if auth.Role == "client" {
		cid, err := middleware.GetClientIdFromAuth(auth)
		if err != nil {
			return nil, err
		}
		clientId = &cid
	}

	err = service.SubscriptionServiceApp.DeleteSubscription(
		input.SubscriptionId,
		memberId,
		clientId,
	)

	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	return util.MakeCommonResponse(struct{ Success bool }{Success: true}), nil
}

// GetDeliveryHistory retrieves webhook delivery history for a subscription
func (sr SubscriptionRouter) GetDeliveryHistory(ctx context.Context, input *GetDeliveryHistoryInput) (*util.CommonResponse[[]model.EventSubscriptionDelivery], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "member", "admin", "client")
	if err != nil {
		return nil, err
	}

	var memberId *string
	var clientId *int64

	if auth.Role == "member" || auth.Role == "admin" {
		memberId = &auth.ID
	} else if auth.Role == "client" {
		cid, err := middleware.GetClientIdFromAuth(auth)
		if err != nil {
			return nil, err
		}
		clientId = &cid
	}

	deliveries, err := service.SubscriptionServiceApp.GetDeliveryHistory(
		input.SubscriptionId,
		memberId,
		clientId,
		input.Limit,
		input.Offset,
	)

	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	return util.MakeCommonResponse(deliveries), nil
}

var SubscriptionRouterApp = SubscriptionRouter{}
