package repo

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

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
