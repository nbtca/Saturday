package repo

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/jackc/pgx/v5/stdlib"

	"github.com/nbtca/saturday/util"
	"github.com/qustavo/sqlhooks/v2"
	"github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var db *sqlx.DB

// Hooks satisfies the sqlhooks.Hooks interface
type Hooks struct{}

// Before hook will print the query with it's args and return the context with the timestamp
func (h *Hooks) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	return context.WithValue(ctx, "begin", time.Now()), nil
}

// After hook will get the timestamp registered on the Before hook and print the elapsed time
func (h *Hooks) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	begin := ctx.Value("begin").(time.Time)
	util.Logger.WithFields(logrus.Fields{
		"query":   query,
		"args":    args,
		"elapsed": time.Since(begin),
		"id":      ctx.Value("uuid"),
	}).Debug("SQL executed")
	return ctx, nil
}

func InitDB() {
	var err error
	sql.Register("mysqlWithHooks", sqlhooks.Wrap(&stdlib.Driver{},&Hooks{}))
	db, err = sqlx.Connect("mysqlWithHooks", os.Getenv("DB_URL"))
	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres", driver)
	m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run

	if err != nil {
		util.Logger.Fatal(err)
	}
	db.SetMaxOpenConns(1000)               // The default is 0 (unlimited)
	db.SetMaxIdleConns(10)                 // defaultMaxIdleConns = 2
	db.SetConnMaxLifetime(time.Minute * 5) // 0, connections are reused forever.
}

func SetDB(dbx *sqlx.DB) {
	db = dbx
}

func CloseDB() {
	db.Close()
}
