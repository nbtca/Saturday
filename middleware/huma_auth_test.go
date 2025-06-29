package middleware_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nbtca/saturday/middleware"
	"github.com/nbtca/saturday/util"
)

// MockHumaContext implements huma.Context for testing
type MockHumaContext struct {
	headers       map[string][]string
	ctx           context.Context
	status        int
	method        string
	urlValue      url.URL
	params        map[string]string
	queries       map[string]string
	bodyReader    io.Reader
	bodyWriter    io.Writer
	host          string
	remoteAddr    string
	operation     *huma.Operation
	multipartForm *multipart.Form
}

func NewMockHumaContext() *MockHumaContext {
	return &MockHumaContext{
		headers:    make(map[string][]string),
		params:     make(map[string]string),
		queries:    make(map[string]string),
		ctx:        context.Background(),
		status:     0,
		method:     "GET",
		urlValue:   url.URL{Path: "/test"},
		host:       "localhost",
		remoteAddr: "127.0.0.1:12345",
		bodyWriter: &bytes.Buffer{},
	}
}

func (m *MockHumaContext) Operation() *huma.Operation {
	return m.operation
}

func (m *MockHumaContext) Context() context.Context {
	return m.ctx
}

func (m *MockHumaContext) TLS() *tls.ConnectionState {
	return nil
}

func (m *MockHumaContext) Version() huma.ProtoVersion {
	return huma.ProtoVersion{Proto: "HTTP/1.1", ProtoMajor: 1, ProtoMinor: 1}
}

func (m *MockHumaContext) Method() string {
	return m.method
}

func (m *MockHumaContext) Host() string {
	return m.host
}

func (m *MockHumaContext) RemoteAddr() string {
	return m.remoteAddr
}

func (m *MockHumaContext) URL() url.URL {
	return m.urlValue
}

func (m *MockHumaContext) Param(name string) string {
	return m.params[name]
}

func (m *MockHumaContext) Query(name string) string {
	return m.queries[name]
}

func (m *MockHumaContext) Header(name string) string {
	values := m.headers[name]
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

func (m *MockHumaContext) EachHeader(cb func(name, value string)) {
	for name, values := range m.headers {
		for _, value := range values {
			cb(name, value)
		}
	}
}

func (m *MockHumaContext) BodyReader() io.Reader {
	return m.bodyReader
}

func (m *MockHumaContext) GetMultipartForm() (*multipart.Form, error) {
	return m.multipartForm, nil
}

func (m *MockHumaContext) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *MockHumaContext) SetStatus(status int) {
	m.status = status
}

func (m *MockHumaContext) Status() int {
	return m.status
}

func (m *MockHumaContext) SetHeader(name, value string) {
	m.headers[name] = []string{value}
}

func (m *MockHumaContext) AppendHeader(name, value string) {
	m.headers[name] = append(m.headers[name], value)
}

func (m *MockHumaContext) BodyWriter() io.Writer {
	return m.bodyWriter
}

func (m *MockHumaContext) WithContext(ctx context.Context) huma.Context {
	m.ctx = ctx
	return m
}

// Helper methods for testing
func (m *MockHumaContext) SetParam(name, value string) {
	m.params[name] = value
}

func (m *MockHumaContext) SetQuery(name, value string) {
	m.queries[name] = value
}

func (m *MockHumaContext) SetBodyReader(r io.Reader) {
	m.bodyReader = r
}

func (m *MockHumaContext) SetMultipartForm(form *multipart.Form) {
	m.multipartForm = form
}

func (m *MockHumaContext) SetAuthHeader(value string) {
	m.SetHeader("Authorization", value)
}

// Test data structure for auth middleware tests
type HumaAuthTestCase struct {
	CaseID        string
	Description   string
	Authorization string
	ExpectedAuth  bool
	ExpectedRole  string
	ShouldFail    bool
}

