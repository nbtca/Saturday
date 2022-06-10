package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
)

var pool *dockertest.Pool
var resource *dockertest.Resource

func isImageExist(name string) (bool, error) {
	out, err := exec.Command("docker", "image", "ls").Output()
	if err != nil {
		return false, err
	}
	outStr := string(out)
	searched := strings.Contains(outStr, "test_db")
	return searched, nil
}

func GetDB() (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error
	pool, err = dockertest.NewPool("")
	pool.MaxWait = time.Minute * 2
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	imageExists, _ := isImageExist("test_db")
	if imageExists {
		resource, err = pool.Run("test_db", "latest", []string{})
	} else {
		resource, err = pool.BuildAndRun("test_db", "../../assets/dockerfile", []string{})
	}
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		log.Println(resource.GetPort("3306/tcp"))
		db, err = sqlx.Connect("mysql", fmt.Sprintf("root:password@(localhost:%s)/saturday_test?multiStatements=true", resource.GetPort("3306/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	return db, nil
}

func CloseResource() {
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func SetSchema(db *sqlx.DB) error {
	b, err := ioutil.ReadFile("../../assets/saturday.sql")
	if err != nil {
		return err
	}
	schema := string(b)
	db.MustExec(schema)
	return nil
}
