package repo

import (
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nbtca/saturday/model"
)

func TestCreateEventPhotosChecksCapacityUnderLock(t *testing.T) {
	dsn := os.Getenv("TEST_DB_DATASOURCE")
	if dsn == "" {
		t.Skip("TEST_DB_DATASOURCE is not set")
	}
	db = sqlx.MustConnect("postgres", dsn)
	sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	var eventId int64
	if err := db.Get(&eventId, `INSERT INTO event (client_id, gmt_create, gmt_modified) VALUES (1, now(), now()) RETURNING event_id`); err != nil {
		t.Fatal(err)
	}
	db.MustExec(`INSERT INTO event_event_status_relation (event_id, event_status_id) SELECT $1, event_status_id FROM event_status WHERE status = 'accepted'`, eventId)
	t.Cleanup(func() {
		db.MustExec(`DELETE FROM event_photo WHERE event_id = $1`, eventId)
		db.MustExec(`DELETE FROM event_event_status_relation WHERE event_id = $1`, eventId)
		db.MustExec(`DELETE FROM event WHERE event_id = $1`, eventId)
	})

	full := errors.New("full")
	batch := make([]model.EventPhoto, 15)
	for i := range batch {
		batch[i] = model.EventPhoto{EventId: eventId, ObjectKey: "k", MemberId: "2333333333", GmtCreate: "2026-10-15T10:00:00+08:00"}
	}
	check := func(status string, count int) error {
		time.Sleep(50 * time.Millisecond)
		if status != "accepted" || count+len(batch) > 20 {
			return full
		}
		return nil
	}

	var wg sync.WaitGroup
	errs := make([]error, 2)
	for i := range errs {
		wg.Go(func() { errs[i] = CreateEventPhotos(eventId, batch, check) })
	}
	wg.Wait()

	var count int
	db.Get(&count, `SELECT count(*) FROM event_photo WHERE event_id = $1`, eventId)
	if count != 15 || (errs[0] == nil) == (errs[1] == nil) || !(errors.Is(errs[0], full) || errors.Is(errs[1], full)) {
		t.Errorf("count=%d errs=%v", count, errs)
	}
}