func TestHumaAuthMiddleware(t *testing.T) {
	testCases := []HumaAuthTestCase{
		{
			CaseID:        "1.1",
			Description:   "Valid legacy JWT token with member role",
			Authorization: util.GenToken("member", "2333333333"),
			ExpectedAuth:  true,
			ExpectedRole:  "member",
			ShouldFail:    false,
		},
		{
			CaseID:        "1.2", 
			Description:   "Valid legacy JWT token with admin role",
			Authorization: util.GenToken("admin", "2333333333"),
			ExpectedAuth:  true,
			ExpectedRole:  "admin",
			ShouldFail:    false,
		},
		{
			CaseID:        "1.3",
			Description:   "Invalid JWT token",
			Authorization: "invalid.jwt.token",
			ExpectedAuth:  false,
			ExpectedRole:  "",
			ShouldFail:    false, // Middleware should continue, let handlers decide
		},
		{
			CaseID:        "1.4",
			Description:   "Missing authorization header",
			Authorization: "",
			ExpectedAuth:  false,
			ExpectedRole:  "",
			ShouldFail:    false, // Middleware should continue
		},
		{
			CaseID:        "1.5",
			Description:   "Bearer prefix only",
			Authorization: "Bearer",
			ExpectedAuth:  false,
			ExpectedRole:  "",
			ShouldFail:    false,
		},
		{
			CaseID:        "1.6",
			Description:   "Malformed bearer token",
			Authorization: "Bearer malformed_token",
			ExpectedAuth:  false,
			ExpectedRole:  "",
			ShouldFail:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.CaseID+"_"+tc.Description, func(t *testing.T) {
			// Create mock context
			mockCtx := NewMockHumaContext()
			if tc.Authorization != "" {
				mockCtx.SetAuthHeader(tc.Authorization)
			}

			// Track if next was called and capture the context
			nextCalled := false
			var capturedCtx huma.Context
			next := func(ctx huma.Context) {
				nextCalled = true
				capturedCtx = ctx
			}

			// Execute middleware
			middleware.HumaAuthMiddleware(mockCtx, next)

			// Verify next was called (middleware should always continue)
			if !nextCalled {
				t.Error("Next middleware should always be called")
			}

			// Check auth context from the captured context
			var authCtx *middleware.AuthContext
			if capturedCtx != nil {
				authCtx = middleware.GetAuthContextFromHuma(capturedCtx)
			}
			
			if tc.ExpectedAuth {
				if authCtx == nil {
					t.Errorf("Expected auth context to be set for case %s", tc.CaseID)
					return
				}
				
				if authCtx.Role != tc.ExpectedRole {
					t.Errorf("Expected role %s, got %s", tc.ExpectedRole, authCtx.Role)
				}
				
				if authCtx.ID == "" {
					t.Error("Expected user ID to be set")
				}
				
				if !authCtx.IsLegacyJWT {
					t.Error("Expected legacy JWT flag to be true")
				}
			} else {
				if authCtx != nil {
					t.Errorf("Expected no auth context for case %s, but got: %+v", tc.CaseID, authCtx)
				}
			}
		})
	}
}

func TestRequireAuthMiddleware(t *testing.T) {
	testCases := []struct {
		CaseID           string
		Description      string
		AuthContext      *middleware.AuthContext
		RequiredRoles    []middleware.Role
		ExpectedStatus   int
		ExpectedContinue bool
	}{
		{
			CaseID:      "2.1",
			Description: "Admin access with admin role required",
			AuthContext: &middleware.AuthContext{
				ID:          "admin123",
				Role:        "admin",
				IsLegacyJWT: true,
			},
			RequiredRoles:    []middleware.Role{"admin"},
			ExpectedStatus:   0,
			ExpectedContinue: true,
		},
		{
			CaseID:      "2.2", 
			Description: "Member access with member or admin role required",
			AuthContext: &middleware.AuthContext{
				ID:          "member123",
				Role:        "member",
				IsLegacyJWT: true,
			},
			RequiredRoles:    []middleware.Role{"member", "admin"},
			ExpectedStatus:   0,
			ExpectedContinue: true,
		},
		{
			CaseID:           "2.3",
			Description:      "No auth context with role required",
			AuthContext:      nil,
			RequiredRoles:    []middleware.Role{"admin"},
			ExpectedStatus:   http.StatusUnauthorized,
			ExpectedContinue: false,
		},
		{
			CaseID:      "2.4",
			Description: "Client access with admin role required",
			AuthContext: &middleware.AuthContext{
				ID:          "client123",
				Role:        "client", 
				IsLegacyJWT: true,
			},
			RequiredRoles:    []middleware.Role{"admin"},
			ExpectedStatus:   http.StatusForbidden,
			ExpectedContinue: false,
		},
		{
			CaseID:      "2.5",
			Description: "Any authenticated user (no role requirements)",
			AuthContext: &middleware.AuthContext{
				ID:          "user123",
				Role:        "client",
				IsLegacyJWT: true,
			},
			RequiredRoles:    []middleware.Role{},
			ExpectedStatus:   0,
			ExpectedContinue: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.CaseID+"_"+tc.Description, func(t *testing.T) {
			// Create mock context
			mockCtx := NewMockHumaContext()
			
			// Set auth context if provided
			if tc.AuthContext != nil {
				mockCtx.ctx = context.WithValue(mockCtx.ctx, middleware.GetAuthContextKey(), tc.AuthContext)
			}

			// Track if next was called
			nextCalled := false
			next := func(ctx huma.Context) {
				nextCalled = true
			}

			// Create and execute RequireAuth middleware
			authMiddleware := middleware.RequireAuth(tc.RequiredRoles...)
			authMiddleware(mockCtx, next)

			// Verify status
			if mockCtx.Status() != tc.ExpectedStatus {
				t.Errorf("Expected status %d, got %d", tc.ExpectedStatus, mockCtx.Status())
			}

			// Verify if next was called
			if nextCalled != tc.ExpectedContinue {
				t.Errorf("Expected next called: %v, got: %v", tc.ExpectedContinue, nextCalled)
			}
		})
	}
}

