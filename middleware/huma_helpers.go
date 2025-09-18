package middleware

import (
	"context"
	"slices"
	"strconv"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/service"
	"github.com/nbtca/saturday/util"
)

// AuthContext represents the authentication context for Huma operations
type AuthContext struct {
	User        AuthContextUser
	ID          string
	Member      interface{}
	Role        string
	IsLegacyJWT bool
}

// GetAuthContext extracts auth context from context
func GetAuthContext(ctx context.Context) *AuthContext {
	return GetAuthContextFromContext(ctx)
}

// AuthenticateUser handles authentication for Huma operations
func AuthenticateUser(authHeader string, acceptableRoles ...Role) (*AuthContext, error) {
	if authHeader == "" {
		return nil, huma.Error401Unauthorized("not authorized, missing token")
	}

	// Handle legacy JWT token (used by wechat mini app)
	if len(strings.Split(authHeader, ".")) > 1 {
		tokenParsed, claims, err := util.ParseToken(authHeader)
		if err != nil || !tokenParsed.Valid {
			return nil, huma.Error401Unauthorized("not authorized, token not valid")
		}

		for _, roleObj := range acceptableRoles {
			if string(roleObj) == claims.Role {
				return &AuthContext{
					ID:          claims.Who,
					Member:      claims.Member,
					Role:        claims.Role,
					IsLegacyJWT: true,
				}, nil
			}
		}
		return nil, huma.Error401Unauthorized("not authorized")
	}

	// Strip bearer prefix
	tokenStr, err := util.GetTokenString(authHeader)
	if err != nil {
		return nil, huma.Error401Unauthorized("invalid token type")
	}

	// Fetch user info from Logto
	userinfo, err := service.LogtoServiceApp.FetchUserInfo(tokenStr)
	if err != nil {
		return nil, huma.Error401Unauthorized("not authorized: " + err.Error())
	}

	// Determine user roles
	var role string
	userRoles := []string{"client"}
	if slices.Contains(userinfo.Roles, "Repair Admin") {
		userRoles = append(userRoles, "admin")
		role = "admin"
	}
	if slices.Contains(userinfo.Roles, "Repair Member") {
		userRoles = append(userRoles, "member")
		if role == "" {
			role = "member"
		}
	}

	// Check if user has acceptable role
	for _, r := range acceptableRoles {
		if slices.Contains(userRoles, string(r)) {
			member, err := service.MemberServiceApp.GetMemberByLogtoId(userinfo.Sub)
			if err != nil {
				return nil, huma.Error401Unauthorized("not authorized")
			}

			user := AuthContextUser{
				Role:     userRoles,
				UserInfo: userinfo,
			}

			return &AuthContext{
				User:   user,
				ID:     member.MemberId,
				Member: member,
				Role:   role,
			}, nil
		}
	}
	return nil, huma.Error401Unauthorized("not authorized")
}

// LoadEvent loads an event by ID
func LoadEvent(eventId int64) (model.Event, error) {
	event, err := service.EventServiceApp.GetEventById(eventId)
	if err != nil {
		return model.Event{}, huma.Error422UnprocessableEntity(err.Error())
	}
	return event, nil
}

// GetClientIdFromAuth extracts client ID from authentication context
func GetClientIdFromAuth(auth *AuthContext) (int64, error) {
	if auth.IsLegacyJWT {
		// For legacy JWT tokens, parse the ID
		return strconv.ParseInt(auth.ID, 10, 64)
	}

	// For Logto tokens, create client by Logto ID if not exists
	if auth.User.UserInfo.Sub != "" {
		client, err := service.ClientServiceApp.CreateClientByLogtoIdIfNotExists(auth.User.UserInfo.Sub)
		if err != nil {
			return 0, err
		}
		return client.ClientId, nil
	}

	// Fallback to parsing ID
	return strconv.ParseInt(auth.ID, 10, 64)
}

// CreateIdentityFromAuth creates a model.Identity from AuthContext
func CreateIdentityFromAuth(auth *AuthContext) model.Identity {
	identity := model.Identity{
		Id:   auth.ID,
		Role: auth.Role,
	}

	if member, ok := auth.Member.(model.Member); ok {
		identity.Member = member
	}

	return identity
}

// MustGetAuthContext extracts auth context from context, panics if not found
func MustGetAuthContext(ctx context.Context) *AuthContext {
	authCtx := GetAuthContextFromContext(ctx)
	if authCtx == nil {
		panic("MustGetAuthContext called but no auth context found")
	}
	return authCtx
}

// GetClientIdFromAuthContext extracts client ID from auth context in context
func GetClientIdFromAuthContext(ctx context.Context) (int64, error) {
	authCtx := GetAuthContextFromContext(ctx)
	if authCtx == nil {
		return 0, huma.Error401Unauthorized("no authentication context found")
	}
	return GetClientIdFromAuth(authCtx)
}

// MustGetEventContext placeholder - not used in new approach
func MustGetEventContext(ctx context.Context) struct{ Event model.Event } {
	// This function is not used in the new approach
	panic("MustGetEventContext should not be called in new Huma operations")
}
