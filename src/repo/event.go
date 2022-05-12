package repo

import (
	"log"
	"saturday/src/model"

	"github.com/Masterminds/squirrel"
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
		LeftJoin("event_action_relation USING (event_log_id)").
		LeftJoin("event_action USING (event_action_id)")
}

func GetEventById(id int64) (model.Event, error) {
	getEventSql, getEventArgs, _ := getEventStatement().Where(squirrel.Eq{"event_id": id}).ToSql()
	getLogSql, getLogArgs, _ := getLogStatement().Where(squirrel.Eq{"event_id": id}).ToSql()
	log.Println(getLogSql)
	event := model.Event{}
	conn, err := db.Beginx()
	if err != nil {
		return model.Event{}, err
	}
	if err := conn.Get(&event, getEventSql, getEventArgs...); err != nil {
		return model.Event{}, err
	}
	if err := conn.Select(&event.Logs, getLogSql, getLogArgs...); err != nil {
		return model.Event{}, err
	}
	// event.Logs = logs
	return event, nil
}
