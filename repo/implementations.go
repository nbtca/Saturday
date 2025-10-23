package repo

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/nbtca/saturday/model"
)

// Concrete implementations of repository interfaces

// memberRepository implements MemberRepository interface
type memberRepository struct {
	db *sqlx.DB
	sq squirrel.StatementBuilderType
}

// NewMemberRepository creates a new member repository instance
func NewMemberRepository(db *sqlx.DB, sq squirrel.StatementBuilderType) MemberRepository {
	return &memberRepository{db: db, sq: sq}
}

func (r *memberRepository) GetMemberById(id string) (model.Member, error) {
	return GetMemberById(id)
}

func (r *memberRepository) GetMemberByLogtoId(logtoId string) (model.Member, error) {
	return GetMemberByLogtoId(logtoId)
}

func (r *memberRepository) GetMemberByGithubId(githubId string) (model.Member, error) {
	return GetMemberByGithubId(githubId)
}

func (r *memberRepository) GetMemberIdByLogtoId(logtoId string) (sql.NullString, error) {
	return GetMemberIdByLogtoId(logtoId)
}

func (r *memberRepository) GetMembers(offset uint64, limit uint64) ([]model.Member, error) {
	return GetMembers(offset, limit)
}

func (r *memberRepository) CreateMember(member *model.Member) error {
	return CreateMember(member)
}

func (r *memberRepository) UpdateMember(member model.Member) error {
	return UpdateMember(member)
}

func (r *memberRepository) ExistMember(id string) (bool, error) {
	return ExistMember(id)
}

// eventRepository implements EventRepository interface
type eventRepository struct {
	db *sqlx.DB
	sq squirrel.StatementBuilderType
}

// NewEventRepository creates a new event repository instance
func NewEventRepository(db *sqlx.DB, sq squirrel.StatementBuilderType) EventRepository {
	return &eventRepository{db: db, sq: sq}
}

func (r *eventRepository) GetEventById(id int64) (model.Event, error) {
	return GetEventById(id)
}

func (r *eventRepository) GetEventByIssueId(issueId int64) (model.Event, error) {
	return GetEventByIssueId(issueId)
}

func (r *eventRepository) CreateEvent(event *model.Event) error {
	return CreateEvent(event)
}

func (r *eventRepository) UpdateEvent(event *model.Event, eventLog *model.EventLog) error {
	return UpdateEvent(event, eventLog)
}

func (r *eventRepository) UpdateEventSize(eventId int64, size string) error {
	return UpdateEventSize(eventId, size)
}

func (r *eventRepository) GetEvents(f EventFilter) ([]model.Event, error) {
	return GetEvents(f)
}

func (r *eventRepository) GetMemberEvents(f EventFilter, memberId string) ([]model.Event, error) {
	return GetMemberEvents(f, memberId)
}

func (r *eventRepository) GetClientEvents(f EventFilter, clientId int64) ([]model.Event, error) {
	return GetClientEvents(f, clientId)
}

func (r *eventRepository) GetClosedEventsByTimeRange(f EventFilter, startTime, endTime string) ([]JoinEvent, error) {
	return GetClosedEventsByTimeRange(f, startTime, endTime)
}

func (r *eventRepository) GetEventClientId(eventId int64) (int64, error) {
	return GetEventClientId(eventId)
}

func (r *eventRepository) CreateEventLog(eventLog *model.EventLog, conn *sqlx.Tx) error {
	return CreateEventLog(eventLog, conn)
}

func (r *eventRepository) ExistEventAction(action string) (bool, error) {
	return ExistEventAction(action)
}

func (r *eventRepository) SetEventAction(eventLogId int64, action string, conn *sqlx.Tx) error {
	return SetEventAction(eventLogId, action, conn)
}

func (r *eventRepository) ExistEventStatus(status string) (bool, error) {
	return ExistEventStatus(status)
}

func (r *eventRepository) SetEventStatus(eventId int64, status string, conn *sqlx.Tx) (sql.Result, error) {
	return SetEventStatus(eventId, status, conn)
}

// clientRepository implements ClientRepository interface
type clientRepository struct {
	db *sqlx.DB
	sq squirrel.StatementBuilderType
}

// NewClientRepository creates a new client repository instance
func NewClientRepository(db *sqlx.DB, sq squirrel.StatementBuilderType) ClientRepository {
	return &clientRepository{db: db, sq: sq}
}

func (r *clientRepository) GetClientByOpenId(openId string) (model.Client, error) {
	return GetClientByOpenId(openId)
}

func (r *clientRepository) GetClientByLogtoId(logtoId string) (model.Client, error) {
	return GetClientByLogtoId(logtoId)
}

func (r *clientRepository) CreateClient(client *model.Client) error {
	return CreateClient(client)
}

// roleRepository implements RoleRepository interface
type roleRepository struct {
	db *sqlx.DB
	sq squirrel.StatementBuilderType
}

// NewRoleRepository creates a new role repository instance
func NewRoleRepository(db *sqlx.DB, sq squirrel.StatementBuilderType) RoleRepository {
	return &roleRepository{db: db, sq: sq}
}

func (r *roleRepository) ExistRole(role string) (bool, error) {
	return ExistRole(role)
}

func (r *roleRepository) SetMemberRole(memberId string, role string, conn *sqlx.Tx) error {
	return SetMemberRole(memberId, role, conn)
}

// databaseManager implements DatabaseManager interface
type databaseManager struct{}

// NewDatabaseManager creates a new database manager instance
func NewDatabaseManager() DatabaseManager {
	return &databaseManager{}
}

func (dm *databaseManager) InitDB() error {
	InitDB()
	return nil
}

func (dm *databaseManager) CloseDB() error {
	CloseDB()
	return nil
}

func (dm *databaseManager) SetDB(dbx *sqlx.DB) {
	SetDB(dbx)
}

func (dm *databaseManager) GetDB() *sqlx.DB {
	return db
}

// GetDB returns the global database connection (used during DI migration)
func GetDB() *sqlx.DB {
	return db
}