package repo

import (
	"database/sql"

	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/util"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var eventFields = []string{"event_id", "client_id", "model", "phone", "qq", "contact_preference",
	"problem", "member_id", "closed_by", "status", "gmt_create", "gmt_modified", "status"}

var EventLogFields = []string{"event_log_id", "description", "gmt_create", "member_id", "action"}

func getEventStatement() squirrel.SelectBuilder {
	prefixedMember := util.Prefixer("member", memberFields)
	prefixedAdmin := util.Prefixer("admin", memberFields)
	prefixedEvent := util.Prefixer("event", eventFields)
	// join columns
	fields := append(prefixedMember, append(prefixedAdmin, prefixedEvent...)...)

	// when column prefix is set, table must set with correspond alias(first letter of alias)
	// for example: when column prefix is "member", the table should be aliased to "m".
	return sq.Select(fields...).
		From("event_view as e").
		LeftJoin("member_view as m USING (member_id)").
		LeftJoin("member_view as a ON (e.closed_by=a.member_id)")
}

func getLogStatement() squirrel.SelectBuilder {
	return sq.Select(EventLogFields...).From("event_log_view")
}

/*
when a struct contains sub struct, the struct's field has the db tag of "struct'dbTag.subStruct.dbTag"
for example: the db tag here for event is event, the db tag for mode.Event has field "EventId" with db tag "event_id"
therefore the JoinEvent.Event.EventId has db tag of event.event_id.
*/
type JoinEvent struct {
	Event  model.Event      `db:"event"`
	Member model.NullMember `db:"member"`
	Admin  model.NullMember `db:"admin"`
}

// TODO this need to be refactored...
func (je JoinEvent) ToEvent() model.Event {
	event := je.Event
	event.Member = je.Member.PublicMember()
	event.ClosedByMember = je.Admin.PublicMember()
	return event
}

func GetEventById(id int64) (model.Event, error) {
	getEventSql, getEventArgs, _ := getEventStatement().Where(squirrel.Eq{"event_id": id}).ToSql()
	getLogSql, getLogArgs, _ := getLogStatement().Where(squirrel.Eq{"event_id": id}).ToSql()
	conn, err := db.Beginx()
	if err != nil {
		return model.Event{}, err
	}
	joinEvent := JoinEvent{}
	if err := conn.Get(&joinEvent, getEventSql, getEventArgs...); err != nil {
		if err == sql.ErrNoRows {
			return model.Event{}, nil
		}
		conn.Rollback()
		return model.Event{}, err
	}
	defer util.RollbackOnErr(err, conn)
	event := joinEvent.ToEvent()
	if err = conn.Select(&event.Logs, getLogSql, getLogArgs...); err != nil {
		return model.Event{}, err
	}
	if err = conn.Commit(); err != nil {
		return model.Event{}, err
	}
	return event, nil
}

type EventFilter struct {
	Offset uint64
	Limit  uint64
	Status string
	Order  string
}

func getEvents(f EventFilter, conditions ...squirrel.Eq) ([]model.Event, error) {
	stat := getEventStatement()
	if f.Status != "" {
		stat = stat.Where(squirrel.Eq{"e.status": f.Status})
	}
	for _, condition := range conditions {
		stat = stat.Where(condition)
	}
	getEventSql, getEventArgs, _ := stat.OrderBy("e.event_id " + f.Order).Offset(f.Offset).Limit(f.Limit).ToSql()
	joinEvent := []JoinEvent{}
	err := db.Select(&joinEvent, getEventSql, getEventArgs...)
	if err != nil {
		return []model.Event{}, err
	}
	events := make([]model.Event, len(joinEvent))
	for i, v := range joinEvent {
		events[i] = v.ToEvent()
	}
	return events, nil
}

func GetEvents(f EventFilter) ([]model.Event, error) {
	return getEvents(f)
}

func GetMemberEvents(f EventFilter, memberId string) ([]model.Event, error) {
	return getEvents(f, squirrel.Eq{"e.member_id": memberId})
}

func GetClientEvents(f EventFilter, clientId string) ([]model.Event, error) {
	return getEvents(f, squirrel.Eq{"e.client_id": clientId})
}

