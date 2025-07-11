package service

import (
	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/model/dto"
)

// MemberServiceInterface defines the contract for member-related operations
type MemberServiceInterface interface {
	// Member retrieval
	GetMemberById(id string) (model.Member, error)
	GetMemberByLogtoId(logtoId string) (model.Member, error)
	GetMemberByGithubId(githubId string) (model.Member, error)
	GetPublicMemberById(id string) (model.PublicMember, error)
	GetPublicMembers(offset uint64, limit uint64) ([]model.PublicMember, error)
	GetMembers(offset uint64, limit uint64) ([]model.Member, error)
	
	// Member management
	CreateMember(member *model.Member) error
	UpdateMember(member model.Member) error
	ActivateMember(member model.Member) error
	
	// Authentication and role management
	CreateToken(member model.Member) (string, error)
	MapLogtoUserRole(roles []LogtoUserRole) string
	
	// New methods to handle business logic from router layer
	GetOrCreateMemberByLogtoId(logtoId string) (model.Member, error)
	SyncMemberProfile(member model.Member, logtoUser *FetchLogtoUsersResponse) error
	AuthenticateAndSyncMember(logtoToken string) (model.Member, string, error)
}

// LogtoServiceInterface defines the contract for Logto authentication operations
type LogtoServiceInterface interface {
	// Token operations
	FetchLogtoToken(resource string, scope string) (map[string]interface{}, error)
	
	// User operations
	FetchUsers(request FetchLogtoUsersRequest) ([]FetchLogtoUsersResponse, error)
	FetchUserById(userId string) (*FetchLogtoUsersResponse, error)
	FetchUserByToken(token string) (*FetchLogtoUsersResponse, error)
	FetchUserInfo(accessToken string) (FetchUserInfoResponse, error)
	PatchUserById(userId string, data dto.PatchLogtoUserRequest) (map[string]interface{}, error)
	
	// Role operations
	FetchUserRole(userId string) (FetchUserRoleResponse, error)
}

// EventServiceInterface defines the contract for event-related operations
type EventServiceInterface interface {
	// Event retrieval with proper authorization
	GetEventByIdWithAuth(eventId string, memberId int, clientId int) (interface{}, error)
	
	// Event filtering and export
	GetEventsWithFilter(filter interface{}, offset uint64, limit uint64) ([]interface{}, error)
	ExportEventsToExcel(filter interface{}) ([]byte, string, error)
	
	// Event management
	CreateEvent(event interface{}) error
	UpdateEvent(event interface{}) error
	DeleteEvent(eventId string, memberId int, clientId int) error
}

// ClientServiceInterface defines the contract for client-related operations
type ClientServiceInterface interface {
	GetClientById(id string) (interface{}, error)
	CreateClient(client interface{}) error
	UpdateClient(client interface{}) error
	DeleteClient(id string) error
}

// AuthServiceInterface defines the contract for authentication operations
type AuthServiceInterface interface {
	// Authentication
	ValidateLogtoToken(token string) (*FetchLogtoUsersResponse, error)
	ValidateLegacyJWT(token string) (interface{}, error)
	
	// Authorization
	CheckPermission(memberId int, resource string, action string) error
	CheckEventAccess(eventId string, memberId int, clientId int) error
	
	// Role management
	MapUserRoles(logtoRoles []LogtoUserRole) string
	SyncUserRole(memberId int, newRole string) error
}