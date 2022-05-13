package repo

import (
	"log"
	"saturday/src/model"
	"saturday/src/util"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

func CreateEventLog(eventLog *model.EventLog) error {
	eventLog.GmtCreate = util.GetDate()
	sql, args, _ := squirrel.Insert("event_log").Columns("event_id", "description", "member_id", "gmt_create").
		Values(eventLog.EventId, eventLog.Description, eventLog.MemberId, util.GetDate()).ToSql()
	conn, err := db.Beginx()
	if err != nil {
		return err
	}
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
	if err = conn.Commit(); err != nil {
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
	log.Println(eventLogId)
	sql := `INSERT INTO event_action_relation VALUES (?,(
		SELECT event_action_id FROM event_action WHERE action=?))
		ON DUPLICATE KEY UPDATE event_action_id=(
		SELECT event_action_id FROM event_action WHERE action= ? )`
	_, err := conn.Exec(sql, eventLogId, action, action)
	log.Println("Event action", err)
	return err
	// if err != nil {
	// 	return nil, err
	// }
	// var p Place
	// err = row.StructScan(&p)
}