func UpdateEvent(event *model.Event, eventLog *model.EventLog) error {
	sql, args, _ := sq.Update("event").
		Set("model", event.Model).
		Set("phone", event.Phone).
		Set("qq", event.QQ).
		Set("contact_preference", event.ContactPreference).
		Set("problem", event.Problem).
		Set("member_id", event.MemberId).
		Set("closed_by", event.ClosedBy).
		// Set("gmt_modified", event.GmtModified).
		Where(squirrel.Eq{"event_id": event.EventId}).ToSql()
	conn, err := db.Beginx()
	if err != nil {
		return err
	}
	defer util.RollbackOnErr(err, conn)
	if _, err = conn.Exec(sql, args...); err != nil {
		return err
	}
	if _, err = SetEventStatus(event.EventId, event.Status, conn); err != nil {
		return err
	}
	if err = CreateEventLog(eventLog, conn); err != nil {
		return err
	}
	if err = conn.Commit(); err != nil {
		return err
	}
	return nil
}

func CreateEvent(event *model.Event) error {
	event.GmtCreate = util.GetDate()
	event.GmtModified = util.GetDate()
	createEventSql, args, _ := sq.Insert("event").Columns(
		"event_id", "client_id", "model", "phone", "qq",
		"contact_preference", "problem", "member_id", "closed_by",
		"gmt_create", "gmt_modified").Values(
		event.EventId, event.ClientId, event.Model, event.Phone, event.QQ,
		event.ContactPreference, event.Problem, event.MemberId, event.ClosedBy,
		event.GmtCreate, event.GmtModified).ToSql()
	conn, err := db.Beginx()
	if err != nil {
		return err
	}
	defer util.RollbackOnErr(err, conn)
	res, err := conn.Exec(createEventSql, args...)
	if err != nil {
		return err
	}
	event.EventId, _ = res.LastInsertId()
	if err = conn.Commit(); err != nil {
		return err

	}
	return nil
}

func CreateEventLog(eventLog *model.EventLog, conn *sqlx.Tx) error {
	sql, args, _ := sq.Insert("event_log").Columns("event_id", "description", "member_id", "gmt_create").
		Values(eventLog.EventId, eventLog.Description, eventLog.MemberId, util.GetDate()).ToSql()
	res, err := conn.Exec(sql, args...)
	if err != nil {
		return err
	}
	eventLogId, _ := res.LastInsertId()
	eventLog.EventLogId = int64(eventLogId)
	err = SetEventAction(eventLogId, eventLog.Action, conn)
	if err != nil {
		return err
	}
	return nil
}

func ExistEventAction(action string) (bool, error) {
	var count int
	sql, args, _ := sq.Select("count(*) as count").From("event_action").Where(squirrel.Eq{"action": action}).ToSql()
	err := db.Get(&count, sql, args...)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func SetEventAction(eventLogId int64, action string, conn *sqlx.Tx) error {
	sql := `INSERT INTO event_event_action_relation VALUES (?,(
		SELECT event_action_id FROM event_action WHERE action=?))
		ON DUPLICATE KEY UPDATE event_action_id=(
		SELECT event_action_id FROM event_action WHERE action= ? )`
	if db.DriverName() == "pqHooked" {
		sql = `INSERT INTO event_event_action_relation VALUES ($1,(
		SELECT event_action_id FROM event_action WHERE action=$2))
		ON CONFLICT (event_log_id) DO UPDATE SET event_action_id=(
		SELECT event_action_id FROM event_action WHERE action=$3 )`
	}
	_, err := conn.Exec(db.Rebind(sql), eventLogId, action, action)
	return err
}

func ExistEventStatus(status string) (bool, error) {
	var count int
	err := db.Get(&count, "SELECT count(*) as count FROM event_status where status = ?", status)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func SetEventStatus(eventId int64, status string, conn *sqlx.Tx) (sql.Result, error) {
	sql := `INSERT INTO event_event_status_relation (event_id, event_status_id)
	VALUES (?, (Select event_status_id from event_status where status = ?))
	ON DUPLICATE KEY UPDATE event_status_id=(SELECT event_status_id FROM event_status WHERE status=?)`

	if db.DriverName() == "pqHooked" {
		sql = `INSERT INTO event_event_status_relation (event_id, event_status_id)
		VALUES ($1, (Select event_status_id from event_status where status = $2))
		ON CONFLICT (event_id) DO UPDATE SET event_status_id=(SELECT event_status_id FROM event_status WHERE status=$3)`
	}

	return conn.Exec(db.Rebind(sql), eventId, status, status)
}

func GetEventClientId(eventId int64) (int64, error) {
	var clientId int64
	sql, args, _ := sq.Select("client_id").From("event").Where(squirrel.Eq{"event_id": eventId}).ToSql()
	err := db.Get(&clientId, sql, args...)
	if err != nil {
		return 0, err
	}
	return clientId, nil
}
