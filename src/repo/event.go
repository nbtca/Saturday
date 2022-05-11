package repo

import (
	"saturday/src/model"

	"github.com/Masterminds/squirrel"
)

var fields = []string{}

func getEventStatement() squirrel.SelectBuilder {
	return squirrel.Select("*").From("event").
		LeftJoin("event_status_relation USING (event_id)").
		LeftJoin("event_status USING (event_status_id)")
}
func getLogStatement() squirrel.SelectBuilder {
	return squirrel.Select("*").From("event").
		LeftJoin("event_action_relation USING (event_id").
		LeftJoin("event_action USING (event_action_id)")
}

func GetEventById(id int64) (model.Event, error) {
	// getEventSql, getEventArgs, _ := getEventStatement().Where(squirrel.Eq{"event_id": id}).ToSql()
	// getLogSql, getLogArgs, _ := getLogStatement().Where(squirrel.Eq{"event_id": id}).ToSql()
	// event := model.Event{}
	// logs := []model.EventLog{}
	// conn, err := db.Begin()
	// if err != nil {
	// 	return model.Event{}, err
	// }
	// conn.Exec(getEventSql, getEventArgs...)
	// if err := db.Get(&event, getEventSql, getEventArgs...); err != nil {
	// 	return model.Event{}, err
	// }
	// return event, nil
	return model.Event{}, nil
}
