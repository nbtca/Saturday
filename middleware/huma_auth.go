package middleware

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nbtca/saturday/service"
	"github.com/nbtca/saturday/util"
)

// AuthContextKey is the key for storing auth context in the request context
type AuthContextKey struct{}

// GetAuthContextKey returns the auth context key for testing purposes
func GetAuthContextKey() AuthContextKey {
	return AuthContextKey{}
}

// HumaAuthMiddleware provides authentication middleware for Huma v2
func HumaAuthMiddleware(ctx huma.Context, next func(huma.Context)) {
	// Get Authorization header
	authHeader := ctx.Header("Authorization")
	if authHeader == "" || authHeader == "Bearer" {
		// No auth header provided, continue without authentication
		// Individual operations will handle authorization requirements
		next(ctx)
		return
	}

	// Authenticate user
	authCtx, err := authenticateToken(authHeader)
	if err != nil {
		// Invalid token - let operation handlers decide if auth is required
		util.Logger.Debugf("Authentication failed in middleware: %v", err)
		next(ctx)
		return
	}

	// Attach auth context to the request context
	ctx = huma.WithValue(ctx, AuthContextKey{}, authCtx)
	next(ctx)
}

// GetAuthContextFromHuma extracts auth context from Huma context
func GetAuthContextFromHuma(ctx huma.Context) *AuthContext {
	if authCtx, ok := ctx.Context().Value(AuthContextKey{}).(*AuthContext); ok {
		return authCtx
	}
	return nil
}

// GetAuthContextFromContext extracts auth context from standard context
func GetAuthContextFromContext(ctx context.Context) *AuthContext {
	if authCtx, ok := ctx.Value(AuthContextKey{}).(*AuthContext); ok {
		return authCtx
	}
	return nil
}

// RequireAuth middleware for operations that require authentication
func RequireAuth(acceptableRoles ...Role) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		authCtx := GetAuthContextFromHuma(ctx)
		if authCtx == nil {
			// Set status and write error response
			ctx.SetStatus(http.StatusUnauthorized)
			return
		}

		// Check if user has acceptable role
		if len(acceptableRoles) > 0 {
			hasRole := false
			if authCtx.IsLegacyJWT {
				// For legacy JWT, check the single role
				for _, role := range acceptableRoles {
					if string(role) == authCtx.Role {
						hasRole = true
						break
					}
				}
			} else {
				// For Logto tokens, check multiple roles
				for _, role := range acceptableRoles {
					if slices.Contains(authCtx.User.Role, string(role)) {
						hasRole = true
						break
					}
				}
			}

			if !hasRole {
				ctx.SetStatus(http.StatusForbidden)
				return
			}
		}

		next(ctx)
	}
}

// authenticateToken handles token authentication (extracted from AuthenticateUser)
func authenticateToken(authHeader string) (*AuthContext, error) {
	// Handle legacy JWT token (used by wechat mini app)
	if len(strings.Split(authHeader, ".")) > 1 {
		tokenParsed, claims, err := util.ParseToken(authHeader)
		if err != nil || !tokenParsed.Valid {
			return nil, err
		}
		
		return &AuthContext{
			ID:          claims.Who,
			Member:      claims.Member,
			Role:        claims.Role,
			IsLegacyJWT: true,
		}, nil
	}

	// Strip bearer prefix
	tokenStr, err := util.GetTokenString(authHeader)
	if err != nil {
		return nil, err
	}

	// Fetch user info from Logto
	userinfo, err := service.LogtoServiceApp.FetchUserInfo(tokenStr)
	if err != nil {
		return nil, err
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

	member, err := service.MemberServiceApp.GetMemberByLogtoId(userinfo.Sub)
	if err != nil {
		return nil, err
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

// AuthenticateUserWithContext is an updated version that uses context-based auth when available
func AuthenticateUserWithContext(ctx context.Context, authHeader string, acceptableRoles ...Role) (*AuthContext, error) {
	// First try to get auth context from context (set by middleware)
	if authCtx := GetAuthContextFromContext(ctx); authCtx != nil {
		// Check if user has acceptable role
		if len(acceptableRoles) > 0 {
			hasRole := false
			if authCtx.IsLegacyJWT {
				// For legacy JWT, check the single role
				for _, role := range acceptableRoles {
					if string(role) == authCtx.Role {
						hasRole = true
						break
					}
				}
			} else {
				// For Logto tokens, check multiple roles
				for _, role := range acceptableRoles {
					if slices.Contains(authCtx.User.Role, string(role)) {
						hasRole = true
						break
					}
				}
			}

			if !hasRole {
				return nil, huma.Error401Unauthorized("not authorized")
			}
		}
		return authCtx, nil
	}

	// Fallback to original implementation for backward compatibility
	return AuthenticateUser(authHeader, acceptableRoles...)
}