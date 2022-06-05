package repo

import (
	"database/sql"
	"log"
	"saturday/model"
	"saturday/util"

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
	fields := append(prefixedMember, append(prefixedAdmin, prefixedEvent...)...)
	return squirrel.Select(fields...).
		From("event_view as e").
		LeftJoin("member_view as m USING (member_id)").
		LeftJoin("member_view as a ON (e.closed_by=a.member_id)")
}

func getLogStatement() squirrel.SelectBuilder {
	return squirrel.Select(EventLogFields...).From("event_log_view")
}

type JoinEvent struct {
	Event  model.Event  `db:"event"`
	Member model.Member `db:"member"`
	Admin  model.Member `db:"admin"`
}

func (je JoinEvent) ToEvent() model.Event {
	event := je.Event
	event.Member = je.Member
	event.ClosedByMember = je.Admin
	return event
}

func GetEventById(id int64) (model.Event, error) {
	getEventSql, getEventArgs, _ := getEventStatement().Where(squirrel.Eq{"event_id": id}).ToSql()
	getLogSql, getLogArgs, _ := getLogStatement().Where(squirrel.Eq{"event_id": id}).ToSql()
	conn, err := db.Beginx()
	if err != nil {
		return model.Event{}, err
	}
	defer func() {
		if err != nil {
			conn.Rollback()
			db.Close()
		}
	}()
	joinEvent := JoinEvent{}
	if err := conn.Get(&joinEvent, getEventSql, getEventArgs...); err != nil {
		return model.Event{}, err
	}
	event := joinEvent.ToEvent()
	if err = conn.Select(&event.Logs, getLogSql, getLogArgs...); err != nil {
		return model.Event{}, err
	}
	if err = conn.Commit(); err != nil {
		return model.Event{}, err
	}
	return event, nil
}

func GetEvents(offset uint64, limit uint64) ([]model.Event, error) {
	getEventSql, getEventArgs, _ := getEventStatement().Offset(offset).Limit(limit).ToSql()
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

func UpdateEvent(event *model.Event, eventLog *model.EventLog) error {
	sql, args, _ := squirrel.Update("event").
		Set("model", event.Model).
		Set("phone", event.Phone).
		Set("qq", event.QQ).
		Set("contact_preference", event.ContactPreference).
		Set("problem", event.Problem).
		Set("member_id", event.MemberId).
		Set("closed_by", event.ClosedBy).
		Set("gmt_modified", event.GmtModified).
		Where(squirrel.Eq{"event_id": event.EventId}).ToSql()
	conn, err := db.Beginx()
	if err != nil {
		return err
	}
	if _, err = conn.Exec(sql, args...); err != nil {
		return err
	}
	log.Println(event.Status)
	if _, err = SetEventStatus(event.EventId, event.Status, conn); err != nil {
		return err
	}
	if err = CreateEventLog(eventLog, conn); err != nil {
		return err
	}
	if err = conn.Commit(); err != nil {
		conn.Rollback()
		return err
	}
	return nil
}

func CreateEventLog(eventLog *model.EventLog, conn *sqlx.Tx) error {
	sql, args, _ := squirrel.Insert("event_log").Columns("event_id", "description", "member_id", "gmt_create").
		Values(eventLog.EventId, eventLog.Description, eventLog.MemberId, util.GetDate()).ToSql()
	res, err := conn.Exec(sql, args...)
	if err != nil {
		return err
	}
	eventLogId, _ := res.LastInsertId()
	eventLog.EventLogId = int64(eventLogId)
	err = SetEventAction(eventLogId, eventLog.Action, conn)
	if err != nil {
		conn.Rollback()
		return err
	}
	return nil
}

func ExistEventAction(action string) (bool, error) {
	var count int
	err := db.Get(&count, "SELECT count(*) as count FROM event_action where action = ?", action)
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
	_, err := conn.Exec(sql, eventLogId, action, action)
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
	res, err := conn.Exec(sql, eventId, status, status)
	return res, err
}

func GetEventClientId(eventId int64) (int64, error) {
	var clientId int64
	err := db.Get(&clientId, "SELECT client_id FROM event WHERE event_id = ?", eventId)
	if err != nil {
		return 0, err
	}
	return clientId, nil
}
