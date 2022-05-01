package util

import (
	"fmt"
	"log"
	"os/exec"
	"reflect"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
)

func FieldsConstructor(q interface{}) []string {
	t := reflect.TypeOf(q)
	i := 0
	var res []string
	shouldAppendField := func(field reflect.StructField) (string, bool) {
		dbTag := t.Field(i).Tag.Get("json")
		// visibleTag := t.Field(i).Tag.Get("visible")
		if dbTag == "" {
			return "", false
		}
		// if visibleTag == "private" && visible != "private" {
		// 	return "", false
		// }
		return dbTag, true
	}
	if reflect.ValueOf(q).Kind() == reflect.Struct {
		for ; i < t.NumField()-1; i++ {
			if dbTag, should := shouldAppendField(t.Field(i)); should {
				res = append(res, dbTag)
			}
		}
	}
	if dbTag, should := shouldAppendField(t.Field(i)); should {
		res = append(res, dbTag)
	}
	return res
}

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
		resource, err = pool.BuildAndRun("test_db", "../../assets/Dockerfile", []string{})
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

func Close() {
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func GetDate() string {
	return time.Now().Format("2006-01-02 15:04:11")
}
