package repo

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Masterminds/squirrel"
)

func TestFilterEventsAddsMemberCondition(t *testing.T) {
	base := squirrel.Select("*").From("event_view as e")
	sql, args, err := filterEvents(base, EventFilter{Status: []string{"closed"}, MemberId: "3200421079"}).ToSql()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(sql, "e.status IN (?)") || !strings.Contains(sql, "e.member_id = ?") {
		t.Fatalf("unexpected sql: %s", sql)
	}
	if want := []any{"closed", "3200421079"}; !reflect.DeepEqual(args, want) {
		t.Fatalf("args = %v, want %v", args, want)
	}
}

func TestFilterEventsLeavesMemberOutWhenEmpty(t *testing.T) {
	sql, _, err := filterEvents(squirrel.Select("*").From("event_view as e"), EventFilter{}).ToSql()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(sql, "WHERE") {
		t.Fatalf("unexpected filter: %s", sql)
	}
}
