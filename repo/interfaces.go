package repo

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/nbtca/saturday/model"
)

// MemberRepository defines the interface for member-related database operations
type MemberRepository interface {
	// Member retrieval
	GetMemberById(id string) (model.Member, error)
	GetMemberByLogtoId(logtoId string) (model.Member, error)
	GetMemberByGithubId(githubId string) (model.Member, error)
	GetMemberIdByLogtoId(logtoId string) (sql.NullString, error)
	GetMembers(offset uint64, limit uint64) ([]model.Member, error)
	
	// Member management
	CreateMember(member *model.Member) error
	UpdateMember(member model.Member) error
	ExistMember(id string) (bool, error)
}

// EventRepository defines the interface for event-related database operations
type EventRepository interface {
	// Core event operations
	GetEventById(id int64) (model.Event, error)
	GetEventByIssueId(issueId int64) (model.Event, error)
	CreateEvent(event *model.Event) error
	UpdateEvent(event *model.Event, eventLog *model.EventLog) error
	UpdateEventSize(eventId int64, size string) error
	
	// Event filtering and retrieval
	GetEvents(f EventFilter) ([]model.Event, error)
	GetMemberEvents(f EventFilter, memberId string) ([]model.Event, error)
	GetClientEvents(f EventFilter, clientId int64) ([]model.Event, error)
	GetClosedEventsByTimeRange(f EventFilter, startTime, endTime string) ([]JoinEvent, error)
	
	// Event utility functions
	GetEventClientId(eventId int64) (int64, error)
	
	// Event log operations
	CreateEventLog(eventLog *model.EventLog, conn *sqlx.Tx) error
	
	// Event status/action management
	ExistEventAction(action string) (bool, error)
	SetEventAction(eventLogId int64, action string, conn *sqlx.Tx) error
	ExistEventStatus(status string) (bool, error)
	SetEventStatus(eventId int64, status string, conn *sqlx.Tx) (sql.Result, error)
}

// ClientRepository defines the interface for client-related database operations
type ClientRepository interface {
	GetClientByOpenId(openId string) (model.Client, error)
	GetClientByLogtoId(logtoId string) (model.Client, error)
	CreateClient(client *model.Client) error
}

// RoleRepository defines the interface for role-related database operations
type RoleRepository interface {
	ExistRole(role string) (bool, error)
	SetMemberRole(memberId string, role string, conn *sqlx.Tx) error
}

// DatabaseManager defines the interface for database lifecycle management
type DatabaseManager interface {
	InitDB() error
	CloseDB() error
	SetDB(dbx *sqlx.DB)
	GetDB() *sqlx.DB
}