func TestGetAuthContextFromHuma(t *testing.T) {
	testCases := []struct {
		CaseID      string
		Description string
		SetContext  bool
		AuthContext *middleware.AuthContext
	}{
		{
			CaseID:      "3.1",
			Description: "Context with valid auth",
			SetContext:  true,
			AuthContext: &middleware.AuthContext{
				ID:          "test123",
				Role:        "admin",
				IsLegacyJWT: true,
			},
		},
		{
			CaseID:      "3.2",
			Description: "Context without auth",
			SetContext:  false,
			AuthContext: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.CaseID+"_"+tc.Description, func(t *testing.T) {
			mockCtx := NewMockHumaContext()
			
			if tc.SetContext && tc.AuthContext != nil {
				mockCtx.ctx = context.WithValue(mockCtx.ctx, middleware.GetAuthContextKey(), tc.AuthContext)
			}

			result := middleware.GetAuthContextFromHuma(mockCtx)

			if tc.SetContext {
				if result == nil {
					t.Error("Expected auth context to be returned")
					return
				}
				if result.ID != tc.AuthContext.ID {
					t.Errorf("Expected ID %s, got %s", tc.AuthContext.ID, result.ID)
				}
				if result.Role != tc.AuthContext.Role {
					t.Errorf("Expected Role %s, got %s", tc.AuthContext.Role, result.Role)
				}
			} else {
				if result != nil {
					t.Error("Expected no auth context to be returned")
				}
			}
		})
	}
}

func TestAuthenticateUserWithContext(t *testing.T) {
	testCases := []struct {
		CaseID        string
		Description   string
		ContextAuth   *middleware.AuthContext
		AuthHeader    string
		RequiredRoles []middleware.Role
		ShouldSucceed bool
		ExpectedRole  string
	}{
		{
			CaseID:      "4.1",
			Description: "Context auth with valid role",
			ContextAuth: &middleware.AuthContext{
				ID:          "ctx123",
				Role:        "admin",
				IsLegacyJWT: true,
			},
			AuthHeader:    util.GenToken("admin", "header123"),
			RequiredRoles: []middleware.Role{"admin"},
			ShouldSucceed: true,
			ExpectedRole:  "admin",
		},
		{
			CaseID:        "4.2",
			Description:   "No context auth, fallback to header",
			ContextAuth:   nil,
			AuthHeader:    util.GenToken("member", "header123"),
			RequiredRoles: []middleware.Role{"member"},
			ShouldSucceed: true,
			ExpectedRole:  "member",
		},
		{
			CaseID:      "4.3",
			Description: "Context auth with insufficient role",
			ContextAuth: &middleware.AuthContext{
				ID:          "ctx123",
				Role:        "client",
				IsLegacyJWT: true,
			},
			AuthHeader:    util.GenToken("client", "header123"),
			RequiredRoles: []middleware.Role{"admin"},
			ShouldSucceed: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.CaseID+"_"+tc.Description, func(t *testing.T) {
			ctx := context.Background()
			
			if tc.ContextAuth != nil {
				ctx = context.WithValue(ctx, middleware.GetAuthContextKey(), tc.ContextAuth)
			}

			result, err := middleware.AuthenticateUserWithContext(ctx, tc.AuthHeader, tc.RequiredRoles...)

			if tc.ShouldSucceed {
				if err != nil {
					t.Errorf("Expected success but got error: %v", err)
					return
				}
				if result == nil {
					t.Error("Expected result but got nil")
					return
				}
				if result.Role != tc.ExpectedRole {
					t.Errorf("Expected role %s, got %s", tc.ExpectedRole, result.Role)
				}
			} else {
				if err == nil {
					t.Error("Expected error but got success")
				}
			}
		})
	}
}