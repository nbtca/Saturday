package repo

import (
	"database/sql"
	"log"
	"saturday/model"
	"saturday/util"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var EventFields = []string{"event_id", "client_id", "model", "phone", "qq", "contact_preference",
	"problem", "member_id", "closed_by", "status", "gmt_create", "gmt_modified", "status"}

var EventLogFields = []string{"event_log_id", "description", "gmt_create", "member_id", "action"}

func getEventStatement() squirrel.SelectBuilder {
	return squirrel.Select(EventFields...).From("event").
		LeftJoin("event_event_status_relation USING (event_id)").
		LeftJoin("event_status USING (event_status_id)")
}
func getLogStatement() squirrel.SelectBuilder {
	return squirrel.Select(EventLogFields...).From("event_log").
		LeftJoin("event_event_action_relation USING (event_log_id)").
		LeftJoin("event_action USING (event_action_id)")
}

func GetEventById(id int64) (model.Event, error) {
	getEventSql, getEventArgs, _ := getEventStatement().Where(squirrel.Eq{"event_id": id}).ToSql()
	getLogSql, getLogArgs, _ := getLogStatement().Where(squirrel.Eq{"event_id": id}).ToSql()
	event := model.Event{}
	conn, err := db.Beginx()
	if err != nil {
		return model.Event{}, err
	}
	if err := conn.Get(&event, getEventSql, getEventArgs...); err != nil {
		return model.Event{}, err
	}
	if err := conn.Select(&event.Logs, getLogSql, getLogArgs...); err != nil {
		log.Println(err)
		return model.Event{}, err
	}
	if err = conn.Commit(); err != nil {
		conn.Rollback()
		return model.Event{}, err
	}
	return event, nil
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
